// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package shorteners provides utilities for retrieving information
// about URL shortening websites.
package shorteners

import (
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"github.com/andrewarchi/urlhero/ia"
)

type Shortener struct {
	Name         string
	Host         string
	Prefix       string
	Alphabet     string
	Pattern      *regexp.Regexp
	CleanFunc    CleanFunc
	IsVanityFunc IsVanityFunc
	HasVanity    bool
}

type CleanFunc func(shortcode string, u *url.URL) string
type IsVanityFunc func(shortcode string) bool

var Shorteners = []*Shortener{
	Allst,
	Bfytw,
	Debli,
	GoHawaiiEdu,
	Qrcx,
	RedHt,
	SUconnEdu,
}

var Lookup = make(map[string]*Shortener)

func init() {
	for _, s := range Shorteners {
		if _, ok := Lookup[s.Name]; ok {
			panic(fmt.Errorf("multiple shorteners with name %s", s.Name))
		}
		if _, ok := Lookup[s.Host]; ok {
			panic(fmt.Errorf("multiple shorteners with host %s", s.Host))
		}
		Lookup[s.Name] = s
		Lookup[s.Host] = s
	}
}

// Clean extracts the shortcode from a URL. An empty string is returned
// when no shortcode can be found.
func (s *Shortener) Clean(shortURL string) (string, error) {
	u, err := url.Parse(shortURL)
	if err != nil {
		return "", err
	}
	return cleanURL(u, s.CleanFunc), nil
}

// CleanURL extracts the shortcode from a URL. An empty string is
// returned when no shortcode can be found.
func (s *Shortener) CleanURL(u *url.URL) string {
	return cleanURL(u, s.CleanFunc)
}

func cleanURL(u *url.URL, clean CleanFunc) string {
	shortcode := strings.TrimPrefix(u.Path, "/")
	// Exclude placeholders:
	//   https://deb.li/<key>
	//   https://deb.li/<name>
	if len(shortcode) >= 2 && shortcode[0] == '<' && shortcode[len(shortcode)-1] == '>' {
		return ""
	}
	// Remove trailing periods:
	//   https://bfy.tw/7JAH.
	//   https://bfy.tw/LOr7...
	shortcode = strings.TrimRight(shortcode, ".")
	// Remove trailing junk:
	//   http://a.ll.st/Instagram","isCrawlable":true,"thumbnail
	//   http://qr.cx/plvd]http:/qr.cx/plvd[/link]
	//   http://qr.cx/)
	//   https://red.ht/sig>
	//   https://red.ht/1zzgkXp&esheet=51687448&newsitemid=20170921005271&lan=en-US&anchor=Red+Hat+blog&index=5&md5=7ea962d15a0e5bf8e35f385550f4decb
	//   https://red.ht/13LslKt&quot
	//   http://go.hawaii.edu/j7L;
	//   https://red.ht/2k3DNz3’
	//   https://deb.li/log%20dari%20training%20Debian%20Women%20dengan%20tema%20%22Debian%20package%20informations%22%20dini%20hari%20tadi%20dapat%20dilihat%20di%20http://meetbot.debian.net/debian-women/2010/debian-women.2010-12-16-20.09.log.html
	//   https://red.ht/21Krw4z%C2%A0   (nbsp)
	// Escape sequences are non-breaking and zero-width spaces
	if i := strings.IndexAny(shortcode, "\"])>&;’ \u00a0\u200B"); i != -1 {
		shortcode = shortcode[:i]
	}
	shortcode = strings.TrimSuffix(shortcode, "/")
	if isCommonFile(shortcode) {
		return ""
	}
	if clean != nil {
		shortcode = clean(shortcode, u)
	}
	shortcode = strings.TrimSuffix(shortcode, "/")
	if isCommonFile(shortcode) {
		return ""
	}
	return shortcode
}

func isCommonFile(shortcode string) bool {
	switch shortcode {
	case "", "favicon.ico", "robots.txt":
		return true
	}
	return false
}

// CleanURLs extracts, deduplicates, and sorts the shortcodes in slice
// of URLs.
func (s *Shortener) CleanURLs(urls []string) ([]string, error) {
	shortcodesMap := make(map[string]struct{})
	var shortcodes []string
	for _, shortURL := range urls {
		shortcode, err := s.Clean(shortURL)
		if err != nil {
			return nil, err
		} else if shortcode == "" {
			continue
		}
		if s.Pattern != nil && !s.Pattern.MatchString(shortcode) {
			return nil, fmt.Errorf("%s: shortcode %q does not match alphabet %s after cleaning: %q", s.Name, shortcode, s.Pattern, shortURL)
		}
		if _, ok := shortcodesMap[shortcode]; !ok {
			shortcodesMap[shortcode] = struct{}{}
			shortcodes = append(shortcodes, shortcode)
		}
	}
	s.Sort(shortcodes)
	return shortcodes, nil
}

// IsVanity returns true when a shortcode is a vanity code. There are
// many false negatives for vanity codes that are programmatically
// indistinguishable from generated codes.
func (s *Shortener) IsVanity(shortcode string) bool {
	return s.IsVanityFunc != nil && s.IsVanityFunc(shortcode)
}

// Sort sorts shorter codes first and generated codes before vanity
// codes.
func (s *Shortener) Sort(shortcodes []string) {
	less := func(a, b string) bool {
		return (len(a) == len(b) && a < b) || len(a) < len(b)
	}
	if s.IsVanityFunc != nil {
		less = func(a, b string) bool {
			aVanity := s.IsVanityFunc(a)
			bVanity := s.IsVanityFunc(b)
			return (aVanity == bVanity && ((len(a) == len(b) && a < b) || len(a) < len(b))) ||
				(!aVanity && bVanity)
		}
	}
	sort.Slice(shortcodes, func(i, j int) bool {
		return less(shortcodes[i], shortcodes[j])
	})
}

// GetIAShortcodes queries all the shortcodes that have been archived on
// the Internet Archive.
func (s *Shortener) GetIAShortcodes() ([]string, error) {
	timemap, err := ia.GetTimemap(s.Host, &ia.TimemapOptions{
		Collapse:    "original",
		Fields:      []string{"original"},
		MatchPrefix: true,
		Limit:       100000,
	})
	if err != nil {
		return nil, err
	}
	urls := make([]string, len(timemap))
	for i, link := range timemap {
		urls[i] = link[0]
	}
	return s.CleanURLs(urls)
}

func splitByte(s string, c byte) (string, string) {
	if i := strings.IndexByte(s, c); i != -1 {
		return s[:i], s[i+1:]
	}
	return s, ""
}

func trimAfterByte(s string, c byte) string {
	if i := strings.IndexByte(s, c); i != -1 {
		return s[:i]
	}
	return s
}
