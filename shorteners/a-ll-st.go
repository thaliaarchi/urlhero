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

// Allst describes the Allstate a.ll.st link shortener.
var Allst = &Shortener{
	Name:   "a-ll-st",
	Host:   "a.ll.st",
	Prefix: "http://a.ll.st/",
	// Underscore is only allowed for vanity URLs.
	Pattern: regexp.MustCompile(`^[0-9A-Za-z_]+$`),
	CleanFunc: func(shortcode string, u *url.URL) string {
		// Remove /scmf/ID/ prefix:
		//   http://a.ll.st/scmf/OrMCe04Lcp0lODk0BD1FrBcO2E4FP0NMEHFGSZ--Pq5q7EdIBj5D0RZwQ0r5O5LJxfQiUmcjxE_yFyVUmcC7Ue52R7KC2DlT6j1Anuut1CVBLh2fal1IZic40eX4xD2dJTg/PrJJpv
		//   http://a.ll.st/scmf/OrMCe04Lcp0lODk2Bzg71hcM2079O8ZJEHE_NJu-wtVr7D9JB0U8qWl1RzYCRZPJxfQiUmcjxE_yF9swgNxdUAkTP4vGed-VJvLu3uityvkzL-5fGDGJnyV0iKf6RXKdJQ/hiddenworldofdata
		if strings.HasPrefix(shortcode, "scmf/") {
			shortcode = shortcode[strings.LastIndexByte(shortcode, '/')+1:]
		}
		return shortcode
	},
	LessFunc: func(a, b string) bool {
		// Sort 6-character generated codes before vanity codes.
		aVanity := len(a) != 6 || strings.Contains(a, "_")
		bVanity := len(b) != 6 || strings.Contains(b, "_")
		return (aVanity == bVanity && ((len(a) == len(b) && a < b) || len(a) < len(b))) ||
			(!aVanity && bVanity)
	},
}
