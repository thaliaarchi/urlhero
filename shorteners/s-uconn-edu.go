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

// SUconnEdu describes the University of Connecticut s.uconn.edu link
// shortener.
var SUconnEdu = &Shortener{
	Name:     "s-uconn-edu",
	Host:     "s.uconn.edu",
	Prefix:   "https://s.uconn.edu/", // Older links use http
	Alphabet: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
	Pattern:  regexp.MustCompile(`^[0-9A-Za-z\-_]+$`),
	CleanFunc: func(shortcode string, u *url.URL) string {
		// Exclude static files
		if strings.ContainsRune(shortcode, '/') {
			return ""
		}
		// Shortcodes are case-insensitive (both generated and vanity)
		return strings.ToLower(shortcode)
	},
	IsVanityFunc: func(shortcode string) bool {
		return strings.ContainsRune(shortcode, '-')
	},
	HasVanity: true,
}
