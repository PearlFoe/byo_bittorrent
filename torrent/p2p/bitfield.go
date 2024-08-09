package p2p

type Bitfield []byte


func (bf Bitfield) HasPiece(index int) bool {
	byteIndex := index / 8
	offset := index % 8
	if byteIndex < 0 || byteIndex >= len(bf) {
		return false
	}
	return bf[byteIndex]>>uint(7-offset)&1 != 0
}

func (bf Bitfield) SetPiece(index int) {
	byteIndex := index / 8
	offset := index % 8
	if byteIndex < 0 || byteIndex >= len(bf) {
		return
	}
	bf[byteIndex] |= 1 << uint(7 - offset)
}

func (bf Bitfield) CountDownloaded() int {
	sum := 0
	for i := 0; i < len(bf) * 8; i++ {
		if bf.HasPiece(i) {
			sum += 1
		}
	}
	return sum
}

func (bf Bitfield) Lenght() int {
	return len(bf)
}