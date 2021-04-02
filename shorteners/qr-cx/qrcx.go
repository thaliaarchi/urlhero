// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package qrcx

import (
	"regexp"
	"strings"

	"github.com/andrewarchi/urlhero/shorteners"
)

// GetIAShortcodes queries all the shortcodes that have been archived on
// the Internet Archive.
func GetIAShortcodes() ([]string, error) {
	alpha := regexp.MustCompile("^[1-9A-HJ-Za-z]+$") // TODO
	clean := func(shortcode string) string {
		// Skip URL in path and files:
		//   http://qr.cx:80/http://qr.cx
		//   http://qr.cx:80/deleted.php
		if strings.Contains(shortcode, ".") || shortcode == "about:blank" {
			return ""
		}
		// Remove closing link:
		//   http://qr.cx/plvd]http:/qr.cx/plvd[/link]
		//   http://qr.cx/plvd]click
		//   http://qr.cx/)
		if i := strings.IndexByte(shortcode, ']'); i != -1 {
			shortcode = shortcode[:i]
		}
		if i := strings.IndexByte(shortcode, ')'); i != -1 {
			shortcode = shortcode[:i]
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
	}
	return shorteners.GetIAShortcodes("qr.cx", alpha, clean, nil)
}
