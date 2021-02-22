// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package tinyback processes URLTeam's first generation TinyBack
// releases.
package tinyback

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/andrewarchi/archive"
	"github.com/andrewarchi/urlteam/beacon"
	"github.com/andrewarchi/urlteam/ia"
)

// Scraper: https://github.com/ArchiveTeam/tinyback
// Tracker and db: https://github.com/ArchiveTeam/tinyarchive
// Releases and tools: https://github.com/ArchiveTeam/urlteam-stuff

func ProcessRelease(dir string) error {
	metaName := filepath.Base(dir) + "_files.xml"
	files, err := ia.ReadFileMeta(filepath.Join(dir, metaName))
	if err != nil {
		return err
	}
	for i := range files {
		file := &files[i]
		if file.Name == metaName {
			continue // checksums of itself are inaccurate
		}
		if strings.HasSuffix(file.Name, ".txt.xz") { // TODO validate other files
			if err := processFile(file, dir); err != nil {
				return err
			}
		}
	}
	return nil
}

func processFile(file *ia.FileMeta, dir string) error {
	fmt.Print(file.Name)
	fv, err := file.OpenValidator(dir)
	if err != nil {
		return err
	}
	defer fv.Close()
	xr, err := archive.NewXZReader(fv)
	if err != nil {
		return err
	}
	defer xr.Close()
	br := beacon.NewURLTeamReader(xr, -1)
	n := 0
	for {
		link, err := br.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		n++
		_ = link
	}
	fmt.Printf(" [%d links]\n", n)
	return nil
}
