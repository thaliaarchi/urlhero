// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package shorteners

import (
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"github.com/andrewarchi/urlhero/ia"
)

// GetIAShortcodes queries all the shortcodes that have been archived on
// the Internet Archive.
func GetIAShortcodes(shortener string, alpha *regexp.Regexp, clean func(shortcode string) string, less func(i, j string) bool) ([]string, error) {
	timemap, err := ia.GetTimemap(shortener, &ia.TimemapOptions{
		Collapse:    "original",
		Fields:      []string{"original"},
		MatchPrefix: true,
		Limit:       100000,
	})
	if err != nil {
		return nil, err
	}
	shortcodesMap := make(map[string]struct{})
	var shortcodes []string
	for _, link := range timemap {
		u, err := url.Parse(link[0])
		if err != nil {
			return nil, err
		}
		shortcode := clean(strings.TrimPrefix(u.Path, "/"))
		switch shortcode {
		case "", "favicon.ico", "robots.txt":
			continue
		}
		if !alpha.MatchString(shortcode) {
			return nil, fmt.Errorf("shorteners: %s shortcode does not match alphabet %s after cleaning: %s", shortener, alpha, shortcode)
		}
		if _, ok := shortcodesMap[shortcode]; !ok {
			shortcodesMap[shortcode] = struct{}{}
			shortcodes = append(shortcodes, shortcode)
		}
	}
	sort.Slice(shortcodes, func(i, j int) bool {
		return less(shortcodes[i], shortcodes[j])
	})
	return shortcodes, nil
}
