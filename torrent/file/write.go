package file

import (
	"os"
	"byo_bittorrent/torrent/metadata/file"
	"byo_bittorrent/torrent/p2p"
)

type Writer struct {
	Torrent *file.TorrentFile
}

func (w *Writer) fileName() string {
	return w.Torrent.Name
}

func (w *Writer) saveBlock(file *os.File, block *p2p.Block) error {
	return nil
}

func (w *Writer) Write(blocks chan p2p.Block) error {
	file, err := os.OpenFile(
		w.fileName(), 
		os.O_CREATE|os.O_WRONLY,
		0644,
	)
    if err != nil {
        panic(err)
    }

	for len(blocks) > 0 {
		block := <- blocks

		if err := w.saveBlock(file, &block); err != nil {
			blocks <- block
		}
	}

	return nil
}
