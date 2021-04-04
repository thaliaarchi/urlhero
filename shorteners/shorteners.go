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

type Shortener struct {
	Name      string
	Host      string
	Prefix    string
	Pattern   *regexp.Regexp
	CleanFunc CleanFunc
	LessFunc  LessFunc
}

type CleanFunc func(shortcode string, u *url.URL) string
type LessFunc func(i, j string) bool

var Shorteners = []*Shortener{
	Allst,
	Debli,
	Qrcx,
	Redht,
}

func (s *Shortener) Clean(shortURL string) (string, error) {
	u, err := url.Parse(shortURL)
	if err != nil {
		return "", err
	}
	shortcode := strings.TrimPrefix(u.Path, "/")
	// Exclude placeholders:
	//   https://deb.li/<key>
	//   https://deb.li/<name>
	if len(shortcode) >= 2 && shortcode[0] == '<' && shortcode[len(shortcode)-1] == '>' {
		return "", nil
	}
	// Remove trailing junk:
	//   http://a.ll.st/Instagram","isCrawlable":true,"thumbnail
	//   http://qr.cx/plvd]http:/qr.cx/plvd[/link]
	//   http://qr.cx/)
	//   https://red.ht/sig>
	//   https://red.ht/1zzgkXp&esheet=51687448&newsitemid=20170921005271&lan=en-US&anchor=Red+Hat+blog&index=5&md5=7ea962d15a0e5bf8e35f385550f4decb
	//   https://red.ht/13LslKt&quot
	//   https://red.ht/2k3DNz3’
	//   https://red.ht/21Krw4z%C2%A0   (nbsp)
	if i := strings.IndexAny(shortcode, "\"])>&’\u00a0"); i != -1 {
		shortcode = shortcode[:i]
	}
	if shortcode == "" {
		return "", nil
	}
	if s.CleanFunc != nil {
		shortcode = s.CleanFunc(shortcode, u)
	}
	shortcode = strings.TrimSuffix(shortcode, "/")
	switch shortcode {
	case "favicon.ico", "robots.txt":
		return "", nil
	}
	return shortcode, nil
}
