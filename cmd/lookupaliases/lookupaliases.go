// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"sort"

	bitly "github.com/andrewarchi/urlhero/shorteners/bit-ly"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s hosts...\n", os.Args[0])
		os.Exit(2)
	}
	hosts := os.Args[1:]

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

		fmt.Printf("  %s: ", host)
		for i, ip := range ips {
			if i != 0 {
				fmt.Print(",")
			}
			fmt.Printf(" %v", ip)
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
				if i != 0 {
					fmt.Print(",")
				}
				fmt.Printf(" %s", host)
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
