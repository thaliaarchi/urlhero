# URLHero

URLHero downloads and processes URLTeam second-generation Terror of Tiny
Town link dumps.

## TODO

### Parsing

- Full BEACON format spec compliance
- Process URLTeam first-generation TinyBack releases

### Dependency issues

- anacrolix/torrent panic ([#465](https://github.com/anacrolix/torrent/issues/465))
- anacrolix/torrent webseed peer issues

## Planned Features

### Database

- Space-efficient storage
- Indexed URLs for regexp query API
- Automatic daily download of latest URLTeam release

### Web Extension

- Automatic shortener redirection to bypass tracking and archive dead
  shorteners
- Tracking parameter trimming
- Proxy unknown shortener requests and contribute back to URLTeam
  dataset
- Fork of [unshort.link](https://github.com/simonfrey/unshort.link)

## License

This project is made available under the
[Mozilla Public License](https://www.mozilla.org/en-US/MPL/2.0/).
