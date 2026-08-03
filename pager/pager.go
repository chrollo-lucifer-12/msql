package pager

import (
	"fmt"
	"io"
	"os"
)

type Pager struct {
	pages    map[uint]Page
	pageSize uint
	input    *os.File
}

func NewPager(input *os.File, pageSize uint) *Pager {
	return &Pager{
		pageSize: pageSize,
		input:    input,
		pages:    make(map[uint]Page),
	}
}

func (p *Pager) ReadPage(n uint) (Page, error) {
	if page, ok := p.pages[n]; ok == true {
		return page, nil
	}

	page, err := p.loadPage(n)
	if err != nil {
		return nil, fmt.Errorf("load page: %w", err)
	}

	return page, nil
}

func (p *Pager) loadPage(n uint) (Page, error) {
	offset := (n - 1) * p.pageSize

	_, err := p.input.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("seek to page start :%w", err)
	}

	buffer := make([]byte, p.pageSize)

	_, err = io.ReadFull(p.input, buffer)
	if err != nil {
		return nil, fmt.Errorf("read page: %w", err)
	}

	page, err := decodePage(buffer, int(n))
	if err != nil {
		return nil, fmt.Errorf("decode page: %w", err)
	}

	return page, nil
}
