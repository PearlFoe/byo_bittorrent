package p2p

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"time"

	"byo_bittorrent/torrent/metadata/file"

	bencode "github.com/jackpal/bencode-go"
)

type PeersResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

type Client struct {
	Url   string
	Peers []Peer
}

func (c *Client) unmarshalPeers(peerBencode string) ([]Peer, error) {
	const peerSize = 6 // 4 for IP, 2 for port

	peerByte := []byte(peerBencode)
	numPeers := len(peerByte) / peerSize

	peers := make([]Peer, numPeers)

	for i := 0; i < numPeers; i++ {
		offset := i * peerSize
		peers[i].IP = net.IP(peerByte[offset : offset+4])
		peers[i].Port = binary.BigEndian.Uint16(peerByte[offset+4 : offset+6])
	}
	return peers, nil
}

func (c *Client) RequestPeers() error {
	client := &http.Client{Timeout: time.Second * 60}

	resp, err := client.Get(c.Url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	trackerResponse := &PeersResponse{}

	if err := bencode.Unmarshal(resp.Body, trackerResponse); err != nil {
		return err
	}

	peers, err := c.unmarshalPeers(trackerResponse.Peers)
	if err != nil {
		return err
	}

	c.Peers = peers

	return nil
}

func (c *Client) SendHandshake(peer *Peer, torrent *file.TorrentFile) error {
	connection, err := net.DialTimeout("tcp", peer.String(), 30*time.Second)
	if err != nil {
		return err
	}

	handshake := &Handshake{
		Pstr:     "BitTorrent protocol",
		InfoHash: torrent.InfoHash,
		PeerID:   torrent.PeerID,
	}

	_, err = connection.Write(handshake.Serialize())
	if err != nil {
		return err
	}

	response, err := ReadHandshake(connection)
	if err != nil {
		return err
	}

	if !bytes.Equal(handshake.InfoHash[:], response.InfoHash[:]) {
		return fmt.Errorf("expected infohash %x but got %x", handshake.InfoHash, response.InfoHash)
	}

	// Логика получения данных и сохранения на диск

	return nil
}
