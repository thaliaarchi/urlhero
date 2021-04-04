// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package shorteners

import (
	"fmt"
	"sort"

	"github.com/andrewarchi/urlhero/ia"
)

// GetIAShortcodes queries all the shortcodes that have been archived on
// the Internet Archive. If alpha, clean, or less are nil, defaults will be
// used.
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
	shortcodesMap := make(map[string]struct{})
	var shortcodes []string
	for _, link := range timemap {
		shortcode, err := s.Clean(link[0])
		if err != nil {
			return nil, err
		} else if shortcode == "" {
			continue
		}
		if s.Pattern != nil && !s.Pattern.MatchString(shortcode) {
			return nil, fmt.Errorf("%s: shortcode %q does not match alphabet %s after cleaning: %q", s.Name, shortcode, s.Pattern, link[0])
		}
		if _, ok := shortcodesMap[shortcode]; !ok {
			shortcodesMap[shortcode] = struct{}{}
			shortcodes = append(shortcodes, shortcode)
		}
	}
	less := s.LessFunc
	if less == nil {
		less = func(a, b string) bool {
			return (len(a) == len(b) && a < b) || len(a) < len(b)
		}
	}
	sort.Slice(shortcodes, func(i, j int) bool {
		return less(shortcodes[i], shortcodes[j])
	})
	return shortcodes, nil
}
