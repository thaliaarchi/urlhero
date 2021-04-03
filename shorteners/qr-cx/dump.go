// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package qrcx

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func DownloadDump(dir string) error {
	url := "https://web.archive.org/web/20151229075230id_/http://qr.cx/dataset/qrcx_all_06eec9b9-1f29-4860-bd91-49c2d517d87d.7z"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	f, err := os.Create(filepath.Join(dir, filepath.Base(url)))
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

type Link struct {
	ShortURL, URL string
	CreationDate  time.Time
}

type Reader struct {
	cr *csv.Reader
}

type ReadCloser struct {
	*Reader
	rc io.ReadCloser
}

func NewReader(r io.Reader) *Reader {
	cr := csv.NewReader(r)
	cr.Comma = '\t'
	cr.Comment = '#'
	cr.FieldsPerRecord = 3
	cr.LazyQuotes = true
	cr.ReuseRecord = true
	return &Reader{cr}
}

func OpenDump(filename string) (*ReadCloser, error) {
	if !strings.HasSuffix(filename, ".csv") {
		return nil, fmt.Errorf("qr-cx: dump is not a CSV file: %s", filename)
	}
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return &ReadCloser{NewReader(f), f}, nil
}

func (r *Reader) Read() (*Link, error) {
	record, err := r.cr.Read()
	if err != nil {
		return nil, err
	}
	t, err := time.Parse("2006-01-02 15:04:05", record[2])
	if err != nil {
		return nil, err
	}
	return &Link{record[0], record[1], t}, nil
}

func (r *Reader) ReadAll() ([]Link, error) {
	var links []Link
	for {
		link, err := r.Read()
		if err == io.EOF {
			return links, nil
		}
		if err != nil {
			return nil, err
		}
		links = append(links, *link)
	}
}

func (r *ReadCloser) Close() error {
	return r.rc.Close()
}
