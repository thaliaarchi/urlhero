// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package beacon processes BEACON-format link dumps as defined by GBV
// and used by URLTeam.
package beacon

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

// Spec: https://gbv.github.io/beaconspec/beacon.html

type Reader struct {
	s          *bufio.Scanner
	header     []Header
	readHeader bool
	line       string
}

type Header struct {
	Field, Value string
}

type Link struct {
	Shortcode, URL string
}

func NewReader(r io.Reader) *Reader {
	s := bufio.NewScanner(r)
	s.Split(scanLines)
	return &Reader{s: s}
}

func (r *Reader) Header() ([]Header, error) {
	if r.readHeader {
		return r.header, nil
	}
	r.readHeader = true
	if !r.s.Scan() {
		return nil, r.s.Err()
	}
	line := r.s.Text()
	if !strings.HasPrefix(line, "#") { // Allow omitted header section
		r.line = line
		return nil, nil
	}
	for {
		i := strings.IndexByte(line, ':')
		if i == -1 {
			return r.header, fmt.Errorf("header line missing colon separator: %s", line)
		}
		value := strings.TrimLeft(line[i+1:], "\t ")
		r.header = append(r.header, Header{line[1:i], value})
		if !r.s.Scan() {
			return r.header, r.s.Err()
		}
		line = r.s.Text()
		if line == "" {
			return r.header, nil
		}
		if !strings.HasPrefix(line, "#") {
			r.line = line
			return nil, fmt.Errorf("blank line must follow header")
		}
	}
}

func (r *Reader) Read() (*Link, error) {
	if !r.readHeader {
		if _, err := r.Header(); err != nil {
			return nil, err
		}
	}
	line := r.line
	r.line = ""
	if line == "" {
		if !r.s.Scan() {
			if err := r.s.Err(); err != nil {
				return nil, err
			}
			return nil, io.EOF
		}
		line = r.s.Text()
	}
	i := strings.IndexByte(line, '|')
	if i == -1 {
		return nil, fmt.Errorf("link missing bar separator: %s", line)
	}
	return &Link{line[:i], line[i+1:]}, nil
}

// scanLines splits by LF, CRLF, or CR.
func scanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexAny(data, "\n\r"); i >= 0 {
		// We have a full newline-terminated line.
		token = data[0:i]
		if data[i] == '\r' && i+1 < len(data) && data[i+1] == '\n' {
			i++
		}
		return i + 1, token, nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

func (m Header) String() string {
	return fmt.Sprintf("#%s: %s", m.Field, m.Value)
}

func (l Link) String() string {
	return fmt.Sprintf("%s|%s", l.Shortcode, l.URL)
}
