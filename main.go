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
		fmt.Println("Ошибка парсинга файла:", err)
	}

	url, err := content.BuildTrackerUrl()
	if err != nil {
		fmt.Println("Ошибка формирования ссылки:", err)
	}

	tracker := &p2p.Tracker{Url: url}
	peers, err := tracker.RequestPeers()
	if err != nil {
		fmt.Println("Ошибка получения пиров:", err)
	}
	
	client := &p2p.Client{Torrent: content}
	if err := client.Start(&peers[rand.Intn(len(peers))]); err != nil {
		fmt.Println(err)
	}
	
}
