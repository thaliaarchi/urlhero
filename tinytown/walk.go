// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tinytown

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/andrewarchi/archive"
	"github.com/andrewarchi/urlteam/beacon"
)

func ProcessReleases(root string) error {
	rootContents, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}
	for _, release := range rootContents {
		if !release.IsDir() {
			continue
		}
		dir := filepath.Join(root, release.Name())
		dirContents, err := ioutil.ReadDir(dir)
		if err != nil {
			return err
		}
		for _, file := range dirContents {
			filename := filepath.Join(dir, file.Name())
			if !strings.HasSuffix(filename, ".zip") {
				continue
			}
			fmt.Fprintln(os.Stderr, filename)
			err := archive.WalkFile(filename, func(f archive.File) error {
				if strings.HasSuffix(f.Name(), ".meta.json.xz") {
					return nil
				}
				fmt.Fprintf(os.Stderr, "\t%s\n", f.Name())
				r, err := f.Open()
				if err != nil {
					return err
				}
				defer r.Close()
				xr, err := archive.NewXZReader(r)
				if err != nil {
					return err
				}
				defer xr.Close()
				br := beacon.NewReader(xr, beacon.URLTeam)
				for {
					link, err := br.Read()
					if err != nil {
						if err == io.EOF {
							return nil
						}
						return err
					}
					_ = link
				}
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type Meta struct {
	Alphabet          string  `json:"alphabet"`
	Autoqueue         bool    `json:"autoqueue"`
	AutoreleaseTime   int     `json:"autorelease_time"`
	BannedCodes       []int   `json:"banned_codes"`
	BodyRegex         string  `json:"body_regex"`
	Enabled           bool    `json:"enabled"`
	LocationAntiRegex string  `json:"location_anti_regex"`
	LowerSequenceNum  int     `json:"lower_sequence_num"`
	MaxNumItems       int     `json:"max_num_items"`
	Method            string  `json:"method"`
	MinClientVersion  int     `json:"min_client_version"`
	MinVersion        int     `json:"min_version"`
	Name              string  `json:"name"`
	NoRedirectCodes   []int   `json:"no_redirect_codes"`
	NumCountPerItem   int     `json:"num_count_per_item"`
	RedirectCodes     []int   `json:"redirect_codes"`
	RequestDelay      float64 `json:"request_delay"`
	UnavailableCodes  []int   `json:"unavailable_codes"`
	URLTemplate       string  `json:"url_template"`
}
