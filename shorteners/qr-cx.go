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

// Qrcx describes the qr.cx link shortener.
var Qrcx = &Shortener{
	Name:    "qr-cx",
	Host:    "qr.cx",
	Prefix:  "http://qr.cx/",
	Pattern: regexp.MustCompile(`^[1-9A-HJ-Za-z]+$`), // TODO verify that I is missing
	CleanFunc: func(shortcode string, u *url.URL) string {
		// Skip URL in path and files:
		//   http://qr.cx:80/http://qr.cx
		//   http://qr.cx:80/deleted.php
		if strings.Contains(shortcode, ".") || shortcode == "about:blank" {
			return ""
		}
		shortcode = strings.TrimSuffix(shortcode, "/")
		shortcode = strings.TrimSuffix(shortcode, "+") // Alias for /get
		dir, path := shortcode, ""
		if i := strings.IndexByte(shortcode, '/'); i != -1 {
			dir, path = shortcode[:i], shortcode[i+1:]
		}
		switch dir {
		case "admin", "api", "dataset", "img", "qr", "twitterjs":
			return ""
		}
		if path == "get" {
			shortcode = dir
		}
		return shortcode
	},
}
