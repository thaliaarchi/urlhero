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
	"os"
	"strings"
)

// Spec: https://gbv.github.io/beaconspec/beacon.html

type Reader struct {
	r        *bufio.Reader
	meta     []MetaField
	metaRead bool
	line     string

	// If LazyBars is true, all links are resolved as SOURCE|TARGET. Any
	// further '|' characters on a line are considered to be part of
	// TARGET.
	LazyBars bool
}

type MetaField struct {
	Name, Value string
}

type Link struct {
	Source, Target, Annotation string
}

func NewReader(r io.Reader) *Reader {
	return &Reader{r: bufio.NewReader(r)}
}

func (r *Reader) Meta() ([]MetaField, error) {
	meta, err := r.readMeta()
	if err == nil || err == io.EOF {
		return meta, nil
	}
	return nil, err
}

func (r *Reader) readMeta() ([]MetaField, error) {
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
		meta, err := splitMeta(line[1:])
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
	if ch == '\uFEFF' {
		return nil
	}
	return r.r.UnreadRune()
}

func splitMeta(meta string) (MetaField, error) {
	for i, ch := range meta {
		switch {
		case 'A' <= ch && ch <= 'Z':
		case ch == ':' || ch == ' ' || ch == '\t':
			field, value := meta[:i], meta[i+1:]
			for value != "" && (value[0] == ' ' || value[0] == '\t') {
				value = value[1:]
			}
			return MetaField{field, value}, nil
		default:
			return MetaField{}, fmt.Errorf("beacon: invalid character %q in meta field: %q", ch, meta)
		}
	}
	return MetaField{}, fmt.Errorf("beacon: meta line missing value: %q", meta)
}

func (r *Reader) Read() (*Link, error) {
	if !r.metaRead {
		if _, err := r.Meta(); err != nil {
			return nil, err
		}
	}
	line := ""
	var err error
	for line == "" {
		line, err = r.readLine()
		if err != nil {
			return nil, err
		}
	}
	if r.LazyBars {
		if i := strings.IndexByte(line, '|'); i != -1 {
			return &Link{line[:i], line[i+1:], ""}, nil
		}
		fmt.Fprintf(os.Stderr, "beacon: link line missing bar separator: %s\n", line)
		return &Link{line, "", ""}, nil
	}
	var link Link
	tokens := strings.SplitN(line, "|", 4)
	switch len(tokens) {
	case 1:
		link.Source = tokens[0]
	case 2:
		link.Source, link.Target = tokens[0], tokens[1]
		// TODO:
		// link.Source, link.Annotation = tokens[0], tokens[1]
	case 3:
		link.Source, link.Annotation, link.Target = tokens[0], tokens[1], tokens[2]
	case 4:
		return nil, fmt.Errorf("beacon: link line has too many bar separators: %q", line)
	}
	return &link, nil
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

func (m MetaField) String() string {
	return fmt.Sprintf("#%s: %s", m.Name, m.Value)
}

func (l Link) String() string {
	if l.Annotation != "" {
		return fmt.Sprintf("%s|%s|%s", l.Source, l.Annotation, l.Target)
	}
	return fmt.Sprintf("%s|%s", l.Source, l.Target)
}
