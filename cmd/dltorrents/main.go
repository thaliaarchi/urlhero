// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/andrewarchi/urlteam/tinytown"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s dir\n", os.Args[0])
		os.Exit(2)
	}
	dir := os.Args[1]
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "No such directory: %s", dir)
		os.Exit(1)
	}

	if err := tinytown.DownloadTorrents(dir); err != nil {
		log.Fatal(err)
	}
}
