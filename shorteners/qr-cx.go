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
	HasVanity: false,
}

func cleanQrcx(shortcode string, u *url.URL) string {
	// Get URL from QR API query string:
	//   http://qr.cx/qr/php/qr_img.php?e=M&s=9&d=http://qr.cx/<shortcode>
	if shortcode == "qr/php/qr_img.php" {
		if d := u.Query().Get("d"); d != "" {
			u2, err := url.Parse(d)
			if err != nil || u2.Hostname() != "qr.cx" {
				return ""
			}
			return cleanURL(u2, "qr.cx", cleanQrcx)
		}
	}
	// Remove file after shortcode
	shortcode = qrcxFiles.ReplaceAllLiteralString(shortcode, "")
	// Remove redirect preview
	return strings.TrimSuffix(shortcode, "+")
}

var qrcxFiles = regexp.MustCompile(`(?:^|/)(?:admin|api|api\.php|dataset|deleted\.php|get|img|piwik\.php|qr|twitterjs)(?:$|/.*)`)
