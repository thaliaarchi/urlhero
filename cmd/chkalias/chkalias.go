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
	"strconv"
	"strings"
	"time"

	bitly "github.com/andrewarchi/urlhero/shorteners/bit-ly"
)

const usage = `Usage:
	chkalias hosts <hosts>...
	chkalias lookup <dns.json>`

func main() {
	if len(os.Args) < 3 {
		printUsage()
	}
	switch os.Args[1] {
	case "host", "hosts":
		lookup(os.Args[2:])
	case "lookup":
		if len(os.Args) != 3 {
			printUsage()
		}
		records, err := lookupProjectSonar(os.Args[2])
		if len(records) != 0 {
			fmt.Println("bit-ly aliases:")
			for _, r := range records {
				fmt.Println(r)
			}
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		printUsage()
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
	ipKeys := make([][16]byte, len(ipMap))
	for ip := range ipMap {
		ipKeys = append(ipKeys, ip)
	}
	sort.Slice(ipKeys, func(i, j int) bool {
		return bytes.Compare(ipKeys[i][:], ipKeys[j][:]) < 0
	})
	for _, ip := range ipKeys {
		hosts := ipMap[ip]
		if len(hosts) > 1 {
			fmt.Printf("  %v:", net.IP(ip[:]))
			for i, host := range hosts {
				fmt.Print(host)
				if i != len(ipKeys)-1 {
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
	Time time.Time
	Host string
	IP   net.IP
}

func (r *DNSRecord) String() string {
	t := r.Time.Format("2006-01-02 15:04:05")
	return fmt.Sprintf("%s\t%v\t%s\n", t, r.Host, r.IP)
}

// lookupProjectSonar finds all hosts that resolve to a bit.ly IP, using
// the Project Sonar forward and reverse DNS datasets from
// https://opendata.rapid7.com/sonar.fdns_v2/ and
// https://opendata.rapid7.com/sonar.rdns_v2/.
func lookupProjectSonar(filename string) ([]DNSRecord, error) {
	// TODO this needs significant work before it can be useful. Forward
	// DNS has all the intermediate pointers, so the IP address of a host
	// cannot be determined via a single read pass and a graph database
	// would be needed. Reverse DNS only contains the authoritative
	// mappings, so all bit.ly aliases are excluded.
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

	type dnsRecord struct {
		Timestamp string `json:"timestamp"` // Unix timestamp in seconds
		Name      string `json:"name"`      // FDNS: host, RDNS: IP
		Value     string `json:"value"`     // FDNS: IP,   RDNS: host
		Type      string `json:"type"`      // FDNS: "ns", RDNS: "ptr"
	}

	d := json.NewDecoder(r)
	d.DisallowUnknownFields()
	var aliases []DNSRecord
	start := time.Now()
	last := start
	n := 1
	for {
		var r dnsRecord
		if err := d.Decode(&r); err != nil {
			if err == io.EOF {
				break
			}
			if _, ok := err.(*json.SyntaxError); ok {
				fmt.Fprintf(os.Stderr, "record %d: %v\n", n, err)
				continue
			}
			return aliases, err
		}

		timestamp, err := strconv.ParseInt(r.Timestamp, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "record %d: %v: %v\n", n, r, err)
		}
		var host, ip string
		switch r.Type {
		case "ns": // forward DNS
			host, ip = r.Name, r.Value
		case "ptr": // reverse DNS
			host, ip = r.Value, r.Name
		default:
			fmt.Fprintf(os.Stderr, "record %d: unrecognized type %s: %v\n", n, r.Type, r)
			continue
		}
		record := DNSRecord{
			Time: time.Unix(timestamp, 0).UTC(),
			Host: host,
			IP:   net.ParseIP(ip),
		}

		if bitly.IsIPAlias(record.IP) {
			aliases = append(aliases, record)
			fmt.Println("Found:", record)
		}
		if n%1000000 == 0 {
			now := time.Now()
			fmt.Printf("Processed %d records in %v (%v elapsed)\n", n, now.Sub(last), now.Sub(start))
			last = now
		}
		n++
	}
	fmt.Printf("Processed %d records (%v elapsed)\n", n-1, time.Since(start))
	return aliases, nil
}
