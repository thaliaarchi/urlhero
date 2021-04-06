// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package shorteners provides utilities for retrieving information
// about URL shortening websites.
package shorteners

import (
	"bytes"
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
	MobyTo,
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
	return s.CleanURL(u)
}

// CleanURL extracts the shortcode from a URL. An empty string is
// returned when no shortcode can be found.
func (s *Shortener) CleanURL(u *url.URL) (string, error) {
	shortcode := cleanURL(u, s.CleanFunc)
	if shortcode != "" && s.Pattern != nil && !s.Pattern.MatchString(shortcode) {
		return "", fmt.Errorf("%s: shortcode %q does not match alphabet %s after cleaning: %q", s.Name, shortcode, s.Pattern, u)
	}
	return shortcode, nil
}

func cleanURL(u *url.URL, clean CleanFunc) string {
	shortcode := strings.TrimLeft(u.Path, "/")
	// Exclude placeholders like <key>
	if len(shortcode) >= 2 {
		s0, s1 := shortcode[0], shortcode[len(shortcode)-1]
		if (s0 == '<' && s1 == '>') || (s0 == '[' && s1 == ']') {
			return ""
		}
	}
	// Remove trailing junk (escapes are nbsp and zwsp)
	if i := strings.IndexAny(shortcode, "\"])>&’” \u00a0\u200B"); i != -1 {
		shortcode = shortcode[:i]
	}
	// Remove trailing punctuation
	shortcode = strings.TrimRight(shortcode, ".;")
	shortcode = strings.TrimRight(shortcode, "/")
	// Remove concatenated URLs
	shortcode = trimAfter(shortcode, "http:/")
	shortcode = trimAfter(shortcode, "https:/")
	if isCommonFile(shortcode) {
		return ""
	}
	if clean != nil {
		shortcode = clean(shortcode, u)
	}
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
	var errs []error
	for _, shortURL := range urls {
		shortcode, err := s.Clean(shortURL)
		if err != nil {
			errs = append(errs, err)
			continue
		} else if shortcode == "" {
			continue
		}
		if _, ok := shortcodesMap[shortcode]; !ok {
			shortcodesMap[shortcode] = struct{}{}
			shortcodes = append(shortcodes, shortcode)
		}
	}
	s.Sort(shortcodes)
	if len(errs) != 0 {
		return shortcodes, &multiError{"CleanURLs", errs}
	}
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

func trimAfter(s string, substr string) string {
	if i := strings.Index(s, substr); i != -1 {
		return s[:i]
	}
	return s
}

func trimAfterByte(s string, c byte) string {
	if i := strings.IndexByte(s, c); i != -1 {
		return s[:i]
	}
	return s
}

type multiError struct {
	tag  string
	errs []error
}

func (merr *multiError) Error() string {
	if len(merr.errs) == 0 {
		return merr.tag
	}
	if len(merr.errs) == 1 {
		return fmt.Sprintf("%s: %s", merr.tag, merr.errs[0])
	}
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s:\n", merr.tag)
	for _, err := range merr.errs {
		fmt.Fprintf(&b, "\t%s\n", err)
	}
	return b.String()
}
