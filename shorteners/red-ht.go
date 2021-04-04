// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package shorteners

import (
	"net/url"
	"strings"

)

// Redht describes the Red Hat red.ht link shortener.
var Redht = &Shortener{
	Name:   "red-ht",
	Host:   "red.ht",
	Prefix: "https://red.ht/",
	// Pattern: regexp.MustCompile("^[0-9A-Za-z]+$"),
	Clean: func(shortcode string, u *url.URL) string {
		if strings.Contains(shortcode, ".") {
			return ""
		}
		if i := strings.IndexRune(shortcode, '’'); i != -1 {
			shortcode = shortcode[:i]
		}
		shortcode = strings.TrimSuffix(shortcode, "/")
		// Remove leftover query string parameters or HTML entities:
		//   http://red.ht/1zzgkXp&esheet=51687448&newsitemid=20170921005271&lan=en-US&anchor=Red+Hat+blog&index=5&md5=7ea962d15a0e5bf8e35f385550f4decb
		//   http://red.ht/13LslKt&quot
		if i := strings.IndexByte(shortcode, '&'); i != -1 {
			shortcode = shortcode[:i]
		}
		return shortcode
	},
	Less: func(a, b string) bool {
		// Sort generated codes before vanity codes.
		aVanity := strings.ContainsAny(a, "-_")
		bVanity := strings.ContainsAny(b, "-_")
		return (aVanity == bVanity && ((len(a) == len(b) && a < b) || len(a) < len(b))) ||
			(!aVanity && bVanity)
	},
}
