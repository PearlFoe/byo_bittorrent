package cli

import (
	"time"

	pb "github.com/schollz/progressbar/v3"

	"byo_bittorrent/torrent/p2p"
)

type ProgressBar struct {
	Bitfield *p2p.Bitfield
}

func (b *ProgressBar) Start() {
	bfLength := b.Bitfield.Lenght()
	bar := pb.Default(int64(bfLength))

	downloaded := 0
	for downloaded < bfLength {
		newDownloaded := b.Bitfield.CountDownloaded()

		if newDownloaded > downloaded {
			bar.Add(newDownloaded - downloaded)
			downloaded = newDownloaded
		}

		time.Sleep(40 * time.Millisecond)
	}
}