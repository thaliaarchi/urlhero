// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package bitly

import (
	"encoding/binary"
	"net"
)

// IsAlias checks whether the domain is a bit.ly alias.
func IsAlias(domain string) (bool, error) {
	// Check if the IP is in the range 67.199.248.10 to 67.199.248.13.
	lo := binary.BigEndian.Uint32([]byte{67, 199, 248, 10})
	hi := binary.BigEndian.Uint32([]byte{67, 199, 248, 13})
	ips, err := net.LookupIP(domain)
	if err != nil {
		return false, err
	}
	for _, ip := range ips {
		ip = ip.To4()
		if ip == nil {
			continue
		}
		ip4 := binary.BigEndian.Uint32(ip)
		if lo <= ip4 && ip4 <= hi {
			return true, nil
		}
	}
	return false, nil
}
