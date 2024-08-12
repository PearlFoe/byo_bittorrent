package cli

import (
	"flag"
	"fmt"
)


type Args struct {
	TorrentFile string
	SaveDir string
}


func HandleArgs() (*Args, error) {
	var torrentFile string
	flag.StringVar(&torrentFile, "torrent", "", "Path to torrent file")
	flag.StringVar(&torrentFile, "t", "", "Path to torrent file")

	var saveDir string
	flag.StringVar(&saveDir, "save", "", "Path to dir where result should be saved")
	flag.StringVar(&saveDir, "s", "", "Path to dir where result should be saved")

	flag.Parse()

	if torrentFile == "" {
		fmt.Println("Torrent file path cant be empty")
		return nil, fmt.Errorf("invalid torrent file path")
	}

	if saveDir == "" {
		fmt.Println("Save dir path cant be empty")
		return nil, fmt.Errorf("invalid save dir path")
	}

	return &Args{TorrentFile: torrentFile, SaveDir: saveDir}, nil
}