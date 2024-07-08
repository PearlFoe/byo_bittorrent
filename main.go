package main

import (
	"fmt"

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

	client := &p2p.Client{Url: url, Torrent: content}

	
	if err := client.RequestPeers(); err != nil {
		fmt.Println("Ошибка получения пиров:", err)
	}

	if err := client.Start(&client.Peers[0]); err != nil {
		fmt.Println("Ошибка подключения к пиру:", err)
	}
	
}
