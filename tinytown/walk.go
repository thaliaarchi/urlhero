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
	Name              string  `json:"name"`
	MinVersion        int     `json:"min_version"`
	MinClientVersion  int     `json:"min_client_version"`
	Alphabet          string  `json:"alphabet"`
	URLTemplate       string  `json:"url_template"`
	RequestDelay      float64 `json:"request_delay"`
	RedirectCodes     []int   `json:"redirect_codes"`    // HTTP codes
	NoRedirectCodes   []int   `json:"no_redirect_codes"` // HTTP codes
	UnavailableCodes  []int   `json:"unavailable_codes"` // HTTP codes
	BannedCodes       []int   `json:"banned_codes"`      // HTTP codes
	BodyRegex         string  `json:"body_regex"`
	LocationAntiRegex string  `json:"location_anti_regex"`
	Method            string  `json:"method"`
	Enabled           bool    `json:"enabled"`
	Autoqueue         bool    `json:"autoqueue"`
	NumCountPerItem   int     `json:"num_count_per_item"`
	MaxNumItems       int     `json:"max_num_items"`
	LowerSequenceNum  int64   `json:"lower_sequence_num"`
	AutoreleaseTime   int     `json:"autorelease_time"`
}

// ProcessFunc is the type of function that is called for each link
// visited.
type ProcessFunc func(l *beacon.Link, m *Meta, shortcodeLen int) error

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
	d := json.NewDecoder(xr)
	d.DisallowUnknownFields()
	if err := d.Decode(&m); err != nil {
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

	shortcodeLen := len(filepath.Base(f.Name)) - len(".txt.xz")
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
