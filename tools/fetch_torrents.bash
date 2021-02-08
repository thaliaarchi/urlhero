#!/bin/bash
# Copyright (c) 2021 Andrew Archibald
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

# fetch_torrents.bash saves the torrent files for all terroroftinytown
# releases.

dir=${1?"Usage: $0 DIR"}

torrents=$(curl 'https://archive.org/services/search/v1/scrape?q=subject:terroroftinytown&count=10000' |
  jq -r '.items[].identifier | ("https://archive.org/download/" + . + "/" + . + "_archive.torrent")')

for url in $torrents; do
  out="$dir/$(basename "$url")"
  test -s "$out" || wget "$url" -O "$out"
done
