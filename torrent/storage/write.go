package storage

import (
	"os"
	"fmt"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"byo_bittorrent/torrent/metadata/file"
	"byo_bittorrent/torrent/p2p"
)

type Writer struct {
	Torrent *file.TorrentFile
	Bitfield p2p.Bitfield
	SaveDir string
}

func (w *Writer) fileName() string {
	return filepath.Join(w.SaveDir, w.Torrent.Name)
}

func (w *Writer) saveBlock(file *os.File, block *p2p.Block) error {
	offset := int64(block.Index * w.Torrent.PieceLength)

	if _, err := file.WriteAt(block.Buffer, offset); err != nil {
		return fmt.Errorf("failed to write block: %s", err)
	}
	log.Info("Wrote block", block.Index)
	return nil
}

func (w *Writer) CreateBitfield() {
	w.Bitfield = make([]byte, len(w.Torrent.PieceHashes))
}

func (w *Writer) Write(blocks chan p2p.Block) {
	log.Println("FILE PATH:", w.fileName())

	file, err := os.OpenFile(
		w.fileName(), 
		os.O_CREATE|os.O_WRONLY,
		0644,
	)
    if err != nil {
		log.Error(err)
        panic(err)
    }
	defer file.Close()

	wroteBlocks := 0
	for wroteBlocks < len(w.Torrent.PieceHashes) {
		block := <- blocks

		if err := w.saveBlock(file, &block); err != nil {
			// if failed to write block to file, return block to the channel
			blocks <- block
		} else {
			w.Bitfield.SetPiece(block.Index)
			wroteBlocks += 1
		}
	}
}
