package cursor

import (
	"fmt"

	"github.com/msql/utils"
)

type RecordFieldType struct {
	Kind int
	Size int
}

const (
	NULL = iota
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

	buffer = buffer[size:int(headerLength)]

	fields := make([]RecordField, 0)
	currentOffset := int(headerLength)

	for len(buffer) > 0 {
		discriminantSize, discriminant := utils.DecodeVarint(buffer, 0)
		buffer = buffer[int(discriminantSize):]

		var (
			fieldType RecordFieldType
			fieldSize int
		)

		switch {
		case discriminant == 0:
			fieldType = RecordFieldType{Kind: NULL}
			fieldSize = 0

		case discriminant == 1:
			fieldType = RecordFieldType{Kind: I8}
			fieldSize = 1

		case discriminant == 2:
			fieldType = RecordFieldType{Kind: I16}
			fieldSize = 2

		case discriminant == 3:
			fieldType = RecordFieldType{Kind: I24}
			fieldSize = 3

		case discriminant == 4:
			fieldType = RecordFieldType{Kind: I32}
			fieldSize = 4

		case discriminant == 5:
			fieldType = RecordFieldType{Kind: I48}
			fieldSize = 6

		case discriminant == 6:
			fieldType = RecordFieldType{Kind: I64}
			fieldSize = 8

		case discriminant == 7:
			fieldType = RecordFieldType{Kind: Float}
			fieldSize = 8

		case discriminant == 8:
			fieldType = RecordFieldType{Kind: Zero}
			fieldSize = 0

		case discriminant == 9:
			fieldType = RecordFieldType{Kind: One}
			fieldSize = 0

		case discriminant >= 12 && discriminant%2 == 0:
			size := int((discriminant - 12) / 2)

			fieldType = RecordFieldType{
				Kind: Blob,
				Size: size,
			}

			fieldSize = size

		case discriminant >= 13 && discriminant%2 == 1:
			size := int((discriminant - 13) / 2)

			fieldType = RecordFieldType{
				Kind: String,
				Size: size,
			}

			fieldSize = size

		default:
			return nil, fmt.Errorf("unsupported field type: %d", discriminant)
		}

		fields = append(fields, RecordField{
			Offset:    uint(currentOffset),
			FieldType: fieldType,
		})

		currentOffset += fieldSize
	}

	return &RecordHeader{
		Fields: fields,
	}, nil
}
