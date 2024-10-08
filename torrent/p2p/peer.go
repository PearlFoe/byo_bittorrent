package p2p

import (
	"net"
	"strconv"
)

type Peer struct {
	IP   net.IP
	Port uint16
}


func (p *Peer) String () string {
	return net.JoinHostPort(p.IP.String(), strconv.Itoa(int(p.Port)))
}

type PeersBencode struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}
