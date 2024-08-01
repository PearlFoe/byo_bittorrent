package file

import (
	"os"
	"io"
	"fmt"
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
	offset := int64(block.Index * w.Torrent.PieceLength)

	if _, err := file.WriteAt(block.Buffer, offset); err != nil {
		return fmt.Errorf("Failed to write block: ", err)
	}

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
	defer file.Close()

	for len(blocks) > 0 {
		block := <- blocks

		if err := w.saveBlock(file, &block); err != nil {
			// if failed to write block to file, return block to the channel
			blocks <- block
		}
	}

	return nil
}
