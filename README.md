# Build your own BitTorrent client
> Inspired by https://github.com/codecrafters-io/build-your-own-x

Simple [Bittorrent](https://www.bittorrent.org/beps/bep_0003.html) client written for learning perposes.

It does:
- allow to choose torrent file and path to place where file would be saved
- show cool progress bar

It does not (may be in some time):
- support downloading multiple files
- seed downloaded files
- daemon process to seed filed 24/7
- support µTP

## Usage
List available cli args
```shell
go run main.go -h
```

Run download
```shell
go run main.go -s . -t some_torrent_file.torrent
```