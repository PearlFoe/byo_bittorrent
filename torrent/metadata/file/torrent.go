package file

import (
	"crypto/rand"
	bencode "github.com/jackpal/bencode-go"
	"net/url"
	"os"
	"strconv"
)

const Port uint16 = 6881

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

func (t *TorrentFile) ReadFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	var bencodeFile = BencodeFile{}

	if err = bencode.Unmarshal(file, &bencodeFile); err != nil {
		return err
	}

	infoHash, err := bencodeFile.Info.Hash()
	if err != nil {
		return err
	}

	pieceHashes, err := bencodeFile.Info.splitPieceHashes()
	if err != nil {
		return err
	}

	t.Announce = bencodeFile.Announce
	t.InfoHash = infoHash
	t.PieceHashes = pieceHashes
	t.PieceLength = bencodeFile.Info.PieceLength
	t.Length = bencodeFile.Info.Length
	t.Name = bencodeFile.Info.Name

	return nil
}

func (t *TorrentFile) BuildTrackerUrl() (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}

	var peerID [20]byte
	_, err = rand.Read(peerID[:])
	if err != nil {
		return "", err
	}

	params := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(Port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Length)},
	}
	base.RawQuery = params.Encode()
	return base.String(), nil
}
