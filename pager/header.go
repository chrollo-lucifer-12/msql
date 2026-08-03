package pager

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type DbHeader struct {
	PageSize uint32
}

const (
	HEADER_SIZE             uint   = 100
	HEADER_PREFIX                  = "SQLite format 3\x00"
	HEADER_PAGE_SIZE_OFFSET uint32 = 16
	PAGE_MAX_SIZE           uint32 = (1 << 16)
)

func errInvalidHeaderPrefix(header string) error {
	return fmt.Errorf("invalid header prefix: %s", header)
}

func ParseHeader(buffer []byte) (DbHeader, error) {
	if strings.HasPrefix(string(buffer), HEADER_PREFIX) == false {
		return DbHeader{}, errInvalidHeaderPrefix(string(buffer)[0:HEADER_PAGE_SIZE_OFFSET])
	}

	pageSizeRaw := readBEWordAt(buffer, int(HEADER_PAGE_SIZE_OFFSET))

	var pageSize uint32

	switch {
	case pageSizeRaw == 1:
		pageSize = PAGE_MAX_SIZE
	case pageSizeRaw != 0 && (pageSizeRaw&(pageSizeRaw-1)) == 0:
		pageSize = uint32(pageSizeRaw)
	default:
		return DbHeader{}, fmt.Errorf("page size is not a power of 2: %d", pageSizeRaw)
	}

	return DbHeader{
		PageSize: pageSize,
	}, nil
}

func readBEWordAt(buffer []byte, offset int) uint16 {
	return binary.BigEndian.Uint16(buffer[offset:])
}
