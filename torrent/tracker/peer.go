package tracker

import (
	"encoding/binary"
	"net"
	"net/http"
	"time"
	bencode "github.com/jackpal/bencode-go"
)


type TrackerResponse struct {
	Interval int 	`bencode:"interval"`
	Peers    string `bencode:"peers"`
}


type Tracker struct {
	Url string
}


type Peer struct {
	IP net.IP
	Port uint16
}

func (t *Tracker) UnmarshalPeers (peerBencode string) ([]Peer, error) {
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

func (t *Tracker) RequestPeers () ([]Peer, error) {
	client := &http.Client{Timeout: time.Second * 60}

	resp, err := client.Get(t.Url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	trackerResponse := &TrackerResponse{}

	if err := bencode.Unmarshal(resp.Body, trackerResponse); err != nil{
		return nil, err
	}

	peers, err := t.UnmarshalPeers(trackerResponse.Peers)
	if err != nil {
		return nil, err
	}

	return peers, nil
}
