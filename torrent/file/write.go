package file

import (
	"byo_bittorrent/torrent/metadata/file"
	"byo_bittorrent/torrent/p2p"
)

type Writer struct {
	Torrent *file.TorrentFile
}

func (w *Writer) fileName() string {
	return w.Torrent.Name
}

func (w *Writer) Write(blocks chan p2p.Block) error {
	return nil
}
