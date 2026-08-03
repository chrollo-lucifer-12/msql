package pager

import "fmt"

type PageType uint8

const (
	TableLeaf PageType = iota
)
const (
	PAGE_LEAF_HEADER_SIZE              uint = 8
	PAGE_FIRST_FREEBLOCK_OFFSET        uint = 1
	PAGE_CELL_COUNT_OFFSET             uint = 3
	PAGE_CELL_CONTENT_OFFSET           uint = 5
	PAGE_FRAGMENTED_BYTES_COUNT_OFFSET uint = 7
	PAGE_LEAF_TABLE_ID                      = 0x0D
)

type Page interface {
	isPage()
}

type PageHeader struct {
	PageType          PageType
	FirstFreeBlock    uint16 // offset of first free block
	CellCount         uint16 // no of cells in page
	CellContentOffset uint32 // offset of the first cell
	FragmentedBytes   uint8  // no of fragmented free bytes
}

type TableLeafCell struct {
	Size    int64 // varint
	RowID   int64 // varint
	Payload []byte
}

type TableLeafPage struct {
	Header       PageHeader
	CellPointers []uint16
	Cells        []TableLeafCell
}

func (TableLeafPage) isPage() {}

func decodePage(buffer []byte, pageNum int) (Page, error) {
	var ptrOffset uint16
	if pageNum == 1 {
		ptrOffset = uint16(HEADER_SIZE)
	}

	switch buffer[0] {
	case PAGE_LEAF_TABLE_ID:
		return decodeTableLeafPage(buffer, ptrOffset), nil

	default:
		return nil, fmt.Errorf("unknown page type: %d", buffer[0])
	}
}

func decodeTableLeafPage(buffer []byte, ptrOffset uint16) TableLeafPage {
	header := decodePageHeader(buffer)
	contentBuffer := buffer[PAGE_LEAF_HEADER_SIZE:]

	cellPointers := decodeCellPointers(contentBuffer, uint(header.CellCount), ptrOffset)

	cells := make([]TableLeafCell, 0, len(cellPointers))

	for _, ptr := range cellPointers {
		cell := decodeTableLeafCell(buffer[ptr:])
		cells = append(cells, cell)
	}

	return TableLeafPage{
		Header:       header,
		CellPointers: cellPointers,
		Cells:        cells,
	}
}

func decodePageHeader(buffer []byte) PageHeader {
	pageType := buffer[0]
	firstFreeBlock := readBEWordAt(buffer, int(PAGE_FIRST_FREEBLOCK_OFFSET))
	cellCount := readBEWordAt(buffer, int(PAGE_CELL_COUNT_OFFSET))
	cellContentOffset := readBEWordAt(buffer, int(PAGE_CELL_CONTENT_OFFSET))
	fragmentedBytesCount := readBEWordAt(buffer, int(PAGE_FRAGMENTED_BYTES_COUNT_OFFSET))

	return PageHeader{
		PageType:          PageType(pageType),
		FirstFreeBlock:    firstFreeBlock,
		CellCount:         cellCount,
		CellContentOffset: uint32(cellContentOffset),
		FragmentedBytes:   uint8(fragmentedBytesCount),
	}
}

func decodeCellPointers(buffer []byte, n uint, offset uint16) []uint16 {
	pointers := make([]uint16, 0, n)

	for i := 0; i < int(n); i++ {
		ptr := readBEWordAt(buffer, 2*i)
		pointers = append(pointers, offset-ptr)
	}

	return pointers
}

func decodeTableLeafCell(buffer []byte) TableLeafCell {
	n, size := decodeVarint(buffer, 0)
	buffer = buffer[n:]

	n, rowId := decodeVarint(buffer, 0)
	buffer = buffer[n:]

	payload := make([]byte, size)
	copy(payload, buffer[:size])

	return TableLeafCell{
		Size:    size,
		RowID:   rowId,
		Payload: payload,
	}
}

func decodeVarint(buffer []byte, offset int) (uint8, int64) {
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
