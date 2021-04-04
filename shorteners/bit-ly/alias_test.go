// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package bitly

import "testing"

func TestIsAlias(t *testing.T) {
	nonAliases := []string{
		"urlte.am",
		"wiki.archiveteam.org",
	}
	for _, host := range Aliases {
		checkAlias(t, host, true)
	}
	for _, host := range nonAliases {
		checkAlias(t, host, false)
	}
}

func checkAlias(t *testing.T, host string, want bool) {
	t.Helper()
	isAlias, err := IsAlias(host)
	if err != nil {
		t.Errorf("IsAlias(%s): %v", host, err)
	} else if isAlias != want {
		t.Errorf("IsAlias(%s) = %t, want %t", host, isAlias, want)
	}
}
