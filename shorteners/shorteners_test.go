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
		{Allst, "http://a.ll.st/1NRwM3", "1NRwM3"},
		{Allst, "http://a.ll.st/Facebook", "Facebook"},
		{Allst, "http://a.ll.st/agentlocatorFB?linkId=104180290", "agentlocatorFB"},
		{Allst, "http://a.ll.st/Instagram%22,%22isCrawlable%22:true,%22thumbnail", "Instagram"}, // "
		{Allst, "http://a.ll.st:80/scmf/OrMCe04Lcp0lODk0BD1FrBcO2E4FP0NMEHFGSZ--Pq5q7EdIBj5D0RZwQ0r5O5LJxfQiUmcjxE_yFyVUmcC7Ue52R7KC2DlT6j1Anuut1CVBLh2fal1IZic40eX4xD2dJTg/PrJJpv", "PrJJpv"},
		{Allst, "http://a.ll.st:80/scmf/OrMCe04Lcp0lODk2Bzg71hcM2079O8ZJEHE_NJu-wtVr7D9JB0U8qWl1RzYCRZPJxfQiUmcjxE_yF9swgNxdUAkTP4vGed-VJvLu3uityvkzL-5fGDGJnyV0iKf6RXKdJQ/hiddenworldofdata", "hiddenworldofdata"},

		{Bfytw, "https://bfy.tw/PanS", "PanS"},
		{Bfytw, "http://bfy.tw/80xn=", "80xn"},
		{Bfytw, "http://bfy.tw:80/7JAH.", "7JAH"},
		{Bfytw, "http://bfy.tw:80/fb/7rt7", "fb"},
		{Bfytw, "http://bfy.tw:80/LOr7...", "LOr7"},
		{Bfytw, "http://bfy.tw:80/3hQy...You", "3hQy"},
		{Bfytw, "http://bfy.tw/BFsxrobots.txt", "BFsx"},
		{Bfytw, "http://bfy.tw/D9lj/robots.txt", "D9lj"},
		{Bfytw, "https://bfy.tw/4jz9ip124.41.235.255", "4jz9"},
		{Bfytw, "http://bfy.tw:80/5PrLhttp://bfy.tw/5PrL", "5PrL"},
		{Bfytw, "http://bfy.tw/Ej4D/wordpress/wp-content/uploads/kisaflo-top-loog.png", "Ej4D"},
		{Bfytw, "https://bfy.tw/Okad%22,%22e%22:%22link%22,%22t%22:%22https://bfy.tw/Okad", "Okad"}, // ""

		{Debli, "https://deb.li/hvPc", "hvPc"},
		{Debli, "http://deb.li:80/p/debian", "debian"},           // redirect preview
		{Debli, "http://deb.li:80/4BE7F84D.5040104@bzed.de", ""}, // mailing list redirect
		{Debli, "http://deb.li:80/imprint.html", ""},
		{Debli, "https://deb.li/static/pics/openlogo-50.png", ""},
		{Debli, "http://deb.li:80/log%20dari%20training%20Debian%20Women%20dengan%20tema%20%22Debian%20package%20informations%22%20dini%20hari%20tadi%20dapat%20dilihat%20di%20http://meetbot.debian.net/debian-women/2010/debian-women.2010-12-16-20.09.log.html", "log"}, // space
		{Debli, "http://deb.li:80/%3Ckey%3E", ""},  // <key> placeholder
		{Debli, "http://deb.li:80/%3Cname%3E", ""}, // <name> placeholder

		{GoHawaiiEdu, "https://go.hawaii.edu/34A", "34A"},
		{GoHawaiiEdu, "http://go.hawaii.edu:80/Vf+", "Vf"}, // redirect preview
		{GoHawaiiEdu, "http://go.hawaii.edu/3P6.", "3P6"},
		{GoHawaiiEdu, "http://go.hawaii.edu/j7L;", "j7L"},
		{GoHawaiiEdu, "http://go.hawaii.edu/fP7)", "fP7"},
		{GoHawaiiEdu, "http://go.hawaii.edu/admin", ""},
		{GoHawaiiEdu, "http://go.hawaii.edu:80/admin/", ""},
		{GoHawaiiEdu, "http://go.hawaii.edu:80/admin/index.php?", ""},
		{GoHawaiiEdu, "http://go.hawaii.edu:80/submit?", ""},
		{GoHawaiiEdu, "http://go.hawaii.edu:80/%E2%80%8Bhttps://www.star.hawaii.edu/studentinterface", ""}, // ZWSP
		{GoHawaiiEdu, "http://go.hawaii.edu:80/robert-j-elisberg/live-from-ces-day-two-the_b_416265.html", ""},

		{MobyTo, "http://moby.to//8dfstt", "8dfstt"},
		{MobyTo, "http://moby.to:80/368eck-", "368eck"},
		{MobyTo, "http://moby.to:80/8f9n7k--", "8f9n7k"},
		{MobyTo, "http://moby.to:80/4rcbqg%E2%80%9D", "4rcbqg"},   // ”
		{MobyTo, "http://moby.to:80/ac35nh%3C%3C", "ac35nh"},      // <<
		{MobyTo, "http://moby.to:80/ac35nh%C2%ABWoW..", "ac35nh"}, // «
		{MobyTo, "http://moby.to/1rrlao:view", "1rrlao"},
		{MobyTo, "https://moby.to/atmkt0:full", "atmkt0"},
		{MobyTo, "http://moby.to/8hmrkj:square", "8hmrkj"},
		{MobyTo, "http://moby.to:80/22ibg5:small", "22ibg5"},
		{MobyTo, "http://moby.to:80/91ttyo:large", "91ttyo"},
		{MobyTo, "http://moby.to:80/1b55uh:thumb", "1b55uh"},
		{MobyTo, "http://moby.to/08dlmz:thumbnail", "08dlmz"},
		{MobyTo, "http://moby.to:80/author/hermioneway/item/3417018", ""},
		{MobyTo, "http://moby.to/*", ""},
		{MobyTo, "http://moby.to:80/***", ""},
		{MobyTo, "http://moby.to:80/******", ""},
		{MobyTo, "http://moby.to/.*", ""},
		{MobyTo, "http://moby.to/.+", ""},

		{Qrcx, "http://qr.cx:80/)", ""},
		{Qrcx, "http://www.qr.cx/mQBM", "mQBM"},
		{Qrcx, "http://qr.cx/tEv/get", "tEv"}, // redirect preview
		{Qrcx, "http://qr.cx/sQ2U+", "sQ2U"},  // redirect preview
		{Qrcx, "http://qr.cx/plvd%5Dclick", "plvd"},
		{Qrcx, "http://qr.cx/plvd%5Dhttp:/qr.cx/plvd%5B/link%5D", "plvd"},
		{Qrcx, "http://qr.cx/yzj/img/301works.png", "yzj"},
		{Qrcx, "http://qr.cx:80/itZ/api.php", "itZ"},
		{Qrcx, "http://qr.cx:80/uqn/piwik.php", "uqn"},
		{Qrcx, "http://qr.cx/img/twitter_icon.png", ""},
		{Qrcx, "http://qr.cx/api.php", ""},
		{Qrcx, "http://qr.cx:80/deleted.php", ""},
		{Qrcx, "http://qr.cx:80/api/?bookmarklet=1&longurl=", ""},
		{Qrcx, "http://qr.cx:80/admin/latest.php?", ""},
		{Qrcx, "http://qr.cx:80/dataset/?flocxshorty=dataset", ""},
		{Qrcx, "http://qr.cx:80/qr/php/qr_img.php?", ""},
		{Qrcx, "http://qr.cx/qr/php/qr_img.php?e=M&s=9&d=http://qr.cx/1oz", "1oz"},
		{Qrcx, "http://qr.cx/qr/php/qr_img.php?e=M&s=9&d=http%3A%2F%2Fqr.cx%2Fyzj", "yzj"},
		{Qrcx, "http://qr.cx:80/http://qr.cx/about:blank", ""},
		{Qrcx, "http://qr.cx:80/http://maps.google.at/maps?", ""},

		{RedHt, "https://red.ht/3tg9nOW", "3tg9nOW"},
		{RedHt, "https://red.ht/3olOq1B@OpenRoboticsOrg", "3olOq1B"},
		{RedHt, "http://red.ht/1H7Wyt1@sklososky@FuturePOV", "1H7Wyt1"},
		{RedHt, "http://www.red.ht/forumswitzerland2017", "forumswitzerland2017"},
		{RedHt, "https://red.ht/SAPvirtualevent?sc_cid=701f2000000tzLzAAI", "SAPvirtualevent"},
		{RedHt, "http://red.ht/sitemap.xml", ""},
		{RedHt, "http://red.ht/static/graphics/fish-404.png", ""},

		{Rbgy, "https://rb.gy/bdb02v+", "bdb02v"}, // redirect preview
		{Rbgy, "https://rb.gy/auvwlc-", "auvwlc"},
		{Rbgy, "https://rb.gy/fpkgmy!", "fpkgmy"},
		{Rbgy, "http://rb.gy/txzznf_", "txzznf"},
		{Rbgy, "https://rb.gy/5xw62x%00", "5xw62x"},
		{Rbgy, "https://rb.gy/qntquc.Questions", "qntquc"},
		// {Rbgy, "http://rb.gy/uku.jog", "ukujog"},
		{Rbgy, "https://rb.gy/ouvdl3It's", "ouvdl3"},
		{Rbgy, "https://rb.gy/ddq3vo/UCaujr", "ddq3vo"},
		{Rbgy, "https://rb.gy/5wsqyxal-text&sr=1-3", "5wsqyx"},
		{Rbgy, "http://rb.gy/bj4..%3C/PAGE_TITLE%3E", "bj4"},
		{Rbgy, "https://rb.gy/gz35r7@YouTubeCreators", "gz35r7"},
		{Rbgy, "https://rb.gy/hjgyijrb.gy/hjgyij", "hjgyij"},
		{Rbgy, "http://rb.gy/ff7gyg/exercise-with-aerobic-video/", "ff7gyg"},
		{Rbgy, "https://rb.gy/wlshcv@drninaansary@A_Tabatabai@EllieGeranmayeh@ebtekarm@araghchi@milanimohsen@JafariPeyman@SadeghKharrazi@ahandjani", "wlshcv"},
		{Rbgy, "https://rb.gy/1zidswv%C3%ADa", "1zidsw"},                                                                 // í
		{Rbgy, "https://rb.gy/mi6dex%E5%B0%8F%E5%BF%83", "mi6dex"},                                                       // 小心
		{Rbgy, "https://rb.gy/vnaknf%E0%B8%AA%E0%B8%A1%E0%B8%B1%E0%B8%84%E0%B8%A3%E0%B8%87%E0%B8%B2%E0%B8%99", "vnaknf"}, // สมัครงาน
		{Rbgy, "https://rb.gy/sef%D0%B4%D0%BE%D0%BBa1x", "sefa1x"},                                                       // дол
		{Rbgy, "https://rb.gy/etqdt2%E2%97%84", "etqdt2"},                                                                // ◄
		{Rbgy, "https://rb.gy/etqdt2%e2%97%84%e2%80%8b%e2%80%8b%e2%80%8b", "etqdt2"},                                     // ◄ ZWSP ZWSP ZWSP

		{RedHt, "https://red.ht/sig%3E", "sig"}, // >
		{RedHt, "https://red.ht/dev-sandbox", "dev-sandbox"},
		{RedHt, "http://red.ht/1zzgkXp&esheet=51687448&newsitemid=20170921005271&lan=en-US&anchor=Red+Hat+blog&index=5&md5=7ea962d15a0e5bf8e35f385550f4decb", "1zzgkXp"},
		{RedHt, "http://red.ht/13LslKt&quot", "13LslKt"},
		{RedHt, "http://red.ht/2k3DNz3%E2%80%99", "2k3DNz3"}, // ’
		{RedHt, "http://red.ht/21Krw4z%C2%A0", "21Krw4z"},    // NBSP

		{ShortIm, "http://www.short.im:80/09u", "09u"},
		{ShortIm, "http://short.im:80/nova", "nova"},
		{ShortIm, "http://short.im:80/Christmas-Corner", "Christmas-Corner"},
		{ShortIm, "http://short.im/api", ""},
		{ShortIm, "http://short.im/api.php?short=http://short.im/1", "1"},
		{ShortIm, "http://short.im/api.php?url=http://example.com/very/long/?url", ""},
		{ShortIm, "http://short.im/donate", ""},
		{ShortIm, "http://short.im/tos", ""},
		{ShortIm, "http://short.im/warn", ""},
		{ShortIm, "http://short.im/feed.rss", ""},
		{ShortIm, "http://short.im:80/index.php", ""},
		{ShortIm, "http://short.im/stats.html", ""},
		{ShortIm, "http://short.im/developer.html", ""},
		{ShortIm, "http://short.im/multishrink.html", ""},
		{ShortIm, "https://short.im/modern_theme/build/img/bg.jpg", ""},
		{ShortIm, "http://short.im/stylesheets/Colaborate-fontfacekit/ColabLig-webfont.eot", ""},
		{ShortIm, "http://short.im/js/standard.js?rte=1&tm=2&dn=short.im&tid=1020", ""},
		{ShortIm, "http://short.im:80/caf/earch/tsc.php?&ses=14159950212c36f8a357f0b866615fe9dab1d7e009&200=MjA0MDg5ODA3&21=MTc0LjEyOS4yMzcuMTU3&681=MTQxNTk5NTAyMTJjMzZmOGEzNTdmMGI4NjY2MTVmZTlkYWIxZDdlMDA5&682=&616=MA==&crc=caada4a2a66dc82a6bcaf8e25c06fff5a7ccc2ec&cv=1", ""},
		// {ShortIm, "http://short.im:80/info/%3C%=urlKeyword%%3E.html?ses=Y3JlPTE0MTU5OTUwMjEmdGNpZD1zaG9ydC5pbTU0NjY1ZThjZTA5ZGEzLjc3OTA3MTEzJmZraT01NDY0JnRhc2s9c2VhcmNoJmRvbWFpbj1zaG9ydC5pbSZzPTZlM2U2YzA2YzdhMWRjN2MxYmRlJmxhbmd1YWdlPWVuJmFfaWQ9Mg==&keyword=%3C%=urlKeyword%%3E&token=%3C%=token%%3E", ""},

		{SUconnEdu, "http://s.uconn.edu/2by", "2by"},
		{SUconnEdu, "http://s.uconn.edu/ctsrc.", "ctsrc"},
		{SUconnEdu, "http://s.uconn.edu/fall-21-letter", "fall-21-letter"},
		{SUconnEdu, "http://s.uconn.edu/PreservingHistoricalResources", "preservinghistoricalresources"},
		{SUconnEdu, "https://s.uconn.edu/css/custom.css", ""},
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
