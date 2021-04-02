// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package qrcx

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func DownloadExport(dir string) error {
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
	cr := csv.NewReader(newCommentReader(r))
	cr.Comma = '\t'
	cr.FieldsPerRecord = 3
	cr.LazyQuotes = true
	cr.ReuseRecord = true
	return &Reader{cr}
}

func OpenExport(filename string) (*ReadCloser, error) {
	if !strings.HasSuffix(filename, ".csv") {
		return nil, fmt.Errorf("qr-cx: export is not a CSV file: %s", filename)
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

// commentReader skips header lines prefixed with '#'.
type commentReader struct {
	br    *bufio.Reader
	buf   []byte
	start bool
}

func newCommentReader(r io.Reader) *commentReader {
	return &commentReader{bufio.NewReader(r), nil, true}
}

func (cr *commentReader) Read(p []byte) (n int, err error) {
	if cr.start {
		for {
			var line []byte
			line, err = cr.br.ReadSlice('\n')
			if len(line) == 0 || line[0] != '#' {
				if err == bufio.ErrBufferFull {
					err = nil
				}
				cr.buf = line
				cr.start = false
				break
			}
			for err == bufio.ErrBufferFull {
				line, err = cr.br.ReadSlice('\n')
			}
			if err != nil {
				return 0, err
			}
		}
	}
	if len(cr.buf) != 0 {
		n = copy(p, cr.buf)
		cr.buf = cr.buf[n:]
		return
	}
	if err != nil {
		return
	}
	return cr.br.Read(p)
}
