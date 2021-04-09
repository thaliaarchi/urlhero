// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ia

import (
	"io"
	"net/http"
	"net/url"
)

type SaveOptions struct {
	CaptureOutlinks    bool
	CaptureAll         bool // save error pages (HTTP status 400-599)
	CaptureScreenshot  bool
	SaveInMyWebArchive bool
	EmailResult        bool
}

func Save(pageURL string, options *SaveOptions) error {
	// Save API, as observed on https://web.archive.org/save

	v := make(url.Values)
	v.Set("url", "https://dumps.wikimedia.org/other/shorturls/shorturls-20210329.gz")
	if options != nil {
		setBool(v, "capture_outlinks", options.CaptureOutlinks)
		setBool(v, "capture_all", options.CaptureAll)
		setBool(v, "capture_screenshot", options.CaptureScreenshot)
		setBool(v, "wm-save-mywebarchive", options.SaveInMyWebArchive)
		setBool(v, "email_result", options.EmailResult)
	}

	resp, err := checkResponse(http.PostForm("https://web.archive.org/save", v))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Ignore HTML body
	_, err = io.Copy(io.Discard, resp.Body)
	return err
}

func setBool(v url.Values, key string, b bool) {
	if b {
		v.Set(key, "on")
	}
}
