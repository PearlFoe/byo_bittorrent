package main

import (
	"fmt"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"

	"byo_bittorrent/cli"
	"byo_bittorrent/torrent/metadata/file"
	"byo_bittorrent/torrent/p2p"
	"byo_bittorrent/torrent/storage"
)

func generateDownloadChannel(torrent *file.TorrentFile) chan p2p.Block {
	blocks := make(chan p2p.Block, len(torrent.PieceHashes))

	for hashIndex, hash := range torrent.PieceHashes {
		blocks <- p2p.Block{
			Index:  hashIndex,
			Length: torrent.CalculatePieceSize(hashIndex),
			Hash:   hash,
		}
	}
	return blocks
}

func generateSaveChannel(torrent *file.TorrentFile) chan p2p.Block {
	return make(chan p2p.Block, len(torrent.PieceHashes))
}

func main() {
	logFile, err := os.OpenFile("logs/main.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.InfoLevel)

	args, err := cli.ParseArgs()
	if err != nil {
		log.Error(err)
		return
	}

	content := &file.TorrentFile{}
	if err := content.ReadFile(args.TorrentFile); err != nil {
		log.Error("Failed to parse file:", err)
		return
	}

	url, err := content.BuildTrackerUrl()
	if err != nil {
		log.Error("Failed to build tracker ulr:", err)
		return
	}

	tracker := &p2p.Tracker{Url: url}
	peers, err := tracker.RequestPeers()
	if err != nil {
		log.Error("Failed to request peers list:", err)
		return
	}

	toDownload := generateDownloadChannel(content)
	toSave := generateSaveChannel(content)

	var wg sync.WaitGroup
	for _, peer := range peers {
		wg.Add(1)

		client := &p2p.Client{Torrent: content}
		go client.Start(&peer, toDownload, toSave, &wg)
	}

	writer := &storage.Writer{Torrent: content, SaveDir: args.SaveDir}
	writer.CreateBitfield()

	progressBar := &cli.ProgressBar{Bitfield: &writer.Bitfield}
	go progressBar.Start()

	writer.Write(toSave)

	wg.Wait()

	fmt.Println("All workers done")
}
