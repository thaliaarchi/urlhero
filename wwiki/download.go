// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package wwiki

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/andrewarchi/urlteam/ia"
)

// DownloadDumps saves all short URL dumps to the given directory.
func DownloadDumps(dir string) error {
	dumps, err := GetDumps()
	if err != nil {
		return err
	}
	for _, dump := range dumps {
		if err := downloadDump(dump, dir); err != nil {
			return err
		}
	}
	return nil
}

func downloadDump(dump DumpInfo, dir string) error {
	name := filepath.Join(dir, path.Base(dump.URL.Path))
	if _, err := os.Stat(name); err == nil {
		// Skip existing
		return nil
	}
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	resp, err := httpGet(dump.URL.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return err
	}

	if mod := resp.Header.Get("Last-Modified"); mod != "" {
		mt, err := time.Parse(time.RFC1123, mod)
		if err != nil {
			return err
		}
		return os.Chtimes(name, mt, mt)
	}
	return nil
}

// IADumpInfo contains basic Internet Archive metadata on a short URL
// dump.
type IADumpInfo struct {
	URL       string
	Timestamp string
	Digest    string
}

// GetIADumps retrieves information on all short URL dumps that have
// been archived by the Internet Archive.
func GetIADumps() ([]IADumpInfo, error) {
	timemap, err := ia.GetTimemap("https://dumps.wikimedia.org/other/shorturls/", &ia.TimemapOptions{
		MatchPrefix: true,
		Collapse:    "digest",
		Fields:      []string{"original", "timestamp", "mimetype", "statuscode", "digest"},
		Limit:       100000,
	})
	if err != nil {
		return nil, err
	}

	dumps := make([]IADumpInfo, 0, len(timemap)-1) // Skip header row
	for _, d := range timemap[1:] {
		original, timestamp, mimetype, statuscode, digest := d[0], d[1], d[2], d[3], d[4]
		// Exclude the index file and include early non-gzipped dumps.
		if statuscode == "200" && (mimetype == "application/octet-stream" || mimetype == "text/plain") {
			dumps = append(dumps, IADumpInfo{original, timestamp, digest})
		}
	}
	return dumps, nil
}
