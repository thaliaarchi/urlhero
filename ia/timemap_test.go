// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ia

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestDecodeDigest(t *testing.T) {
	tests := []struct {
		digest, sha1 string
	}{
		{"TS3WOHL6SGIAF7FIMPABIV7CO27YXCM7", "9cb7671d7e919002fca863c01457e276bf8b899f"},
		{"7DNQJBSVVVSST6ZRKPCIEE6VNJWOP3UE", "f8db048655ad6529fb3153c48213d56a6ce7ee84"},
		{"NNZV4FHGZW2OUD5KFR6EN3P4EERARKXU", "6b735e14e6cdb4ea0faa2c7c46edfc212208aaf4"},
		{"KKWAWM37XU5PE3SZJY626H76D6NSMRLO", "52ac0b337fbd3af26e594e3daf1ffe1f9b26456e"},
		{"76POXGGRGPS6NJXWUIM4WHD5SNZ5CA6Q", "ff9eeb98d133e5e6a6f6a219cb1c7d9373d103d0"},
		{"4PVTM2ICN6HDOXXJ4YIX44DR66IA5RV4", "e3eb3669026f8e375ee9e6117e7071f7900ec6bc"},
	}
	for _, test := range tests {
		sha1, err := hex.DecodeString(test.sha1)
		if err != nil {
			t.Error(err)
			continue
		}
		b, err := DecodeDigest(test.digest)
		if err != nil {
			t.Errorf("DecodeDigest(%q) %v", test.digest, err)
			continue
		}
		if !bytes.Equal(b[:], sha1) {
			t.Errorf("got %x, want %x", b[:], sha1)
		}
	}
}
