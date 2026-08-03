package cursor

import (
	"fmt"

	"github.com/msql/utils"
)

type RecordFieldType int

const (
	NULL RecordFieldType = iota
	I8
	I16
	I24
	I32
	I48
	I64
	Float
	Zero
	One
	String
	Blob
)

type RecordField struct {
	Offset    uint
	FieldType RecordFieldType
}

type RecordHeader struct {
	Fields []RecordField
}

func parseRecordHeader(buffer []byte) (*RecordHeader, error) {
	size, headerLength := utils.DecodeVarint(buffer, 0)
	buffer = buffer[size:headerLength]

	fields := make([]RecordField, 0)
	currentOffset := headerLength

	for len(buffer) > 0 {
		discriminantSize, discriminant := utils.DecodeVarint(buffer, 0)
		buffer = buffer[discriminantSize:]

		var (
			fieldType RecordFieldType
			fieldSize int
		)

		switch {
		case discriminant == 0:
			fieldType = NULL
			fieldSize = 0

		case discriminant == 1:
			fieldType = I8
			fieldSize = 1

		case discriminant == 2:
			fieldType = I16
			fieldSize = 2

		case discriminant == 3:
			fieldType = I24
			fieldSize = 3

		case discriminant == 4:
			fieldType = I32
			fieldSize = 4

		case discriminant == 5:
			fieldType = I48
			fieldSize = 6

		case discriminant == 6:
			fieldType = I64
			fieldSize = 8

		case discriminant == 7:
			fieldType = Float
			fieldSize = 8

		case discriminant == 8:
			fieldType = Zero
			fieldSize = 0

		case discriminant == 9:
			fieldType = One
			fieldSize = 0

		case discriminant >= 12 && discriminant%2 == 0:
			size := int((discriminant - 12) / 2)
			fieldType = Blob
			fieldSize = size

		case discriminant >= 13 && discriminant%2 == 1:
			size := int((discriminant - 13) / 2)
			fieldType = String
			fieldSize = size

		default:
			return nil, fmt.Errorf("unsupported field type: %d", discriminant)
		}

		fields = append(fields,
			RecordField{
				Offset:    uint(currentOffset),
				FieldType: fieldType,
			},
		)

		currentOffset += int64(fieldSize)
	}

	return &RecordHeader{
		Fields: fields,
	}, nil
}
