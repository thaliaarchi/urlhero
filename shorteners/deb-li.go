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

// Source: https://github.com/bzed/go.debian.net
// Beta test announcement: https://lists.debian.org/debian-devel/2010/05/msg00248.html
//
// deb.li has three types of URLs:
//   generated shortcode:
//     https://deb.li/ijEl
//   vanity shortcode (case sensitive):
//     https://deb.li/DTAuthors
//   Debian mailing list redirect:
//     https://deb.li/4BE7F84D.5040104@bzed.de
//     https://deb.li/<message-id> -> https://lists.debian.org/msgid-search/<message-id>

// Debli describes the Debian deb.li link shortener.
var Debli = &Shortener{
	Name:     "deb-li",
	Host:     "deb.li",
	Prefix:   "https://deb.li/", // Older links use http
	Alphabet: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
	Pattern:  regexp.MustCompile(`^(?:[0-9A-Za-z]+|.+@.+)$`),
	CleanFunc: func(shortcode string, u *url.URL) string {
		// Keep mailing list redirects as-is
		if strings.ContainsRune(shortcode, '@') {
			return shortcode
		}
		// Remove redirect preview
		shortcode = strings.TrimPrefix(shortcode, "p/")
		// Exclude static files and strange URLs
		if strings.ContainsRune(shortcode, '/') || strings.ContainsRune(shortcode, '.') {
			return ""
		}
		return shortcode
	},
	HasVanity: true,
}
