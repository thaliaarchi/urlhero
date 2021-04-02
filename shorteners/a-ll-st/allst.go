// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package allst handles Allstates's a.ll.st link shortener.
package allst

import (
	"net/url"
	"sort"
	"strings"

	"github.com/andrewarchi/urlhero/ia"
)

func GetIAShortcodes() ([]string, error) {
	timemap, err := ia.GetTimemap("a.ll.st", &ia.TimemapOptions{
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
		// Remove query parameters:
		//   http://a.ll.st/agentlocatorTW?linkId=88456333
		//   http://a.ll.st/PP45AJ?utm_sourcex3dplus.url.google.comx26utm_mediumx3dreferral
		u, err := url.Parse(link[0])
		if err != nil {
			return nil, err
		}
		shortcode := strings.TrimPrefix(u.Path, "/")
		// Remove trailing JSON for some social media shortcodes:
		//   http://a.ll.st/Facebook","navigationEndpoint
		//   http://a.ll.st/Instagram","isCrawlable":true,"thumbnail
		if i := strings.IndexByte(shortcode, '"'); i != -1 {
			shortcode = shortcode[:i]
		}
		switch shortcode {
		case "", "robots.txt", "favicon.ico":
			continue
		}
		// Remove /scmf/ID/ prefix:
		//   http://a.ll.st/scmf/OrMCe04Lcp0lODk0BD1FrBcO2E4FP0NMEHFGSZ--Pq5q7EdIBj5D0RZwQ0r5O5LJxfQiUmcjxE_yFyVUmcC7Ue52R7KC2DlT6j1Anuut1CVBLh2fal1IZic40eX4xD2dJTg/PrJJpv
		//   http://a.ll.st/scmf/OrMCe04Lcp0lODk2Bzg71hcM2079O8ZJEHE_NJu-wtVr7D9JB0U8qWl1RzYCRZPJxfQiUmcjxE_yF9swgNxdUAkTP4vGed-VJvLu3uityvkzL-5fGDGJnyV0iKf6RXKdJQ/hiddenworldofdata
		if strings.HasPrefix(shortcode, "scmf/") {
			shortcode = shortcode[strings.LastIndexByte(shortcode, '/')+1:]
		}
		if _, ok := shortcodesMap[shortcode]; !ok {
			shortcodesMap[shortcode] = struct{}{}
			shortcodes = append(shortcodes, shortcode)
		}
	}
	sort.Slice(shortcodes, func(i, j int) bool {
		// Sort 6-character generated codes before vanity codes.
		si, sj := shortcodes[i], shortcodes[j]
		return (len(si) == len(sj) && si < sj) ||
			(len(si) == 6 && len(sj) != 6) || len(si) < len(sj)
	})
	return shortcodes, nil
}
