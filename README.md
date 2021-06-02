# URLHero

URLHero is a link resolver for current and defunct URL shorteners. It
uses link mappings from [URLTeam](https://wiki.archiveteam.org/index.php/URLTeam)
archives, dumps provided by shortener operators, and links captured by
the Internet Archive.

## Planned features

### Downloader

- Automatically download and process daily URLTeam releases.
- Hopefully gain access to [301Works dumps](301works.md).
- Switch to a torrent client that can scale to handle 1500 webseed
  items. [anacrolix/torrent](https://github.com/anacrolix/torrent) has
  [less mature webseed support](https://github.com/anacrolix/torrent/issues/465)
  and is relatively slow. [Transmission](https://transmissionbt.com/)
  was unable to handle all torrents, in simple tests.
- Support Internet Archive API authentication. For example,
  [URLTeamTorrentRelease2013July](https://archive.org/download/URLTeamTorrentRelease2013July)
  can only be downloaded when signed in.

### Link resolver

- Create link resolving website and API.
- Create Web Extension that redirects dead short links using URLHero.
- Proxy unknown shortener requests and contribute back to URLTeam
  dataset.
- Possibly fork [unshort.link](https://github.com/simonfrey/unshort.link).

### Parsing

- Process URLTeam first-generation TinyBack releases.
- Write custom CSV parser for qr-cx datasets to handle unescaped quotes.
- Full BEACON format spec compliance.

### Database

- Find a relational or key-value database with efficient compression.

## Contributing

There are many ways to contribute:

- File an issue or PR to submit a feature or bug report.
- Send link mappings for a URL shortener that you operate or have
  archived.
- Join URLTeam and help us archive at-risk shorteners by running the
  terroroftinytown project [in Docker](https://wiki.archiveteam.org/index.php/Running_Archive_Team_Projects_with_Docker#Basic_usage)
  or via the [Archive Team Warrior](https://wiki.archiveteam.org/index.php/ArchiveTeam_Warrior#Installing_and_running_with_Docker).

If you want to get in touch, join the
[#urlteam](https://webirc.hackint.org/#irc://irc.hackint.org/#urlteam)
channel on hackint or email me.

## License

This project is made available under the
[Mozilla Public License, v. 2.0](https://www.mozilla.org/en-US/MPL/2.0/).
