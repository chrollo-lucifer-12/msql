package cursor

import (
	"fmt"

	"github.com/msql/pager"
)

type Scanner struct {
	pager *pager.Pager
	page  uint
	cell  uint
}

func NewScanner(pager *pager.Pager, page uint) *Scanner {
	return &Scanner{
		pager: pager,
		page:  page,
		cell:  0,
	}
}

func (s *Scanner) NextRecord() (*Cursor, error) {
	page, err := s.pager.ReadPage(s.page)
	if err != nil {
		return nil, err
	}

	leaf, ok := page.(*pager.TableLeafPage)
	if !ok {
		return nil, fmt.Errorf("wrong page type")
	}

	if int(s.cell) >= len(leaf.Cells) {
		return nil, nil
	}

	cell := leaf.Cells[s.cell]
	header, err := parseRecordHeader(cell.Payload)
	if err != nil {
		return nil, err
	}

	c := &Cursor{
		header:    header,
		pager:     s.pager,
		pageIndex: s.page,
		pageCell:  s.cell,
	}

	s.cell++

	return c, nil

}
