// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"os"

	wwiki "github.com/andrewarchi/urlteam/shorteners/w-wiki"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s dir\n", os.Args[0])
		os.Exit(2)
	}
	dir := os.Args[1]
	try(wwiki.DownloadDumps(dir))
	try(wwiki.DownloadIADumps(dir))
}

func try(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
