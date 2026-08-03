package pager

import (
	"fmt"

	"github.com/msql/utils"
)

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

	PAGE_LEAF_TABLE_ID = 0x0D
)

type Page interface {
	isPage()
}

type PageHeader struct {
	PageType          PageType
	FirstFreeBlock    uint16
	CellCount         uint16
	CellContentOffset uint32
	FragmentedBytes   uint8
}

type TableLeafCell struct {
	Size    int64
	RowID   int64
	Payload []byte
}

type TableLeafPage struct {
	Header       PageHeader
	CellPointers []uint16
	Cells        []TableLeafCell
}

func (TableLeafPage) isPage() {}

func decodePage(buffer []byte, pageNum int) (Page, error) {
	var pageOffset uint
	if pageNum == 1 {
		pageOffset = HEADER_SIZE
	}
	if len(buffer) <= int(pageOffset) {
		return nil, fmt.Errorf("page too small")
	}
	switch buffer[pageOffset] {
	case PAGE_LEAF_TABLE_ID:
		page := decodeTableLeafPage(buffer, pageOffset)
		return &page, nil
	default:
		return nil, fmt.Errorf(
			"unknown page type: %d",
			buffer[pageOffset],
		)
	}
}

func decodeTableLeafPage(buffer []byte, pageOffset uint) TableLeafPage {
	page := buffer[pageOffset:]
	header := decodePageHeader(page)
	cellPointers := decodeCellPointers(
		page[PAGE_LEAF_HEADER_SIZE:],
		uint(header.CellCount),
	)
	cells := make([]TableLeafCell, 0, len(cellPointers))
	for _, ptr := range cellPointers {
		cellOffset := uint(ptr)
		cell := decodeTableLeafCell(buffer[cellOffset:])
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

	firstFreeBlock := utils.ReadBEWordAt(
		buffer,
		int(PAGE_FIRST_FREEBLOCK_OFFSET),
	)

	cellCount := utils.ReadBEWordAt(
		buffer,
		int(PAGE_CELL_COUNT_OFFSET),
	)

	cellContentOffset := utils.ReadBEWordAt(
		buffer,
		int(PAGE_CELL_CONTENT_OFFSET),
	)

	fragmentedBytes := buffer[PAGE_FRAGMENTED_BYTES_COUNT_OFFSET]

	return PageHeader{
		PageType:          PageType(pageType),
		FirstFreeBlock:    firstFreeBlock,
		CellCount:         cellCount,
		CellContentOffset: uint32(cellContentOffset),
		FragmentedBytes:   fragmentedBytes,
	}
}

func decodeCellPointers(buffer []byte, n uint) []uint16 {
	pointers := make([]uint16, 0, n)

	for i := 0; i < int(n); i++ {
		ptr := utils.ReadBEWordAt(
			buffer,
			i*2,
		)

		pointers = append(pointers, ptr)
	}

	return pointers
}

func decodeTableLeafCell(buffer []byte) TableLeafCell {
	n, payloadSize := utils.DecodeVarint(buffer, 0)
	buffer = buffer[n:]

	n, rowID := utils.DecodeVarint(buffer, 0)
	buffer = buffer[n:]

	payload := make([]byte, payloadSize)
	copy(payload, buffer[:payloadSize])

	return TableLeafCell{
		Size:    payloadSize,
		RowID:   rowID,
		Payload: payload,
	}
}
