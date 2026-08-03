package cursor

import (
	"github.com/msql/pager"
	"github.com/msql/utils"
)

type Cursor struct {
	header    *RecordHeader
	pager     *pager.Pager
	pageIndex uint
	pageCell  uint
}

func (c *Cursor) Field(n int) *Value {
	if n >= len(c.header.Fields) {
		return nil
	}

	recordField := c.header.Fields[n]

	page, err := c.pager.ReadPage(c.pageIndex)
	if err != nil {
		return nil
	}

	leaf, ok := page.(*pager.TableLeafPage)
	if !ok {
		return nil
	}

	payload := leaf.Cells[c.pageCell].Payload

	switch recordField.FieldType.Kind {
	case NULL:
		v := NewNull()
		return &v

	case I8:
		v := NewInt(utils.ReadI8At(payload, int(recordField.Offset)))
		return &v

	case I16:
		v := NewInt(utils.ReadI16At(payload, int(recordField.Offset)))
		return &v

	case I24:
		v := NewInt(utils.ReadI24At(payload, int(recordField.Offset)))
		return &v

	case I32:
		v := NewInt(utils.ReadI32At(payload, int(recordField.Offset)))
		return &v

	case I48:
		v := NewInt(utils.ReadI48At(payload, int(recordField.Offset)))
		return &v

	case I64:
		v := NewInt(utils.ReadI64At(payload, int(recordField.Offset)))
		return &v

	case Float:
		v := NewFloat(utils.ReadF64At(payload, int(recordField.Offset)))
		return &v

	case String:
		start := int(recordField.Offset)
		end := start + recordField.FieldType.Size

		v := NewString(string(payload[start:end]))
		return &v

	case Blob:
		start := int(recordField.Offset)
		end := start + recordField.FieldType.Size

		v := NewBlob(payload[start:end])
		return &v

	default:
		panic("unimplemented")
	}
}
