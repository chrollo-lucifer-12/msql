package db

import (
	"fmt"
	"io"
	"os"

	"github.com/msql/cursor"
	"github.com/msql/pager"
)

type Db struct {
	Header pager.DbHeader
	pager  *pager.Pager
}

func (d *Db) FromFile(filename string) (*Db, error) {
	file, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	headerBuffer := make([]byte, pager.HEADER_SIZE)

	_, err = io.ReadFull(file, headerBuffer)
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("read db header: %w", err)
	}

	header, err := pager.ParseHeader(headerBuffer)
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("parse db header: %w", err)
	}

	p := pager.NewPager(file, uint(header.PageSize))

	return &Db{
		Header: header,
		pager:  p,
	}, nil
}

func (d *Db) GetScanner(page uint) *cursor.Scanner {
	return cursor.NewScanner(d.pager, page)
}
