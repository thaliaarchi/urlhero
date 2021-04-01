// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package wwiki

import (
	"fmt"
	"io"
	"net/url"
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
		out := filepath.Join(dir, path.Base(dump.URL.Path))
		if err := downloadDump(dump.URL.String(), out, nil); err != nil {
			return err
		}
	}
	return nil
}

// DownloadIADumps saves all short URL dumps that have been archived by
// the Internet Archive to the given directory.
func DownloadIADumps(dir string) error {
	dumps, err := GetIADumps()
	if err != nil {
		return err
	}
	for _, dump := range dumps {
		iaURL := fmt.Sprintf("https://web.archive.org/web/%sif_/%s", dump.Timestamp, dump.URL)
		u, err := url.Parse(dump.URL)
		if err != nil {
			return err
		}
		out := filepath.Join(dir, path.Base(u.Path))
		if err := downloadDump(iaURL, out, dump.SHA1[:]); err != nil {
			return err
		}
	}
	return nil
}

func downloadDump(url, out string, sha1Sum []byte) error {
	fmt.Println("Downloading", url)
	// Skip existing
	if _, err := os.Stat(out); err == nil {
		if sha1Sum != nil {
			return ia.ValidateFile(out, nil, sha1Sum, nil)
		}
		// TODO check ETag, if it is a checksum
		return nil
	}

	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()

	resp, err := httpGet(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r io.Reader = resp.Body
	if sha1Sum != nil {
		r = ia.NewReadValidator(r, url, nil, sha1Sum, nil)
	}
	if _, err := io.Copy(f, r); err != nil {
		return err
	}

	mod := resp.Header.Get("Last-Modified")
	if mod == "" {
		mod = resp.Header.Get("X-Archive-Orig-Last-Modified")
	}
	if mod != "" {
		mt, err := time.Parse(time.RFC1123, mod)
		if err != nil {
			return err
		}
		return os.Chtimes(out, mt, mt)
	}
	return nil
}

// IADumpInfo contains basic Internet Archive metadata on a short URL
// dump.
type IADumpInfo struct {
	URL       string
	Timestamp string
	SHA1      [20]byte
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
		// Exclude the index file and include early non-gzipped dumps
		if statuscode == "200" && (mimetype == "application/octet-stream" || mimetype == "text/plain") {
			sha1, err := ia.DecodeDigest(digest)
			if err != nil {
				return nil, err
			}
			dumps = append(dumps, IADumpInfo{original, timestamp, *sha1})
		}
	}
	return dumps, nil
}
