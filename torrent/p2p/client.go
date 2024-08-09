package p2p

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"byo_bittorrent/torrent/metadata/file"
)


type Client struct {
	Torrent *file.TorrentFile
	Choked  bool
}

func (c *Client) sendHandshake(connection net.Conn) error {
	handshake := &Handshake{
		Pstr:     "BitTorrent protocol",
		InfoHash: c.Torrent.InfoHash,
		PeerID:   c.Torrent.PeerID,
	}

	if _, err := connection.Write(handshake.Serialize()); err != nil {
		return err
	}

	response, err := ReadHandshake(connection)
	if err != nil {
		return err
	}

	if !bytes.Equal(handshake.InfoHash[:], response.InfoHash[:]) {
		return fmt.Errorf("expected infohash %x but got %x", handshake.InfoHash, response.InfoHash)
	}

	return nil
}

func (c *Client) waitBitfield(connection net.Conn) (*Bitfield, error) {
	message, err := ReadMessage(connection)
	if err != nil {
		return nil, nil
	}

	if message.ID != MsgBitfield {
		return nil, fmt.Errorf("recieved invalid message code %d, expected %d", message.ID, MsgBitfield)
	}

	bf := Bitfield(message.Payload)
	return &bf, nil
}

func (c *Client) waitUnchoke(connection net.Conn) error {
	message, err := ReadMessage(connection)
	if err != nil {
		return err
	}

	if message.ID != MsgUnchoke {
		return fmt.Errorf("recieved invalid message code %d, expected %d", message.ID, MsgUnchoke)
	}

	return nil
}

func (c *Client) sendInterested(connection net.Conn) error {
	message := &Message{ID: MsgInterested}
	if err := SendMessage(connection, message); err != nil {
		return err
	}
	return nil
}

func (c *Client) requestBlock(connection net.Conn, index, begin, length int) error {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], uint32(index))
	binary.BigEndian.PutUint32(payload[4:8], uint32(begin))
	binary.BigEndian.PutUint32(payload[8:12], uint32(length))

	message := &Message{ID: MsgRequest, Payload: payload}
	if err := SendMessage(connection, message); err != nil {
		return err
	}

	return nil
}

func (c *Client) downloadBlock(connection net.Conn, block *Block) error {
	if err := c.requestBlock(connection, block.Index, 0, MaxBlockSize); err != nil {
		return err
	}

	begin := 0
	buff := make([]byte, block.Length)

	for begin < block.Length-1 {
		message, err := ReadMessage(connection)
		if err != nil {
			return err
		}

		switch message.ID {
		case MsgChoke:
			c.Choked = true
		case MsgUnchoke:
			c.Choked = false
		default:
		}

		if c.Choked {
			continue
		}

		if message.ID == MsgPiece {
			piece, err := message.ParsePiece(block.Index)
			if err != nil {
				return err
			}
			copy(buff[begin:], piece)
			begin += len(piece)

			if err := c.requestBlock(connection, block.Index, begin, MaxBlockSize); err != nil {
				return err
			}
		}
	}

	if block.CheckHash(buff) {
		if block.Buffer == nil {
			block.Buffer = make([]byte, block.Length)
		}
		copy(block.Buffer, buff)
	}

	return nil
}


func (c *Client) Start(peer *Peer, toDownload, toSave chan Block, wg *sync.WaitGroup) {
	defer wg.Done()
	
	log.Info("Connecting to peer", peer.String())

	connection, err := net.DialTimeout("tcp", peer.String(), 60*time.Second)
	if err != nil {
		log.Warn(err)
		return
	}
	defer connection.Close()

	log.Info("Connected to peer", peer.String())

	if err := c.sendHandshake(connection); err != nil {
		log.Warn(err)
		return
	}

	c.Choked = true

	log.Info("Handshaked peer", peer.String())

	bitfield, err := c.waitBitfield(connection)
	if err != nil {
		log.Warn(err)
		return
	}

	log.Info("Recieved bitfield")

	// TODO: Find out why i get EOF sometimes
	if err := c.waitUnchoke(connection); err != nil {
		log.Warn(err)
		return
	}

	c.Choked = false
	log.Info("Unchoked")

	if err := c.sendInterested(connection); err != nil {
		log.Warn(err)
		return
	}

	log.Info("Sent interested")

	for len(toDownload) > 0 {
		block := <- toDownload
		if !bitfield.HasPiece(block.Index) {
			toDownload <- block
			continue
		}

		if err := c.downloadBlock(connection, &block); err != nil {
			toDownload <- block
			log.Warn(err)
			return
		}

		if len(block.Buffer) > 0{
			toSave <- block
			log.Printf("Finished download: %d / %d \n", block.Index + 1, len(c.Torrent.PieceHashes))
		}
	}
}
