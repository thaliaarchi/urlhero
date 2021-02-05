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
	"fmt"
	"io"
	"strings"
)

// Spec: https://gbv.github.io/beaconspec/beacon.html

type Reader struct {
	r        *bufio.Reader
	meta     []Meta
	metaRead bool
	line     string
}

type Meta struct {
	Field, Value string
}

type Link struct {
	Shortcode, URL string
}

func NewReader(r io.Reader) *Reader {
	return &Reader{r: bufio.NewReader(r)}
}

func (r *Reader) Meta() ([]Meta, error) {
	meta, err := r.readMeta()
	if err == nil || err == io.EOF {
		return meta, nil
	}
	return nil, err
}

func (r *Reader) readMeta() ([]Meta, error) {
	if r.metaRead {
		return r.meta, nil
	}
	r.metaRead = true
	if err := r.consumeBOM(); err != nil {
		return nil, err
	}
	// Allow omitted header section
	if b, err := r.r.Peek(1); err != nil || b[0] != '#' {
		return nil, err
	}

	// Read meta lines until the first blank line or non-#-prefixed line
	for {
		line, err := r.readLine()
		if err != nil || line == "" {
			return r.meta, err
		}
		if line[0] != '#' {
			r.line = line
			return r.meta, nil
		}
		meta, err := splitMeta(line)
		if err != nil {
			return nil, err
		}
		r.meta = append(r.meta, meta)
	}
}

// consumeBOM skips a UTF-8 byte order mark as permitted by section 3.1.
func (r *Reader) consumeBOM() error {
	ch, _, err := r.r.ReadRune()
	if err != nil {
		return err
	}
	if ch == '\ufeff' {
		return nil
	}
	return r.r.UnreadRune()
}

func splitMeta(meta string) (Meta, error) {
	for i, ch := range meta[1:] {
		switch {
		case 'A' <= ch && ch <= 'Z':
		case ch == ':' || ch == ' ' || ch == '\t':
			field, value := meta[:i], meta[i+1:]
			for value != "" && (value[0] == ' ' || value[0] == '\t') {
				value = value[1:]
			}
			return Meta{field, value}, nil
		default:
			return Meta{}, fmt.Errorf("beacon: invalid character %q in meta field: %q", ch, meta)
		}
	}
	return Meta{}, fmt.Errorf("beacon: meta line missing value: %q", meta)
}

func (r *Reader) Read() (*Link, error) {
	if !r.metaRead {
		if _, err := r.Meta(); err != nil {
			return nil, err
		}
	}
	line, err := r.readLine()
	if err != nil {
		return nil, err
	}
	i := strings.IndexByte(line, '|')
	if i == -1 {
		return nil, fmt.Errorf("beacon: link line missing bar separator: %s", line)
	}
	return &Link{line[:i], line[i+1:]}, nil
}

func (r *Reader) readLine() (string, error) {
	if l := r.line; l != "" {
		r.line = ""
		return l, nil
	}
	line, err := r.r.ReadString('\n')
	if err != nil {
		return "", err
	}
	if len(line) >= 2 && line[len(line)-2] == '\r' {
		return line[:len(line)-2], nil
	}
	return line[:len(line)-1], nil
}

func (m Meta) String() string {
	return fmt.Sprintf("#%s: %s", m.Field, m.Value)
}

func (l Link) String() string {
	return fmt.Sprintf("%s|%s", l.Shortcode, l.URL)
}
