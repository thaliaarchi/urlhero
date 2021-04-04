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
	Prefix:   "https://deb.li/",
	Alphabet: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
	Pattern:  regexp.MustCompile(`^(?:[0-9A-Za-z]+|.+@.+)$`),
	CleanFunc: func(shortcode string, u *url.URL) string {
		// Keep mailing list redirects as-is:
		//   https://deb.li/4BE7F84D.5040104@bzed.de
		if strings.Contains(shortcode, "@") {
			return shortcode
		}
		// Remove preview prefix:
		//   https://deb.li/p/debian
		//   https://deb.li/p/1r8d
		shortcode = strings.TrimPrefix(shortcode, "p/")
		// Exclude static files and strange URLs:
		//   https://deb.li/static/pics/openlogo-50.png
		//   https://deb.li/imprint.html
		//   https://deb.li/log%20dari%20training%20Debian%20Women%20dengan%20tema%20%22Debian%20package%20informations%22%20dini%20hari%20tadi%20dapat%20dilihat%20di%20http://meetbot.debian.net/debian-women/2010/debian-women.2010-12-16-20.09.log.html
		if strings.Contains(shortcode, "/") || strings.Contains(u.Path, ".") {
			return ""
		}
		return shortcode
	},
}
