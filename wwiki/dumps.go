// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package wwiki processes dumps for Wikimedia's w.wiki link shortener.
package wwiki

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// DumpInfo contains information on a short URL dump.
type DumpInfo struct {
	URL  *url.URL
	Time time.Time
	Size int64
}

// GetDumps retrieves information on all short URL dumps.
func GetDumps() ([]DumpInfo, error) {
	const indexURL = "https://dumps.wikimedia.org/other/shorturls/"
	baseURL, err := url.Parse(indexURL)
	if err != nil {
		return nil, err
	}

	resp, err := httpGet(indexURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	pre := findFirst(doc, atom.Pre)
	if pre == nil {
		return nil, errors.New("wwiki: no pre element")
	}
	var dumps []DumpInfo
	err = eachChild(pre, atom.A, func(a *html.Node) error {
		href, _ := attr(a, "href")
		if href == "../" {
			return nil
		}
		rel, err := url.Parse(href)
		if err != nil {
			return err
		}
		u := baseURL.ResolveReference(rel)

		if a.NextSibling.Type != html.TextNode {
			return fmt.Errorf("eeiki: no time and size for %s", href)
		}
		text := strings.TrimSpace(a.NextSibling.Data)

		i := strings.LastIndexByte(text, ' ')
		if i == -1 {
			return fmt.Errorf("wwiki: cannot split time and size for %s", href)
		}
		timeStr := strings.TrimSpace(text[:i])
		sizeStr := text[i+1:]
		t, err := time.Parse("02-Jan-2006 15:04", timeStr)
		if err != nil {
			return err
		}
		size, err := strconv.ParseInt(sizeStr, 10, 64)
		if err != nil {
			return err
		}

		dumps = append(dumps, DumpInfo{u, t, size})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return dumps, nil
}

func httpGet(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("wwiki: http status %s", resp.Status)
	}
	return resp, nil
}

func findFirst(n *html.Node, tag atom.Atom) *html.Node {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.DataAtom == tag {
			return c
		}
		if n := findFirst(c, tag); n != nil {
			return n
		}
	}
	return nil
}

func eachChild(n *html.Node, tag atom.Atom, fn func(*html.Node) error) error {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.DataAtom == tag {
			if err := fn(c); err != nil {
				return err
			}
		}
	}
	return nil
}

func attr(n *html.Node, key string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false
}
