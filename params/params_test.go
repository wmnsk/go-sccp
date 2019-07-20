// Copyright 2019 go-sccp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params_test

import (
	"encoding"
	"testing"

	"github.com/pascaldekloe/goe/verify"
	"github.com/wmnsk/go-sccp/params"
)

type serializable interface {
	encoding.BinaryMarshaler
	MarshalLen() int
}

type decodeFunc func([]byte) (serializable, error)

var testcases = []struct {
	description string
	structured  serializable
	serialized  []byte
	decodeFunc
}{
	{
		description: "PartyAddress",
		structured: params.NewPartyAddress(
			0x12, 0, 6, 0, // GTI, SPC, SSN, TT
			1, 1, 4, // NP, ES, NAI
			[]byte{
				0x21, 0x43, 0x65, 0x87, 0x09,
			},
		),
		serialized: []byte{
			0x0a, 0x12, 0x06, 0x00, 0x11, 0x04, 0x21, 0x43, 0x65, 0x87, 0x09,
		},
		decodeFunc: func(b []byte) (serializable, error) {
			v, err := params.ParsePartyAddress(b)
			if err != nil {
				return nil, err
			}

			return v, nil
		},
	},
}

func TestStructuredParams(t *testing.T) {
	t.Helper()

	for _, c := range testcases {
		t.Run(c.description, func(t *testing.T) {
			t.Run("Decode", func(t *testing.T) {
				prm, err := c.decodeFunc(c.serialized)
				if err != nil {
					t.Fatal(err)
				}

				if got, want := prm, c.structured; !verify.Values(t, "", got, want) {
					t.Fail()
				}
			})

			t.Run("Serialize", func(t *testing.T) {
				b, err := c.structured.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}

				if got, want := b, c.serialized; !verify.Values(t, "", got, want) {
					t.Fail()
				}
			})

			t.Run("Len", func(t *testing.T) {
				if got, want := c.structured.MarshalLen(), len(c.serialized); got != want {
					t.Fatalf("got %v want %v", got, want)
				}
			})
		})
	}
}
