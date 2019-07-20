// Copyright 2019 go-sccp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package sccp_test

import (
	"encoding"
	"testing"

	"github.com/pascaldekloe/goe/verify"
	"github.com/wmnsk/go-sccp"
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
		description: "Header",
		structured: sccp.NewHeader(
			1, // Code
			[]byte{0xde, 0xad, 0xbe, 0xef},
		),
		serialized: []byte{
			0x01, 0xde, 0xad, 0xbe, 0xef,
		},
		decodeFunc: func(b []byte) (serializable, error) {
			v, err := sccp.ParseHeader(b)
			if err != nil {
				return nil, err
			}

			return v, nil
		},
	}, {
		description: "UDT",
		structured: sccp.NewUDT(
			1,    // Protocol Class
			true, // Message handling
			params.NewPartyAddress( // CalledPartyAddress
				0x12, 0, 6, 0x00, // Indicator, SPC, SSN, TT
				0x01, 0x01, 0x04, // NP, ES, NAI
				[]byte{ // Data
					0x21, 0x43, 0x65, 0x87, 0x09, 0x21, 0x43, 0x65,
				},
			),
			params.NewPartyAddress( // CallingPartyAddress
				0x12, 0, 7, 0x00, // Indicator, SPC, SSN, TT
				0x01, 0x02, 0x04, // NP, ES, NAI
				[]byte{ // Data
					0x89, 0x67, 0x45, 0x23, 0x01,
				},
			),
			[]byte{ // Data
				0xde, 0xad, 0xbe, 0xef,
			},
		),
		serialized: []byte{
			0x09, 0x81, 0x03, 0x10, 0x1a, 0x0d, 0x12, 0x06,
			0x00, 0x11, 0x04, 0x21, 0x43, 0x65, 0x87, 0x09,
			0x21, 0x43, 0x65, 0x0a, 0x12, 0x07, 0x00, 0x12,
			0x04, 0x89, 0x67, 0x45, 0x23, 0x01, 0x04, 0xde,
			0xad, 0xbe, 0xef,
		},
		decodeFunc: func(b []byte) (serializable, error) {
			v, err := sccp.ParseUDT(b)
			if err != nil {
				return nil, err
			}

			return v, nil
		},
	},
}

func TestMessages(t *testing.T) {
	t.Helper()

	for _, c := range testcases {
		t.Run(c.description, func(t *testing.T) {
			t.Run("Decode", func(t *testing.T) {
				msg, err := c.decodeFunc(c.serialized)
				if err != nil {
					t.Fatal(err)
				}

				if got, want := msg, c.structured; !verify.Values(t, "", got, want) {
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

			t.Run("Interface", func(t *testing.T) {
				// Ignore *Header and Generic in this tests.
				if _, ok := c.structured.(*sccp.Header); ok {
					return
				}

				decoded, err := sccp.ParseMessage(c.serialized)
				if err != nil {
					t.Fatal(err)
				}

				if got, want := decoded.MessageType(), c.structured.(sccp.Message).MessageType(); got != want {
					t.Fatalf("got %v want %v", got, want)
				}
				if got, want := decoded.MessageTypeName(), c.structured.(sccp.Message).MessageTypeName(); got != want {
					t.Fatalf("got %v want %v", got, want)
				}
			})
		})
	}
}
