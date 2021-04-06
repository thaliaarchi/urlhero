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

// Redirect preview:
//   https://rb.gy/<shortcode>+
// Shortcodes appear to max at 6 characters and be case insensitive.

// Rbgy describes the rb-gy link shortener, which is powered by
// Rebrandly.
var Rbgy = &Shortener{
	Name:     "rb-gy",
	Host:     "rb.gy",
	Prefix:   "https://rb.gy/", // Older links use http
	Alphabet: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
	Pattern:  regexp.MustCompile(`^[0-9A-Za-z]+$`),
	CleanFunc: func(shortcode string, u *url.URL) string {
		// Remove redirect preview + and trailing junk
		shortcode = trimAfterAny(shortcode, "+/._-:$@◄")
		// Remove characters not in alphabet
		shortcode = rbgyNonAlpha.ReplaceAllLiteralString(shortcode, "")
		// Remove directly-concatenated junk
		if len(shortcode) > 6 {
			shortcode = shortcode[:6]
		}
		return strings.ToLower(shortcode)
	},
	HasVanity: false,
}

var rbgyNonAlpha = regexp.MustCompile("[^0-9A-Za-z]+")
