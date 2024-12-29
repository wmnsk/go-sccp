// Copyright 2019-2024 go-sccp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

import (
	"testing"
)

var cases = []struct {
	description   string
	typed         ProtocolClass
	bin           uint8
	class         int
	returnOnError bool
}{
	{
		"Class 1, No ReturnOnError",
		NewProtocolClass(1, false),
		0x01,
		1,
		false,
	},
	{
		"Class 2, No ReturnOnError",
		NewProtocolClass(2, false),
		0x02,
		2,
		false,
	},
	{
		"Class 3, No ReturnOnError",
		NewProtocolClass(3, false),
		0x03,
		3,
		false,
	},
	{
		"Class 4, No ReturnOnError",
		NewProtocolClass(4, false),
		0x04,
		4,
		false,
	},
	{
		"Class 1, ReturnOnError",
		NewProtocolClass(1, true),
		0x81,
		1,
		true,
	},
	{
		"Class 2, ReturnOnError",
		NewProtocolClass(2, true),
		0x82,
		2,
		true,
	},
	{
		"Class 3, ReturnOnError",
		NewProtocolClass(3, true),
		0x83,
		3,
		true,
	},
	{
		"Class 4, ReturnOnError",
		NewProtocolClass(4, true),
		0x84,
		4,
		true,
	},
}

func TestProtocolClass(t *testing.T) {
	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			if got, want := uint8(c.typed), c.bin; got != want {
				t.Errorf("unexpected ProtocolClass: got: %v, want: %v", got, want)
			}

			if got, want := c.typed.Class(), c.class; got != want {
				t.Errorf("unexpected Class: got: %v, want: %v", got, want)
			}

			if got, want := c.typed.ReturnOnError(), c.returnOnError; got != want {
				t.Errorf("unexpected ReturnOnError: got: %v, want: %v", got, want)
			}
		})
	}
}
