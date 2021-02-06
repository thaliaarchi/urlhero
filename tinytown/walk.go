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

	"github.com/andrewarchi/archive"
	"github.com/andrewarchi/urlteam/beacon"
)

// Meta contains link dump metadata from a *.meta.json.xz file.
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

// ProcessFunc is the type of function that is called for each link
// visited.
type ProcessFunc func(link *beacon.Link, meta *Meta, shortcodeLen int) error

// ProcessReleases processes every release in a directory by calling fn
// on every link.
func ProcessReleases(root string, fn ProcessFunc) error {
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
			if err := ProcessRelease(filename, fn); err != nil {
				return err
			}
		}
	}
	return nil
}

// ProcessRelease processes every link dump in a release by calling fn
// on every link.
func ProcessRelease(filename string, fn ProcessFunc) error {
	zr, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer zr.Close()
	if len(zr.File) == 0 {
		return fmt.Errorf("tinytown: empty archive: %q", filename)
	}

	meta, err := readMeta(zr.File[0])
	if err != nil {
		return err
	}
	for _, f := range zr.File[1:] {
		fmt.Fprintf(os.Stderr, "\t%s\n", f.Name)
		if err := processLinkDump(f, meta, fn); err != nil {
			return err
		}
	}
	return nil
}

func readMeta(f *zip.File) (*Meta, error) {
	if !strings.HasSuffix(f.Name, ".meta.json.xz") {
		return nil, fmt.Errorf("tinytown: not a meta file: %q", f.Name)
	}
	fr, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer fr.Close()
	xr, err := archive.NewXZReader(fr)
	if err != nil {
		return nil, err
	}
	defer xr.Close()
	var m Meta
	if err := json.NewDecoder(xr).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func processLinkDump(f *zip.File, meta *Meta, fn ProcessFunc) error {
	if !strings.HasSuffix(f.Name, ".txt.xz") {
		return fmt.Errorf("tinytown: not a link dump: %q", f.Name)
	}
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

	shortcodeLen := len(f.Name) - len(".txt.xz")
	br := beacon.NewURLTeamReader(xr, shortcodeLen)
	for {
		link, err := br.Read()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if err := fn(link, meta, shortcodeLen); err != nil {
			return err
		}
	}
}
