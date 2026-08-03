package pager

import (
	"os"
	"testing"
)

func TestParseHeader(t *testing.T) {

	data, err := os.ReadFile("../temp/test.db")
	if err != nil {
		t.Fatalf("failed to read test db: %v", err)
	}

	header, err := ParseHeader(data)
	if err != nil {
		t.Fatalf("failed to parse header: %v", err)
	}

	if header.PageSize != 4096 {
		t.Errorf("PageSize = %d, want %d", header.PageSize, 4096)
	}
}
