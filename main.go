package main

import (
	"fmt"
	"sync"
	"byo_bittorrent/torrent/metadata/file"
	"byo_bittorrent/torrent/p2p"
)


func generateDownloadChannel(torrent *file.TorrentFile) chan p2p.Block{
	blocks := make(chan p2p.Block, len(torrent.PieceHashes))

	for hashIndex, hash := range torrent.PieceHashes {
		blocks <- p2p.Block{
			Index: hashIndex,
			Length: torrent.CalculatePieceSize(hashIndex),
			Hash: hash,
		}
	} 
	return blocks
}


func generateSaveChannel(torrent *file.TorrentFile) chan p2p.Block{
	return make(chan p2p.Block, len(torrent.PieceHashes))
}


func main() {
	filePath := "data/debian.iso.torrent"
	content := &file.TorrentFile{}

	err := content.ReadFile(filePath)

	if err != nil {
		fmt.Println("Failed to parse file:", err)
	}

	url, err := content.BuildTrackerUrl()
	if err != nil {
		fmt.Println("Failed to build tracker ulr:", err)
	}

	tracker := &p2p.Tracker{Url: url}
	peers, err := tracker.RequestPeers()
	if err != nil {
		fmt.Println("Failed to request peers list:", err)
	}
	
	toDownload := generateDownloadChannel(content)
	toSave := generateSaveChannel(content)

	var wg sync.WaitGroup
	for _, peer := range peers {
		wg.Add(1)

		client := &p2p.Client{Torrent: content}
		go client.Start(&peer, toDownload, toSave, &wg)
	}
	wg.Wait()
	fmt.Println("All workers done")
}
