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

	"github.com/andrewarchi/urlhero/beacon"
	"github.com/andrewarchi/urlhero/tinytown"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s DIR PATTERN\n", os.Args[0])
		os.Exit(2)
	}
	dir, pattern := os.Args[1], os.Args[2]
	re, err := regexp.Compile(pattern)
	try(err)

	processLink := func(l *beacon.Link, m *tinytown.Meta, shortcodeLen int, releaseFilename, dumpFilename string) error {
		if re.MatchString(l.Target) {
			fmt.Println(l.Target)
		}
		return nil
	}
	try(tinytown.ProcessReleases(dir, processLink))
}

func try(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
