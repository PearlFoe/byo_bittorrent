package p2p

import (
	"encoding/binary"
	"net"
	"net/http"
	"time"

	bencode "github.com/jackpal/bencode-go"
)

type Tracker struct {
	Url string
}

func (t *Tracker) unmarshalPeers(peerBencode string) ([]Peer, error) {
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

func (c *Tracker) RequestPeers() ([]Peer, error) {
	client := &http.Client{Timeout: time.Second * 60}

	response, err := client.Get(c.Url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	trackerResponse := &PeersBencode{}
	if err := bencode.Unmarshal(response.Body, trackerResponse); err != nil {
		return nil, err
	}

	peers, err := c.unmarshalPeers(trackerResponse.Peers)
	if err != nil {
		return nil, err
	}

	return peers, nil
}
