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

type Reader struct {
	r         *bufio.Reader
	meta      []MetaField
	metaRead  bool
	line      string
	format    Format
	sourceLen int
}

type MetaField struct {
	Name, Value string
}

type Link struct {
	Source, Target, Annotation string
}

// Format defines the format of the BEACON link dump.
type Format uint8

const (
	// RFC link dumps follow draft-003 of the BEACON format RFC submitted
	// December 2017 at https://gbv.github.io/beaconspec/beacon.html.
	RFC Format = iota

	// URLTeam link dumps resolve all links as SOURCE|TARGET. Any further
	// '|' characters on a line are considered to be part of TARGET.
	URLTeam
)

func NewReader(r io.Reader) *Reader {
	return &Reader{r: bufio.NewReader(r)}
}

func NewURLTeamReader(r io.Reader, shortcodeLen int) *Reader {
	return &Reader{r: bufio.NewReader(r), format: URLTeam, sourceLen: shortcodeLen}
}

func (r *Reader) Meta() ([]MetaField, error) {
	if r.metaRead {
		return r.meta, nil
	}
	r.metaRead = true
	meta, err := r.readMeta()
	if err == nil || err == io.EOF {
		return meta, nil
	}
	return nil, err
}

func (r *Reader) readMeta() ([]MetaField, error) {
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
		if err != nil {
			return r.meta, err
		}
		if trimLeftSpace(line) == "" {
			break
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

	// Consume empty lines
	for {
		line, err := r.readLine()
		if err != nil {
			return r.meta, err
		}
		if trimLeftSpace(line) != "" {
			r.line = line
			return r.meta, nil
		}
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
			return MetaField{meta[:i], trimLeftSpace(meta[i+1:])}, nil
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

	line, err := r.readLine()
	if err != nil {
		return nil, err
	}

	if r.format == URLTeam {
		if i := strings.IndexByte(line, '|'); i != -1 {
			shortcode, target := line[:i], line[i+1:]
			if len(shortcode) != r.sourceLen {
				fmt.Fprintf(os.Stderr, "beacon: shortcode not %d characters: %q\n", r.sourceLen, line)
			}
			return &Link{shortcode, target, ""}, nil
		}
		fmt.Fprintf(os.Stderr, "beacon: link line missing bar separator: %q\n", line)
		return &Link{"", line, ""}, nil
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
	if err != nil && !(err == io.EOF && line != "") {
		return "", err
	}
	return dropLineBreak(line), nil
}

func dropLineBreak(line string) string {
	if len(line) >= 1 && line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
		if len(line) >= 1 && line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}
	}
	return line
}

func trimLeftSpace(s string) string {
	for s != "" && (s[0] == ' ' || s[0] == '\t') {
		s = s[1:]
	}
	return s
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
