// Copyright 2019-2024 go-sccp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

import (
	"testing"
)

var testProtocolClassBytes = []byte{
	0x01, 0x02, 0x03, 0x04,
	0x81, 0x82, 0x83, 0x84,
}

func TestNewProtocolClass(t *testing.T) {
	protocolClasses := []ProtocolClass{
		NewProtocolClass(1, false),
		NewProtocolClass(2, false),
		NewProtocolClass(3, false),
		NewProtocolClass(4, false),
		NewProtocolClass(1, true),
		NewProtocolClass(2, true),
		NewProtocolClass(3, true),
		NewProtocolClass(4, true),
	}
	for i, p := range protocolClasses {
		x := testProtocolClassBytes[i]
		if uint8(p) != x {
			t.Errorf("Bytes doesn't match. Want: %#x, Got: %#x at %dth", x, p, i)
		}
	}
}

func TestClass(t *testing.T) {
	protocolClasses := []ProtocolClass{
		NewProtocolClass(1, false),
		NewProtocolClass(2, false),
		NewProtocolClass(3, false),
		NewProtocolClass(4, false),
		NewProtocolClass(1, true),
		NewProtocolClass(2, true),
		NewProtocolClass(3, true),
		NewProtocolClass(4, true),
	}
	for i, p := range protocolClasses {
		if i < 4 {
			if got, want := p.Class(), i+1; got != want {
				t.Errorf("Class doesn't match. Want: %d, Got: %d when ProtocolClass is %x", want, got, p)
			}
		} else if i >= 4 {
			if got, want := p.Class(), i-3; got != want {
				t.Errorf("Class doesn't match. Want: %d, Got: %d when ProtocolClass is %x", want, got, p)
			}
		}
	}
}

func TestReturnOnError(t *testing.T) {
	protocolClasses := []ProtocolClass{
		NewProtocolClass(1, false),
		NewProtocolClass(2, false),
		NewProtocolClass(3, false),
		NewProtocolClass(4, false),
		NewProtocolClass(1, true),
		NewProtocolClass(2, true),
		NewProtocolClass(3, true),
		NewProtocolClass(4, true),
	}
	for i, p := range protocolClasses {
		if i < 4 {
			if got := p.ReturnOnError(); !got {
				t.Errorf("ReturnOnError doesn't match. Want: %v, Got: %v when ProtocolClass is %x", false, got, p)
			}
		} else if i >= 4 {
			if got := p.ReturnOnError(); got {
				t.Errorf("ReturnOnError doesn't match. Want: %v, Got: %v when ProtocolClass is %x", true, got, p)
			}
		}
	}
}
