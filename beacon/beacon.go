package beacon

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

// https://gbv.github.io/beaconspec/beacon.html

type Reader struct {
	s        *bufio.Scanner
	readMeta bool
	meta     []Meta
	line     string
	lines    int
}

type Meta struct {
	Field, Value string
}

type Link struct {
	Source, Target, Annotation string
}

func NewReader(r io.Reader) *Reader {
	s := bufio.NewScanner(r)
	s.Split(scanLines)
	return &Reader{s: s}
}

func (r *Reader) ReadMeta() ([]Meta, error) {
	if r.readMeta {
		return nil, errors.New("meta already read")
	}
	r.readMeta = true
	if !r.scan() {
		return nil, nil
	}
	line := r.text()
	if !strings.HasPrefix(line, "#") { // Allow omitted meta section
		r.line = line
		return nil, nil
	}
	for {
		i := strings.IndexByte(line, ':')
		if i == -1 {
			return r.meta, fmt.Errorf("meta line missing separator: %s", line)
		}
		value := strings.TrimLeft(line[i+1:], "\t ")
		r.meta = append(r.meta, Meta{line[1:i], value})
		if !r.scan() {
			return r.meta, nil
		}
		line = r.text()
		if line == "" {
			return r.meta, nil
		}
		if !strings.HasPrefix(line, "#") {
			r.line = line
			return nil, fmt.Errorf("blank line must follow meta section")
		}
	}
}

func (r *Reader) Read() (*Link, error) {
	if !r.readMeta {
		if _, err := r.ReadMeta(); err != nil {
			return nil, err
		}
	}
	if !r.scan() {
		return nil, r.err()
	}
	line := r.text()
	fields := strings.Split(line, "|")
	switch len(fields) {
	case 3: // source|annotation|target
		return &Link{fields[0], fields[2], fields[1]}, nil
	case 2: // source|target
		return &Link{fields[0], "", fields[1]}, nil
	case 1: // source
		return &Link{fields[0], "", fields[0]}, nil
	}
	return nil, fmt.Errorf("link line has too many bars: %s", line)
}

func (r *Reader) scan() bool {
	if r.line != "" {
		return true
	}
	r.lines++
	return r.s.Scan()
}

func (r *Reader) text() string {
	if l := r.line; l != "" {
		r.line = ""
		return l
	}
	return r.s.Text()
}

func (r *Reader) err() error {
	if err := r.s.Err(); err != nil {
		return err
	}
	return io.EOF
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
