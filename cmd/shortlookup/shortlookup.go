// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"os"

	"github.com/andrewarchi/urlhero/tinytown"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s dir shortener shortcodes...\n", os.Args[0])
		os.Exit(2)
	}
	dir, shortener, shortcodes := os.Args[1], os.Args[2], os.Args[3:]
	if _, err := tinytown.SearchReleases(dir, shortener, shortcodes); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
