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
	Name:     "a-ll-st",
	Host:     "a.ll.st",
	Prefix:   "http://a.ll.st/",
	Alphabet: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
	// Underscore is only allowed for vanity URLs.
	Pattern: regexp.MustCompile(`^[0-9A-Za-z_]+$`),
	CleanFunc: func(shortcode string, u *url.URL) string {
		if strings.HasPrefix(shortcode, "scmf/") {
			shortcode = shortcode[strings.LastIndexByte(shortcode, '/')+1:]
		}
		return shortcode
	},
	IsVanityFunc: func(shortcode string) bool {
		return len(shortcode) != 6 || strings.ContainsRune(shortcode, '_')
	},
	HasVanity: true,
}
