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

// Redht describes the Red Hat red.ht link shortener.
var Redht = &Shortener{
	Name:   "red-ht",
	Host:   "red.ht",
	Prefix: "https://red.ht/",
	// TODO meaning of @ is unknown; should it be stripped?
	Pattern: regexp.MustCompile(`^[0-9A-Za-z\-_@]+$`),
	CleanFunc: func(shortcode string, u *url.URL) string {
		if strings.Contains(shortcode, ".") {
			return ""
		}
		shortcode = strings.TrimSuffix(shortcode, "/")
		return shortcode
	},
	LessFunc: func(a, b string) bool {
		// Sort generated codes before vanity codes.
		aVanity := strings.ContainsAny(a, "-_")
		bVanity := strings.ContainsAny(b, "-_")
		return (aVanity == bVanity && ((len(a) == len(b) && a < b) || len(a) < len(b))) ||
			(!aVanity && bVanity)
	},
}
