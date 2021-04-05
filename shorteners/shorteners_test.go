// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package shorteners

import "testing"

func TestIAGetShortcodes(t *testing.T) {
	t.Skip()
	for _, s := range Shorteners {
		shortcodes, err := s.GetIAShortcodes()
		if err != nil {
			t.Errorf("%s: %v", s.Name, err)
		} else if len(shortcodes) == 0 {
			t.Errorf("%s: no shortcodes", s.Name)
		}
	}
}

func TestClean(t *testing.T) {
	tests := []struct {
		s              *Shortener
		url, shortcode string
	}{
		{GoHawaiiEdu, "http://go.hawaii.edu/37c", "37c"},
		{GoHawaiiEdu, "https://go.hawaii.edu/34A", "34A"},
		{GoHawaiiEdu, "http://go.hawaii.edu/Vf+", "Vf"}, // Link preview
		{GoHawaiiEdu, "http://go.hawaii.edu/j7L;", "j7L"},
		{GoHawaiiEdu, "http://go.hawaii.edu/admin", ""},
		{GoHawaiiEdu, "http://go.hawaii.edu/admin/", ""},
		{GoHawaiiEdu, "http://go.hawaii.edu/admin/index.php?", ""},
		{GoHawaiiEdu, "http://go.hawaii.edu/submit?", ""},
		{GoHawaiiEdu, "http://go.hawaii.edu/robert-j-elisberg/live-from-ces-day-two-the_b_416265.html", ""}, // ???
		{GoHawaiiEdu, "http://go.hawaii.edu/%E2%80%8Bhttps://www.star.hawaii.edu/studentinterface", ""},     // zero-width space
	}
	for i, tt := range tests {
		shortcode, err := tt.s.Clean(tt.url)
		if err != nil {
			t.Errorf("#%d: %v", i, err)
		} else if shortcode != tt.shortcode {
			t.Errorf("#%d: (%s).Clean(%q) = %q, want %q", i, tt.s.Name, tt.url, shortcode, tt.shortcode)
		}
	}
}
