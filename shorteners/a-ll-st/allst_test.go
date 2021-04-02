// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package allst

import (
	"fmt"
	"testing"
)

func TestGetIAShortcodes(t *testing.T) {
	t.Skip()
	shortcodes, err := GetIAShortcodes()
	if err != nil {
		t.Fatal(err)
	}
	if len(shortcodes) == 0 {
		t.Fatal("no shortcodes")
	}
	for _, s := range shortcodes {
		fmt.Println(s)
	}
}
