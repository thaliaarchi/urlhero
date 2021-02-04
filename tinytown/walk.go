// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tinytown

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/andrewarchi/urlteam/beacon"
	"github.com/ulikunitz/xz"
)

func processAll(root string) error {
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
			name := file.Name()
			if !strings.HasSuffix(name, ".zip") {
				continue
			}
			if _, err := processZip(filepath.Join(dir, name), file.Size()); err != nil {
				return err
			}
		}
	}
	return nil
}

func processZip(filename string, size int64) (*Meta, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r, err := zip.NewReader(f, size)
	if err != nil {
		return nil, err
	}
	var meta *Meta
	for _, zf := range r.File {
		fmt.Printf("%s\n", filepath.Join(filename, zf.Name))
		if !strings.HasSuffix(zf.Name, ".xz") || zf.FileInfo().IsDir() {
			continue
		}

		zr, err := zf.Open()
		if err != nil {
			return nil, err
		}
		defer zr.Close()
		xr, err := xz.NewReader(zr)
		if err != nil {
			return nil, err
		}

		switch {
		case strings.HasSuffix(zf.Name, ".meta.json.xz"):
			jd := json.NewDecoder(xr)
			jd.DisallowUnknownFields()
			var m Meta
			if err := jd.Decode(&m); err != nil {
				return nil, err
			}
			meta = &m
			continue
		case strings.HasSuffix(zf.Name, ".txt.xz"):
			br := beacon.NewReader(xr)
			_, err := br.Header()
			if err != nil {
				return nil, err
			}
			for {
				_, err := br.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					fmt.Fprintln(os.Stderr, filepath.Join(filepath.Base(filename), zf.Name), err)
					continue
				}
			}
		}
	}
	return meta, nil
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
