package p2p

import (
	"crypto/sha1"
	"bytes"
)

const MaxBlockSize = 16384

type Block struct {
	Index  int
	Length int
	Hash   [20]byte
	Buffer []byte
}

func (b *Block) CheckHash(buffer []byte) bool {
	hash := sha1.Sum(buffer)
	return bytes.Equal(hash[:], b.Hash[:])
}