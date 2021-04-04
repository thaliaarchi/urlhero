// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package shorteners

import (
	"net/url"
	"regexp"
)

// Bfytw describes the bfy.tw link shortener. All URLs redirect to
// https://lmgtfy.app/?q=<query>
var Bfytw = &Shortener{
	Name:     "bfy-tw",
	Host:     "bfy.tw",
	Prefix:   "https://bfy.tw/",
	Alphabet: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
	Pattern:  regexp.MustCompile(`^[0-9A-Za-z]+$`),
	CleanFunc: func(shortcode string, u *url.URL) string {
		// Remove trailing junk:
		//   https://bfy.tw/80xn=
		//   https://bfy.tw/fb/7rt7
		//   https://bfy.tw/3hQy...You
		//   https://bfy.tw/BFsxrobots.txt
		//   https://bfy.tw/D9lj/favicon.ico
		//   https://bfy.tw/5PrLhttps://bfy.tw/5PrL
		//   https://bfy.tw/Ej4D/wordpress/wp-content/uploads/kisaflo-top-loog.png
		// This URL is incorrectly cleaned:
		//   https://bfy.tw/4jz9ip124.41.235.255
		return bfytwPattern.ReplaceAllLiteralString(shortcode, "")
	},
}

var bfytwPattern = regexp.MustCompile("/?(?:https?://.+|[/.].*|favicon.ico|robots.txt)?=?$")
