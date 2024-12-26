package sccp

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

var mockCCs = [][]byte{
	{0x2, 0x0, 0x3, 0x75, 0x3, 0x20, 0x48, 0x2, 0x0},
	{0x2, 0x0, 0x70, 0x3e, 0x0, 0x0, 0x5, 0x2, 0x1, 0x3, 0x4, 0x43, 0x1c, 0x2d, 0xfe, 0x0},
}

func TestCC(t *testing.T) {
	for i, v := range mockCCs {
		cc, err := ParseCC(v)
		if err != nil {
			t.Fatal(i, err)
		}
		b, err := cc.MarshalBinary()
		if err != nil {
			t.Fatal(i, err)
		}
		if !bytes.Equal(v, b) {
			fmt.Println(hex.EncodeToString(v))
			fmt.Println(hex.EncodeToString(b))

			t.Fatal(i, err)
		}
		fmt.Println(cc)
	}
}
