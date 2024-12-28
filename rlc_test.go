package sccp

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

var mockRLC = []byte{0x5, 0x0, 0x70, 0x3e, 0x0, 0x0, 0x5}

func TestRLC(t *testing.T) {
	r, err := ParseRLC(mockRLC)
	if err != nil {
		t.Fatal(err)
	}
	b, err := r.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(mockRLC, b) {
		fmt.Println(hex.EncodeToString(mockRLC))
		fmt.Println(hex.EncodeToString(b))

		t.Fatal(err)
	}
	fmt.Println(r)
}
