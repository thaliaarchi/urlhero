// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package bitly

import "testing"

func TestIsAlias(t *testing.T) {
	tests := []struct {
		domain  string
		isAlias bool
	}{
		// Known bit.ly aliases
		// https://wiki.archiveteam.org/index.php?title=URLTeam#bit.ly_aliases
		{"0-60.in", true},
		{"030mag.de", true},
		{"1.azdhs.gov", true},
		{"20ss.nyc", true},
		{"360myte.ch", true},
		{"4sq.com", true},
		{"511t.ac", true},
		{"6.wate.com", true},
		{"7gee.se", true},
		{"8legs.co", true},
		{"9l7gf7x5o9v.xyz", true},
		{"IMDB.to", true},
		{"abcn.ws", true},
		{"accntu.re", true},
		{"actf.tv", true},
		{"adobe.ly", true},
		{"adweek.it", true},
		{"amzn.to", true},
		{"apple.co", true},
		{"arfo.sk", true},
		{"audi.us", true},
		{"autism.link", true},
		{"azc.cc", true},
		{"bbc.in", true},
		{"beta.works", true},
		{"binged.it", true},
		{"bitly.is", true},
		{"blizz.ly", true},
		{"bloom.bg", true},
		{"bo.st", true},
		{"bravo.ly", true},
		{"bzfd.it", true},
		{"canva.link", true},
		{"cb.com", true},
		{"chfstps.co", true},
		{"chzb.gr", true},
		{"cnet.co", true},
		{"cnnmon.ie", true},
		{"comca.st", true},
		{"conta.cc", true},
		{"crks.me", true},
		{"d-mk.co", true},
		{"don8blood.com", true},
		{"econ.st", true},
		{"engri.sh", true},
		{"eonli.ne", true},
		{"es.pn", true},
		{"fltsim.me", true},
		{"fxn.ws", true},
		{"got.cr", true},
		{"hrkey.co", true},
		{"hub.am", true},
		{"huff.to", true},
		{"idle-empi.re", true},
		{"ift.tt", true},
		{"kck.st", true},
		{"kore.us", true},
		{"krg.bz", true},
		{"lat.ms", true},
		{"laurc.ro", true},
		{"lemde.fr", true},
		{"lft.to", true},
		{"m.ttmask.com", true},
		{"marsdd.it", true},
		{"mbist.ro", true},
		{"mojo.ly", true},
		{"mttr.io", true},
		{"nydn.us", true},
		{"nyti.ms", true},
		{"oculta.bit.ly", true},
		{"on.fb.me", true},
		{"on.natgeo.com", true},
		{"red.ht", true},
		{"reut.rs", true},
		{"sdut.us", true},
		{"snd.sc", true},
		{"spoti.fi", true},
		{"stanford.io", true},
		{"tcrn.ch", true},
		{"ti.me", true},
		{"usat.ly", true},
		{"wapo.st", true},
		{"wgrd.tech", true},
		{"wttw.me", true},
		{"xar.ph", true},
		{"yhoo.it", true},

		// Returns false
		//   1.usa.gov
		//   aje.me
		//   arcg.is
		//   atfp.co
		//   bbybgrl.com
		//   bitly.com
		//   bitlymail.com
		//   corb.is
		//   cot.ag
		//   curbed.cc
		//   gph.is
		//   j.mp
		//   lego.build
		//   on.rare.us
		//   on.si.com
		//   qr.cm
		//   theatln.tc
		//   vstphl.ly
		//   whrt.it
		//   www.bitly.com
		//   zmb.me

		// Lookup: no such host
		//   carrot.cr
		//   emarketee.rs
		//   grd.to
		//   jrnl.to
		//   stnfy.com

		// Not a host
		//   feedly.com/k/

		{"urlte.am", false},
		{"wiki.archiveteam.org", false},
	}
	for _, test := range tests {
		isAlias, err := IsAlias(test.domain)
		if err != nil {
			t.Errorf("IsAlias(%s): %v", test.domain, err)
			continue
		}
		if isAlias != test.isAlias {
			t.Errorf("IsAlias(%s) = %t, want %t", test.domain, isAlias, test.isAlias)
		}
	}
}
