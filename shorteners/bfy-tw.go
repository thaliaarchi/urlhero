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
	Prefix:   "https://bfy.tw/", // Older links use http
	Alphabet: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
	Pattern:  regexp.MustCompile(`^[0-9A-Za-z]+$`),
	CleanFunc: func(shortcode string, u *url.URL) string {
		return bfytwTrailing.ReplaceAllLiteralString(shortcode, "")
	},
	HasVanity: false,
}

var bfytwTrailing = regexp.MustCompile(`/?(?:https?://.+|[/.].*|favicon.ico|robots.txt|ip\d+\.\d+\.\d+\.\d+)?=?$`)
