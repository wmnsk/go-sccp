// Copyright 2019 go-sccp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package sccp

import (
	"encoding/hex"
	"testing"
)

var testHeaderBytes = []byte{
	0x01, 0xde, 0xad, 0xbe, 0xef,
}

func TestSerializeHeader(t *testing.T) {
	h := NewHeader(
		1, // Code
		[]byte{
			0xde, 0xad, 0xbe, 0xef,
		},
	)

	serialized, err := h.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize Header %s", err)
	}
	for i, s := range serialized {
		x := testHeaderBytes[i]
		if s != x {
			t.Errorf("Bytes doesn't match. Expected: %#x, Got: %#x at %dth", x, s, i)
		}
	}
	t.Logf("%x", serialized)
}

func TestDecodeHeader(t *testing.T) {
	h, err := DecodeHeader(testHeaderBytes)
	if err != nil {
		t.Fatalf("Failed to decode Header: %s", err)
	}

	dummyStr := hex.EncodeToString(h.Payload)
	switch {
	case h.Type != 1:
		t.Errorf("Type doesn't match. Expected: %d, Got: %d", 1, h.Type)
	case dummyStr != "deadbeef":
		t.Errorf("Payload doesn't match. Expected: %s, Got: %s", "deadbeef", dummyStr)
	}
	t.Log(h.String())
}
