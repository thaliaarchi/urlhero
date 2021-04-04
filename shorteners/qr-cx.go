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
	Name:      "qr-cx",
	Host:      "qr.cx",
	Prefix:    "http://qr.cx/",
	Alphabet:  "123456789ABCDEFGHJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
	Pattern:   regexp.MustCompile(`^[1-9A-HJ-Za-z]+$`),
	CleanFunc: cleanQrcx,
}

func cleanQrcx(shortcode string, u *url.URL) string {
	// Get URL from QR API query string:
	//   http://qr.cx/qr/php/qr_img.php?e=M&s=9&d=http://qr.cx/1oz
	if shortcode == "qr/php/qr_img.php" {
		if d := u.Query().Get("d"); d != "" {
			u2, err := url.Parse(d)
			if err != nil || u2.Hostname() != "qr.cx" {
				return ""
			}
			return cleanURL(u2, cleanQrcx)
		}
	}
	// Remove file after shortcode:
	//   http://qr.cx/itZ/api.php
	//   http://qr.cx/uqn/piwik.php
	if i := strings.LastIndexByte(shortcode, '/'); i != -1 {
		switch shortcode[i+1:] {
		case "api.php", "piwik.php":
			shortcode = shortcode[:i]
		}
	}
	// Skip URL in path and files:
	//   http://qr.cx/http://qr.cx
	//   http://qr.cx/deleted.php
	if strings.Contains(shortcode, ".") || shortcode == "about:blank" {
		return ""
	}
	// Remove link previews:
	//   http://qr.cx/tEv/get
	//   http://qr.cx/sQ2U+
	shortcode = strings.TrimSuffix(shortcode, "/get")
	shortcode = strings.TrimSuffix(shortcode, "+")
	// Exclude static files:
	switch shortcode {
	case "admin", "api", "dataset", "img", "qr", "twitterjs":
		return ""
	}
	if strings.Contains(shortcode, "/") {
		return ""
	}
	return shortcode
}
