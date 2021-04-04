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

type Shortener struct {
	Name    string
	Host    string
	Prefix  string
	Pattern *regexp.Regexp
	Clean   CleanFunc
	Less    LessFunc
}

type CleanFunc func(shortcode string, u *url.URL) string
type LessFunc func(i, j string) bool

var Shorteners = []*Shortener{
	Allst,
	Debli,
	Qrcx,
	Redht,
}
