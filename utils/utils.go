package utils

import (
	"encoding/binary"
	"math"
)

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

func ReadI8At(input []byte, offset int) int64 {
	return int64(int8(input[offset]))
}

func ReadI16At(input []byte, offset int) int64 {
	return int64(int16(binary.BigEndian.Uint16(input[offset : offset+2])))
}

func ReadI24At(input []byte, offset int) int64 {
	v := int32(input[offset])<<16 |
		int32(input[offset+1])<<8 |
		int32(input[offset+2])

	if v&0x800000 != 0 {
		v |= ^0xFFFFFF
	}

	return int64(v)
}

func ReadI32At(input []byte, offset int) int64 {
	return int64(int32(binary.BigEndian.Uint32(input[offset : offset+4])))
}

func ReadI48At(input []byte, offset int) int64 {
	v := int64(input[offset])<<40 |
		int64(input[offset+1])<<32 |
		int64(input[offset+2])<<24 |
		int64(input[offset+3])<<16 |
		int64(input[offset+4])<<8 |
		int64(input[offset+5])

	if v&(1<<47) != 0 {
		v |= ^int64(0xFFFFFFFFFFFF)
	}

	return v
}

func ReadI64At(input []byte, offset int) int64 {
	return int64(binary.BigEndian.Uint64(input[offset : offset+8]))
}

func ReadF64At(input []byte, offset int) float64 {
	bits := binary.BigEndian.Uint64(input[offset : offset+8])
	return math.Float64frombits(bits)
}
