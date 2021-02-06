#!/bin/bash

# fetch_torrents.bash saves the torrent files for all terroroftinytown
# releases.

dir=${1?"Usage: $0 DIR"}

torrents=$(curl 'https://archive.org/services/search/v1/scrape?q=subject:terroroftinytown&count=10000' |
  jq -r '.items[].identifier | ("https://archive.org/download/" + . + "/" + . + "_archive.torrent")')

for url in $torrents; do
  out="$dir/$(basename "$url")"
  test -s "$out" || wget "$url" -O "$out"
done
