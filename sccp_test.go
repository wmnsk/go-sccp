// Copyright 2019-2024 go-sccp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package sccp_test

import (
	"encoding"
	"io"
	"strings"
	"testing"

	"github.com/pascaldekloe/goe/verify"
	"github.com/wmnsk/go-sccp"
	"github.com/wmnsk/go-sccp/params"
)

type serializable interface {
	encoding.BinaryMarshaler
	MarshalTo([]byte) error
	MarshalLen() int
}

var testcases = []struct {
	description string
	structured  serializable
	serialized  []byte
	parseFunc   func([]byte) (serializable, error)
}{
	{
		description: "UDT",
		structured: sccp.NewUDT(
			1,    // Protocol Class
			true, // Message handling
			params.NewCalledPartyAddress(
				params.NewAddressIndicator(false, true, false, params.GTITTNPESNAI),
				0, 6, // SPC, SSN
				params.NewGlobalTitle(
					params.GTITTNPESNAI,
					params.TranslationType(0),
					params.NPISDNTelephony,
					params.ESBCDOdd,
					params.NAIInternationalNumber,
					[]byte{0x21, 0x43, 0x65, 0x87, 0x09, 0x21, 0x43, 0x65},
				),
			),
			params.NewCallingPartyAddress(
				params.NewAddressIndicator(false, true, false, params.GTITTNPESNAI),
				0, 7, // SPC, SSN
				params.NewGlobalTitle(
					params.GTITTNPESNAI,
					params.TranslationType(0),
					params.NPISDNTelephony,
					params.ESBCDEven,
					params.NAIInternationalNumber,
					[]byte{0x89, 0x67, 0x45, 0x23, 0x01},
				),
			),
			[]byte{0xde, 0xad, 0xbe, 0xef},
		),
		serialized: []byte{
			0x09,
			0x81,
			0x03, 0x10, 0x1a,
			0x0d, 0x12, 0x06, 0x00, 0x11, 0x04, 0x21, 0x43, 0x65, 0x87, 0x09, 0x21, 0x43, 0x65,
			0x0a, 0x12, 0x07, 0x00, 0x12, 0x04, 0x89, 0x67, 0x45, 0x23, 0x01,
			0x04, 0xde, 0xad, 0xbe, 0xef,
		},
		parseFunc: func(b []byte) (serializable, error) {
			return sccp.ParseUDT(b)
		},
	},
	{
		description: "UDT-2Bytes-PartyAddress",
		structured: sccp.NewUDT(
			1,    // Protocol Class
			true, // Message handling
			params.NewCalledPartyAddress(0x42, 0, 6, nil),
			params.NewCallingPartyAddress(0x42, 0, 7, nil),
			nil,
		),
		serialized: []byte{
			0x09, 0x81, 0x03, 0x05, 0x07, 0x02, 0x42, 0x06, 0x02, 0x42, 0x07, 0x00,
		},
		parseFunc: func(b []byte) (serializable, error) {
			return sccp.ParseUDT(b)
		},
	},
	{
		description: "XUDT/No optionals",
		structured: sccp.NewXUDT(
			1,    // Protocol Class
			true, // Message handling
			2,    // Hop Counter
			params.NewCalledPartyAddress(
				params.NewAddressIndicator(false, true, false, params.GTITTNPESNAI),
				0, 6, // SPC, SSN
				params.NewGlobalTitle(
					params.GTITTNPESNAI,
					params.TranslationType(0),
					params.NPISDNTelephony,
					params.ESBCDOdd,
					params.NAIInternationalNumber,
					[]byte{0x21, 0x43, 0x65, 0x87, 0x09, 0x21, 0x43, 0x65},
				),
			),
			params.NewCallingPartyAddress(
				params.NewAddressIndicator(false, true, false, params.GTITTNPESNAI),
				0, 7, // SPC, SSN
				params.NewGlobalTitle(
					params.GTITTNPESNAI,
					params.TranslationType(0),
					params.NPISDNTelephony,
					params.ESBCDEven,
					params.NAIInternationalNumber,
					[]byte{0x89, 0x67, 0x45, 0x23, 0x01},
				),
			),
			[]byte{0xde, 0xad, 0xbe, 0xef},
		),
		serialized: []byte{
			0x11,                   // MsgType
			0x81,                   // Protocol Class
			0x02,                   // Hop Counter
			0x04, 0x11, 0x1b, 0x00, // Pointers
			0x0d, 0x12, 0x06, 0x00, 0x11, 0x04, 0x21, 0x43, 0x65, 0x87, 0x09, 0x21, 0x43, 0x65, // CdPA
			0x0a, 0x12, 0x07, 0x00, 0x12, 0x04, 0x89, 0x67, 0x45, 0x23, 0x01, // CgPA
			0x04, 0xde, 0xad, 0xbe, 0xef, // Data
		},
		parseFunc: func(b []byte) (serializable, error) {
			return sccp.ParseXUDT(b)
		},
	},
	{
		description: "XUDT/with optionals",
		structured: sccp.NewXUDT(
			1,    // Protocol Class
			true, // Message handling
			2,    // Hop Counter
			params.NewCalledPartyAddress(
				params.NewAddressIndicator(false, true, false, params.GTITTNPESNAI),
				0, 6, // SPC, SSN
				params.NewGlobalTitle(
					params.GTITTNPESNAI,
					params.TranslationType(0),
					params.NPISDNTelephony,
					params.ESBCDOdd,
					params.NAIInternationalNumber,
					[]byte{0x21, 0x43, 0x65, 0x87, 0x09, 0x21, 0x43, 0x65},
				),
			),
			params.NewCallingPartyAddress(
				params.NewAddressIndicator(false, true, false, params.GTITTNPESNAI),
				0, 7, // SPC, SSN
				params.NewGlobalTitle(
					params.GTITTNPESNAI,
					params.TranslationType(0),
					params.NPISDNTelephony,
					params.ESBCDEven,
					params.NAIInternationalNumber,
					[]byte{0x89, 0x67, 0x45, 0x23, 0x01},
				),
			),
			[]byte{0xde, 0xad, 0xbe, 0xef},
			params.NewSegmentation(true, 1, 2, 0xffffff),
			params.NewImportance(2),
		),
		serialized: []byte{
			0x11,                   // MsgType
			0x81,                   // Protocol Class
			0x02,                   // Hop Counter
			0x04, 0x11, 0x1b, 0x1f, // Pointers
			0x0d, 0x12, 0x06, 0x00, 0x11, 0x04, 0x21, 0x43, 0x65, 0x87, 0x09, 0x21, 0x43, 0x65, // CdPA
			0x0a, 0x12, 0x07, 0x00, 0x12, 0x04, 0x89, 0x67, 0x45, 0x23, 0x01, // CgPA
			0x04, 0xde, 0xad, 0xbe, 0xef, // Data
			0x10, 0x04, 0xc2, 0xff, 0xff, 0xff, // Segmentation
			0x12, 0x01, 0x02, // Importance
			0x00, // End of optional parameters
		},
		parseFunc: func(b []byte) (serializable, error) {
			return sccp.ParseXUDT(b)
		},
	},
	{
		description: "SCMG SSA",
		structured:  sccp.NewSCMG(sccp.SCMGTypeSSA, 9, 405, 0, 0),
		serialized:  []byte{0x1, 0x09, 0x95, 0x01, 0x00},
		parseFunc: func(b []byte) (serializable, error) {
			return sccp.ParseSCMG(b)
		},
	},
	{
		description: "SCMG SSC",
		structured:  sccp.NewSCMG(sccp.SCMGTypeSSC, 9, 405, 0, 4),
		serialized:  []byte{0x6, 0x09, 0x95, 0x01, 0x00, 0x04},
		parseFunc: func(b []byte) (serializable, error) {
			return sccp.ParseSCMG(b)
		},
	},
}

func TestMessages(t *testing.T) {
	t.Helper()

	for _, c := range testcases {
		t.Run(c.description, func(t *testing.T) {
			t.Run("Decode", func(t *testing.T) {
				msg, err := c.parseFunc(c.serialized)
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
				if _, ok := c.structured.(*sccp.SCMG); ok {
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

func TestPartialStructuredMessages(t *testing.T) {
	for _, c := range testcases {
		if strings.Contains(c.description, "SCMG") {
			continue
		}
		for i := range c.serialized {
			partial := c.serialized[:i]
			_, err := c.parseFunc(partial)
			if err != io.ErrUnexpectedEOF {
				t.Errorf("parse %v / %#x: got error %v, want unexpected EOF", c.description, partial, err)
			}
		}

		for i := range c.serialized {
			if i == len(c.serialized) {
				continue
			}
			b := make([]byte, i)
			if err := c.structured.MarshalTo(b); err != io.ErrUnexpectedEOF {
				t.Errorf("marshal %v / %#x: got error %v, want unexpected EOF", c.description, b, err)
			}
		}
	}
}
