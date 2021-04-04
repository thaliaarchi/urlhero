// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"sort"
	"strings"
	"time"

	bitly "github.com/andrewarchi/urlhero/shorteners/bit-ly"
)

const usage = `Usage:
	lookupaliases hosts <hosts>...
	lookupaliases reverse <rdns.json>`

func main() {
	if len(os.Args) < 3 {
		printUsage()
	}
	switch os.Args[1] {
	case "hosts":
		lookup(os.Args[2:])
	case "reverse":
		if len(os.Args) != 3 {
			printUsage()
		}
		records, err := reverseLookup(os.Args[2])
		if len(records) != 0 {
			for _, r := range records {
				fmt.Println(r)
			}
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, usage)
	os.Exit(2)
}

func lookup(hosts []string) {
	type hostIP struct {
		host string
		ips  []net.IP
	}
	var hostInfo []hostIP
	ipMap := make(map[[16]byte][]string)
	var errs []error

	fmt.Println("Resolved IP addresses:")
	for _, host := range hosts {
		ips, err := net.LookupIP(host)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		sort.Slice(ips, func(i, j int) bool {
			return bytes.Compare(ips[i], ips[j]) < 0
		})
		hostInfo = append(hostInfo, hostIP{host, ips})

		fmt.Printf("  %s:\t", host)
		for i, ip := range ips {
			fmt.Print(ip)
			if i != len(ips)-1 {
				fmt.Print(", ")
			}
			ip6 := ip.To16()
			if ip6 == nil {
				continue
			}
			var b [16]byte
			copy(b[:], ip6)
			ipMap[b] = append(ipMap[b], host)
		}
		fmt.Println()
	}

	fmt.Println("Shared IP addresses:")
	ipsSorted := make([][16]byte, len(ipMap))
	for ip := range ipMap {
		ipsSorted = append(ipsSorted, ip)
	}
	sort.Slice(ipsSorted, func(i, j int) bool {
		return bytes.Compare(ipsSorted[i][:], ipsSorted[j][:]) < 0
	})
	for _, ip := range ipsSorted {
		hosts := ipMap[ip]
		if len(hosts) > 1 {
			fmt.Printf("  %v:", net.IP(ip[:]))
			for i, host := range hosts {
				fmt.Print(host)
				if i != len(ipsSorted)-1 {
					fmt.Print(", ")
				}
			}
			fmt.Println()
		}
	}

	fmt.Println("bit.ly aliases:")
	for _, h := range hostInfo {
		if bitly.IsIPAlias(h.ips...) {
			fmt.Printf("  %s\n", h.host)
		}
	}

	if len(errs) != 0 {
		fmt.Println("Errors:")
		for _, err := range errs {
			fmt.Printf("  %v\n", err)
		}
	}
}

type DNSRecord struct {
	Timestamp int64  `json:"timestamp,string"` // Unix timestamp in seconds
	IP        net.IP `json:"name"`
	Host      string `json:"value"`
	Type      string `json:"type"` // "ptr"
}

func (r *DNSRecord) String() string {
	t := time.Unix(r.Timestamp, 0).UTC().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("%s\t%v\t%s\n", t, r.IP, r.Host)
}

// reverseLookup finds all hosts that resolve to a bit.ly IP, using the
// Project Sonar reverse DNS dataset from
// https://opendata.rapid7.com/sonar.rdns_v2/.
func reverseLookup(filename string) ([]DNSRecord, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := io.Reader(f)
	if strings.HasSuffix(filename, ".gz") {
		gr, err := gzip.NewReader(f)
		if err != nil {
			return nil, err
		}
		defer gr.Close()
		r = gr
	}

	d := json.NewDecoder(r)
	d.DisallowUnknownFields()
	var aliases []DNSRecord
	var record DNSRecord
	start := time.Now()
	last := start
	n := 0
	fmt.Println("bit-ly aliases:")
	for {
		if err := d.Decode(&record); err != nil {
			if err == io.EOF {
				break
			}
			return aliases, err
		}
		if bitly.IsIPAlias(record.IP) {
			aliases = append(aliases, record)
			fmt.Println(record)
		}
		n++
		if n%1000000 == 0 {
			now := time.Now()
			fmt.Printf("Processed %d records in %v (%v elapsed)\n", n, now.Sub(last), now.Sub(start))
			last = now
		}
	}
	fmt.Printf("Processed %d records (%v elapsed)\n", n, time.Since(start))
	return aliases, nil
}
