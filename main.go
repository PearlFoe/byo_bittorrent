package main

import (
	"fmt"
	"math/rand"
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
	
	client := &p2p.Client{Torrent: content}
	if err := client.Start(&peers[rand.Intn(len(peers))]); err != nil {
		fmt.Println(err)
	}
	
}
