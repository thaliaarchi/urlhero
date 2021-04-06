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

// RedHt describes the Red Hat red.ht link shortener.
var RedHt = &Shortener{
	Name:     "red-ht",
	Host:     "red.ht",
	Prefix:   "https://red.ht/", // Older links use http
	Alphabet: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
	// Underscore and dash are only allowed for vanity URLs.
	Pattern: regexp.MustCompile(`^[0-9A-Za-z\-_]+$`),
	CleanFunc: func(shortcode string, u *url.URL) string {
		// Exclude static files
		if strings.ContainsRune(shortcode, '.') {
			return ""
		}
		// Remove social media @
		return trimAfterByte(shortcode, '@')
	},
	IsVanityFunc: func(shortcode string) bool {
		return strings.ContainsAny(shortcode, "-_")
	},
	HasVanity: true,
}
