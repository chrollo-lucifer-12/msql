package utils

import "encoding/binary"

func DecodeVarint(buffer []byte, offset int) (uint8, int64) {
	var size uint8
	var result int64

	for size < 9 {
		currentByte := buffer[offset]
		continuationBit := (currentByte >> 7) & 1
		dataBits := currentByte & 0x7F

		if size == 8 {
			result = (result << 8) | int64(currentByte)
		} else {
			result = (result << 7) | int64(dataBits)
		}

		size++
		offset++

		if continuationBit == 0 {
			break
		}
	}

	return size, result
}

func ReadBEWordAt(buffer []byte, offset int) uint16 {
	return binary.BigEndian.Uint16(buffer[offset:])
}
