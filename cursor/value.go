package cursor

import "strconv"

type ValueType int

const (
	NullValue ValueType = iota
	StringValue
	BlobValue
	IntValue
	FloatValue
)

type Value struct {
	Type ValueType

	Str   string
	Blob  []byte
	Int   int64
	Float float64
}

func NewNull() Value {
	return Value{
		Type: NullValue,
	}
}

func NewString(s string) Value {
	return Value{
		Type: StringValue,
		Str:  s,
	}
}

func NewBlob(b []byte) Value {
	blob := make([]byte, len(b))
	copy(blob, b)

	return Value{
		Type: BlobValue,
		Blob: blob,
	}
}

func NewInt(v int64) Value {
	return Value{
		Type: IntValue,
		Int:  v,
	}
}

func NewFloat(v float64) Value {
	return Value{
		Type:  FloatValue,
		Float: v,
	}
}

func (v Value) AsString() (string, bool) {
	if v.Type != StringValue {
		return "", false
	}
	return v.Str, true
}

func (v Value) AsBlob() ([]byte, bool) {
	if v.Type != BlobValue {
		return nil, false
	}
	return v.Blob, true
}

func (v Value) AsInt() (int64, bool) {
	if v.Type != IntValue {
		return 0, false
	}
	return v.Int, true
}

func (v Value) AsFloat() (float64, bool) {
	if v.Type != FloatValue {
		return 0, false
	}
	return v.Float, true
}

func (v Value) IsNull() bool {
	return v.Type == NullValue
}

func (v Value) String() string {
	switch v.Type {
	case NullValue:
		return "NULL"
	case StringValue:
		return v.Str
	case BlobValue:
		return string(v.Blob)
	case IntValue:
		return strconv.FormatInt(v.Int, 10)
	case FloatValue:
		return strconv.FormatFloat(v.Float, 'g', -1, 64)
	default:
		return "<unknown>"
	}
}
