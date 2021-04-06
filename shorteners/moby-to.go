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

// moby.to is a link shortener for sharing images.
//
// These suffixes show the picture at various sizes:
//   http://moby.to/<shortcode>:view
//   http://moby.to/<shortcode>:full
//   http://moby.to/<shortcode>:square
//   http://moby.to/<shortcode>:small
//   http://moby.to/<shortcode>:large
//   http://moby.to/<shortcode>:thumb
//   http://moby.to/<shortcode>:thumbnail

// MobyTo describes the Mobypicture moby.to link shortener.
var MobyTo = &Shortener{
	Name:     "moby-to",
	Host:     "moby.to",
	Prefix:   "http://moby.to/",
	Alphabet: "0123456789abcdefghijklmnopqrstuvwxyz",
	Pattern:  regexp.MustCompile("^[0-9a-z]+$"),
	CleanFunc: func(shortcode string, u *url.URL) string {
		if strings.ContainsRune(shortcode, '/') {
			return ""
		}
		// Remove : suffix and trailing junk
		if i := strings.IndexAny(shortcode, ":<«-+*."); i != -1 {
			shortcode = shortcode[:i]
		}
		return strings.ToLower(shortcode)
	},
	HasVanity: false,
}
