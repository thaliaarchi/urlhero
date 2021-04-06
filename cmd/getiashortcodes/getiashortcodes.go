// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/andrewarchi/urlhero/shorteners"
)

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, "Usage: getiashortcodes <shortener> [alphabet]")
		os.Exit(2)
	}
	shortener := os.Args[1]
	var alpha string
	if len(os.Args) >= 3 {
		alpha = os.Args[2]
	}

	s, ok := shorteners.Lookup[shortener]
	if !ok {
		var name, host string
		if strings.ContainsRune(shortener, '.') {
			name, host = strings.ReplaceAll(shortener, ".", "-"), shortener
		} else {
			name, host = shortener, strings.ReplaceAll(shortener, "-", ".")
		}
		s = &shorteners.Shortener{
			Name:     name,
			Host:     host,
			Alphabet: alpha,
		}
		if alpha != "" {
			pattern, err := regexp.Compile("^[" + alpha + "]+$")
			try(err)
			s.Pattern = pattern
		}
	}

	shortcodes, err := s.GetIAShortcodes()
	try(err)
	for _, shortcode := range shortcodes {
		fmt.Println(shortcode)
	}
}

func try(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
