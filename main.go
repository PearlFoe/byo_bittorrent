package main

import (
	"byo_bittorrent/torrent/file"
	"fmt"
)

func main() {
	filePath := "data/debian.iso.torrent"
	content := new(file.TorrentFile)

	err := content.ReadFile(filePath)

	if err != nil {
		fmt.Println("Ошибка парсинга файла", err)
	}

	url, err := content.BuildTrackerUrl()
	if err != nil {
		fmt.Println("Ошибка формирования ссылки", err)
	}

	fmt.Println(url)
}
