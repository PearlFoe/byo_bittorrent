package main

import (
	"fmt"
	"bytes"
	"strconv"
	"crypto/rand"
	"crypto/sha1"
	"net/url"
	"os"
	bencode "github.com/jackpal/bencode-go"
)


const Port uint16 = 6881 


type BencodeInfo struct {
    Pieces      string `bencode:"pieces"`
    PieceLength int    `bencode:"piece length"`
    Length      int    `bencode:"length"`
    Name        string `bencode:"name"`
}

func (i *BencodeInfo) Hash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}

type BencodeFile struct {
    Announce string      `bencode:"announce"`
    Info     BencodeInfo `bencode:"info"`
}


func (b * BencodeFile) UnmarshalTorrent(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	if err = bencode.Unmarshal(file, &b); err != nil {
		return err
	}
	
	return nil
}


func (b * BencodeFile) BuildTracerUrl() (string, error) {
	base, err := url.Parse(b.Announce)
	if err != nil {
		return "", err
	}

	var peerID [20]byte
	_, err = rand.Read(peerID[:])
	if err != nil {
		return "", err
	}

	infoHash, err := b.Info.Hash()
	if err != nil {
		return "", err
	}

    params := url.Values{
        "info_hash":  []string{string(infoHash[:])},
        "peer_id":    []string{string(peerID[:])},
        "port":       []string{strconv.Itoa(int(Port))},
        "uploaded":   []string{"0"},
        "downloaded": []string{"0"},
        "compact":    []string{"1"},
        "left":       []string{strconv.Itoa(b.Info.Length)},
    }
    base.RawQuery = params.Encode()
    return base.String(), nil
}

func main(){
	filePath := "data/debian.iso.torrent"
	content := new(BencodeFile)

	err := content.UnmarshalTorrent(filePath)

	if err != nil {
		fmt.Println("Ошибка парсинга файла", err)
	}

	url, err := content.BuildTracerUrl()
	if err != nil {
		fmt.Println("Ошибка формирования ссылки", err)
	}

	fmt.Println(url)
}