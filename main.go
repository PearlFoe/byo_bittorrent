package main

import (
	"fmt"
	//"io"
	"os"
	bencode "github.com/jackpal/bencode-go"
)

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}


func UnmarshalTorrent(path string) (*TorrentFile, error) {
	content := TorrentFile{}

	file, err := os.Open(path)
	if err != nil {
		return &content, err
	}

	if err = bencode.Unmarshal(file, &content); err != nil {
		return &content, err
	}
	
	return &content, err
}


func main(){
	filePath := "data/debian.iso.torrent"

	torrent, err := UnmarshalTorrent(filePath)

	if err != nil {
		fmt.Println("Ошибка парсинга файла", err)
	}

	fmt.Println(*torrent)
	
}