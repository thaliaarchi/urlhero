// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tinytown

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/andrewarchi/urlhero/beacon"
)

func SearchReleases(root, shortener string, shortcodes []string) ([]*beacon.Link, error) {
	// TODO merge with ProcessReleases.
	shortcodeMap := make(map[string]struct{})
	for _, shortcode := range shortcodes {
		shortcodeMap[shortcode] = struct{}{}
	}

	rootContents, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}
	var links []*beacon.Link
	for _, release := range rootContents {
		if !release.IsDir() {
			continue
		}
		dir := filepath.Join(root, release.Name())
		dirContents, err := ioutil.ReadDir(dir)
		if err != nil {
			return nil, err
		}
		for _, file := range dirContents {
			name := file.Name()
			if !strings.HasPrefix(name, shortener+".") || !strings.HasSuffix(name, ".zip") {
				continue
			}
			filename := filepath.Join(dir, name)
			// TODO only search link dumps with shortcode length in the set of
			// lengths being searched for.
			fn := func(l *beacon.Link, m *Meta, shortcodeLen int, releaseFilename, dumpFilename string) error {
				if _, ok := shortcodeMap[l.Source]; ok {
					fmt.Println(l)
					links = append(links, l)
				}
				return nil
			}
			if err := ProcessProject(filename, fn); err != nil {
				return nil, err
			}
		}
	}
	return links, nil
}
