package main

import (
	"byo_bittorrent/torrent/file"
	"byo_bittorrent/torrent/tracker"
	"fmt"
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

	tracker := &tracker.Tracker{Url: url}

	response, err := tracker.RequestPeers()
	if err != nil {
		fmt.Println("Ошибка получения пиров:", err)
	}

	fmt.Println(response)
}
