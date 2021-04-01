// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ia

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// TimemapOptions contains options for a timemap API call.
type TimemapOptions struct {
	MatchPrefix bool     // whether url is a prefix (* wildcard is appended)
	Collapse    string   // field to collapse by; earliest captures with unique field is kept
	Fields      []string // i.e. urlkey,timestamp,endtimestamp,original,mimetype,statuscode,digest,redirect,robotflags,length,offset,filename,groupcount,uniqcount
	Limit       int      // i.e. 100000
}

// GetTimemap gets a list of Internet Archive captures of the given URL.
func GetTimemap(pageURL string, options *TimemapOptions) ([][]string, error) {
	// Timemap API, as observed on
	// https://web.archive.org/web/*/https://dumps.wikimedia.org/other/shorturls/*

	q := make(url.Values)
	q.Set("url", pageURL)
	q.Set("output", "json") // other values: "csv" and omitted
	if options.MatchPrefix {
		q.Set("matchType", "prefix") // other values unknown
	}
	if options.Collapse != "" {
		q.Set("collapse", options.Collapse)
	}
	if len(options.Fields) != 0 {
		q.Set("fl", strings.Join(options.Fields, ","))
	}
	if options.Limit > 0 {
		q.Set("limit", strconv.Itoa(options.Limit))
	}

	resp, err := http.Get("https://web.archive.org/web/timemap/?" + q.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ia: http status %s", resp.Status)
	}

	var timemap [][]string
	if err := json.NewDecoder(resp.Body).Decode(&timemap); err != nil {
		return nil, err
	}
	return timemap, nil
}
