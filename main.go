package main

import (
	"fmt"
	"math/rand"
	"byo_bittorrent/torrent/metadata/file"
	"byo_bittorrent/torrent/p2p"
)

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
