// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package shorteners

import (
	"net/url"
	"regexp"
	"strings"
)

// GoHawaiiEdu describes the University of Hawaii go.hawaii.edu link
// shortener.
var GoHawaiiEdu = &Shortener{
	Name:     "go-hawaii-edu",
	Host:     "go.hawaii.edu",
	Prefix:   "https://go.hawaii.edu/", // Older links use http
	Alphabet: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
	Pattern:  regexp.MustCompile(`^[0-9A-Za-z]+$`),
	CleanFunc: func(shortcode string, u *url.URL) string {
		switch trimAfterByte(shortcode, '/') {
		case "admin", "submit":
			return ""
		}
		// Remove link preview
		return strings.TrimSuffix(shortcode, "+")
	},
	HasVanity: false,
}
