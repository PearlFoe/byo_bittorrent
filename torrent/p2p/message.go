package p2p

import (
	"encoding/binary"
	"fmt"
	"io"
)

type messageID uint8

const (
    MsgChoke         messageID = 0
    MsgUnchoke       messageID = 1
    MsgInterested    messageID = 2
    MsgNotInterested messageID = 3
    MsgHave          messageID = 4
    MsgBitfield      messageID = 5
    MsgRequest       messageID = 6
    MsgPiece         messageID = 7
    MsgCancel        messageID = 8
)

type Message struct {
	ID      messageID
	Payload []byte
}


func (m *Message) Serialize() []byte {
	length := uint32(len(m.Payload) + 1) // +1 for id
	buf := make([]byte, length + 4)  // +4 for message length value

	binary.BigEndian.PutUint32(buf[:4], length)
	buf[4] = byte(m.ID)
	copy(buf[5:], m.Payload)

	return buf
}

func (m *Message) ParsePiece(pieceIndex int) ([]byte, error) {
	if m.ID != MsgPiece {
		return nil, fmt.Errorf("recieved invalid message code %d, expected %d", m.ID, MsgPiece)
	}

	index := int(binary.BigEndian.Uint32(m.Payload[:4]))
	if index != pieceIndex {
		return nil,  fmt.Errorf("recieved invalid piece index %d, expected %d", index, pieceIndex)
	}
	data := m.Payload[8:]

	return data, nil
}


func ReadMessage(r io.Reader) (*Message, error) {
	lengthBuffer := make([]byte, 4)
	if _, err := io.ReadFull(r, lengthBuffer); err != nil {
		return nil, err
	}

	messageLength := binary.BigEndian.Uint32(lengthBuffer)

	// keep-alive message
	if messageLength == 0 {
		return nil, nil
	}

	messageBuffer := make([]byte, messageLength)
	if _, err := io.ReadFull(r, messageBuffer); err != nil {
		return nil, err
	}

	message := Message{
		ID: messageID(messageBuffer[0]),
		Payload: messageBuffer[1:],
	}

	return &message, nil
}


func SendMessage(w io.Writer, m *Message) error {
	if _, err := w.Write(m.Serialize()); err != nil {
		return err
	}
	return nil
}
