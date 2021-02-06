// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package beacon

import "testing"

func TestSplitMeta(t *testing.T) {
	tests := []struct {
		raw  string
		meta MetaField
		err  bool
	}{
		{"FORMAT: BEACON", MetaField{"FORMAT", "BEACON"}, false},
		{"PREFIX: http://example.org/id/", MetaField{"PREFIX", "http://example.org/id/"}, false},
		{"TARGET: http://example.com/about/", MetaField{"TARGET", "http://example.com/about/"}, false},
		{"FORMAT BEACON", MetaField{"FORMAT", "BEACON"}, false},
		{"FORMAT\tBEACON", MetaField{"FORMAT", "BEACON"}, false},
		{"FORMAT:\t \tBEACON", MetaField{"FORMAT", "BEACON"}, false},
		{"FORMAT : BEACON", MetaField{"FORMAT", ": BEACON"}, false},
		{"FORMAT", MetaField{}, true},
		{"http://example.org/id/", MetaField{}, true},
	}
	for i, test := range tests {
		meta, err := splitMeta(test.raw)
		if (err != nil) != test.err {
			t.Errorf("#%d: splitMeta(%q) got err %v, want %t", i, test.raw, err, test.err)
			continue
		}
		if meta != test.meta {
			t.Errorf("#%d: splitMeta(%q) got %v, want %v", i, test.raw, meta, test.meta)
		}
	}
}
