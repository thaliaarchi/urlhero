// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package bitly handles the bit.ly link shortener and its aliases.
package bitly

import (
	"encoding/binary"
	"net"
)

// Aliases is a list of the hostnames of known bit.ly aliases.
// See https://wiki.archiveteam.org/index.php?title=URLTeam#bit.ly_aliases
var Aliases = []string{
	"0-60.in",
	"030mag.de",
	"1.azdhs.gov",
	"20ss.nyc",
	"360myte.ch",
	"4sq.com",
	"511t.ac",
	"6.wate.com",
	"7gee.se",
	"8legs.co",
	"9l7gf7x5o9v.xyz",
	"imdb.to",
	"abcn.ws",
	"accntu.re",
	"actf.tv",
	"adobe.ly",
	"adweek.it",
	"amzn.to",
	"apple.co",
	"arfo.sk",
	"audi.us",
	"autism.link",
	"azc.cc",
	"bbc.in",
	"beta.works",
	"binged.it",
	"bitly.is",
	"blizz.ly",
	"bloom.bg",
	"bo.st",
	"bravo.ly",
	"bzfd.it",
	"canva.link",
	"cb.com",
	"chfstps.co",
	"chzb.gr",
	"cnet.co",
	"cnnmon.ie",
	"comca.st",
	"conta.cc",
	"crks.me",
	"d-mk.co",
	"don8blood.com",
	"econ.st",
	"engri.sh",
	"eonli.ne",
	"es.pn",
	"fltsim.me",
	"fxn.ws",
	"got.cr",
	"hrkey.co",
	"hub.am",
	"huff.to",
	"idle-empi.re",
	"ift.tt",
	"kck.st",
	"kore.us",
	"krg.bz",
	"lat.ms",
	"laurc.ro",
	"lemde.fr",
	"lft.to",
	"m.ttmask.com",
	"marsdd.it",
	"mbist.ro",
	"mojo.ly",
	"mttr.io",
	"nydn.us",
	"nyti.ms",
	"oculta.bit.ly",
	"on.fb.me",
	"on.natgeo.com",
	"red.ht",
	"reut.rs",
	"sdut.us",
	"snd.sc",
	"spoti.fi",
	"stanford.io",
	"tcrn.ch",
	"ti.me",
	"usat.ly",
	"wapo.st",
	"wgrd.tech",
	"wttw.me",
	"xar.ph",
	"yhoo.it",
}

/*
Dead or problematic aliases from wiki list:

DNS not in the range 67.199.248.10 to 67.199.248.13
	1.usa.gov
	aje.me
	arcg.is
	atfp.co
	bbybgrl.com
	bitly.com
	bitlymail.com
	corb.is
	cot.ag
	curbed.cc
	gph.is
	j.mp
	lego.build
	on.rare.us
	on.si.com
	qr.cm
	theatln.tc
	vstphl.ly
	whrt.it
	www.bitly.com
	zmb.me

Lookup: no such host
	carrot.cr
	emarketee.rs
	grd.to
	jrnl.to
	stnfy.com

Not a host
	feedly.com/k/
*/

// IsAlias checks whether the domain is a bit.ly alias.
func IsAlias(domain string) (bool, error) {
	// Check if the IP is in the range 67.199.248.10 to 67.199.248.13.
	lo := binary.BigEndian.Uint32([]byte{67, 199, 248, 10})
	hi := binary.BigEndian.Uint32([]byte{67, 199, 248, 13})
	ips, err := net.LookupIP(domain)
	if err != nil {
		return false, err
	}
	for _, ip := range ips {
		ip = ip.To4()
		if ip == nil {
			continue
		}
		ip4 := binary.BigEndian.Uint32(ip)
		if lo <= ip4 && ip4 <= hi {
			return true, nil
		}
	}
	return false, nil
}
