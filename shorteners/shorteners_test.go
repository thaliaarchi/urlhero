// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package shorteners

import (
	"testing"
)

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
