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
	Name:    "red-ht",
	Host:    "red.ht",
	Prefix:  "https://red.ht/",
	Pattern: regexp.MustCompile(`^[0-9A-Za-z\-_]+$`),
	CleanFunc: func(shortcode string, u *url.URL) string {
		// Exclude static files:
		//   http://red.ht/static/graphics/fish-404.png
		if strings.Contains(shortcode, ".") {
			return ""
		}
		// Remove social media @:
		//   https://red.ht/1H7Wyt1@sklososky@FuturePOV
		//   https://red.ht/3olOq1B@OpenRoboticsOrg
		if i := strings.IndexByte(shortcode, '@'); i != -1 {
			shortcode = shortcode[:i]
		}
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
