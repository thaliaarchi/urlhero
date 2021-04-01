// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package wwiki

import "testing"

func TestGetDumps(t *testing.T) {
	dumps, err := GetDumps()
	if err != nil {
		t.Fatal(err)
	}
	if len(dumps) == 0 {
		t.Fatalf("no dumps")
	}
}
