package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/andrewarchi/urlteam/beacon"
	"github.com/ulikunitz/xz"
)

func main() {
	dir := os.Args[1]
	if err := processAll(dir); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func processAll(root string) error {
	contents, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}
	for _, fi := range contents {
		if !fi.IsDir() {
			continue
		}
		if err := process(filepath.Join(root, fi.Name())); err != nil {
			return err
		}
	}
	return nil
}

func process(dirname string) error {
	dir, err := os.Open(dirname)
	if err != nil {
		return err
	}
	defer dir.Close()
	contents, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	for _, fi := range contents {
		name := fi.Name()
		fmt.Println("check", fi.Name())
		if !strings.HasSuffix(name, ".zip") {
			continue
		}
		f, err := os.Open(filepath.Join(dirname, name))
		if err != nil {
			return err
		}
		defer f.Close()

		f.Stat()
		r, err := zip.NewReader(f, fi.Size())
		if err != nil {
			return err
		}
		var siteMeta SiteMeta
		for _, zf := range r.File {
			fmt.Println("  ", zf.Name)
			if !strings.HasSuffix(zf.Name, ".xz") {
				fmt.Printf("non-xz skipped: %s\n", zf.Name)
				continue
			}
			zr, err := zf.Open()
			if err != nil {
				return err
			}
			xr, err := xz.NewReader(zr)
			if err != nil {
				return err
			}

			switch {
			case strings.HasSuffix(zf.Name, ".meta.json.xz"):
				jd := json.NewDecoder(xr)
				jd.DisallowUnknownFields()
				if err := jd.Decode(&siteMeta); err != nil {
					return err
				}
				continue
			case strings.HasSuffix(zf.Name, ".txt.xz"):
				br := beacon.NewReader(xr)
				meta, err := br.ReadMeta()
				if err != nil {
					return err
				}
				for _, m := range meta {
					fmt.Println(m)
				}
				for {
					link, err := br.Read()
					if err == io.EOF {
						break
					}
					if err != nil {
						return err
					}
					fmt.Println(link)
				}
			}
		}
	}
	return nil
}

type SiteMeta struct {
	Alphabet          string  `json:"alphabet"`
	Autoqueue         bool    `json:"autoqueue"`
	AutoreleaseTime   int     `json:"autorelease_time"`
	BannedCodes       []int   `json:"banned_codes"`
	BodyRegex         string  `json:"body_regex"`
	Enabled           bool    `json:"enabled"`
	LocationAntiRegex string  `json:"location_anti_regex"`
	LowerSequenceNum  int     `json:"lower_sequence_num"`
	MaxNumItems       int     `json:"max_num_items"`
	Method            string  `json:"method"`
	MinClientVersion  int     `json:"min_client_version"`
	MinVersion        int     `json:"min_version"`
	Name              string  `json:"name"`
	NoRedirectCodes   []int   `json:"no_redirect_codes"`
	NumCountPerItem   int     `json:"num_count_per_item"`
	RedirectCodes     []int   `json:"redirect_codes"`
	RequestDelay      float64 `json:"request_delay"`
	UnavailableCodes  []int   `json:"unavailable_codes"`
	URLTemplate       string  `json:"url_template"`
}
