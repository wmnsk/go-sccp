// Copyright 2019 go-sccp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

/*
Package sccp provides encoding/decoding feature of Signalling Connection Control Part used in SS7/SIGTRAN protocol stack.

This is still an experimental project, and currently in its very early stage of development. Any part of implementations
(including exported APIs) may be changed before released as v1.0.0.
*/
package sccp

import (
	"encoding"
	"fmt"

	"github.com/pkg/errors"
)

// MsgType is type of SCCP message.
type MsgType uint8

// Message Type definitions.
const (
	_ MsgType = iota
	MsgTypeCR
	MsgTypeCC
	MsgTypeCREF
	MsgTypeRLSD
	MsgTypeRLC
	MsgTypeDT1
	MsgTypeDT2
	MsgTypeAK
	MsgTypeUDT
	MsgTypeUDTS
	MsgTypeED
	MsgTypeEA
	MsgTypeRSR
	MsgTypeRSC
	MsgTypeERR
	MsgTypeIT
	MsgTypeXUDT
	MsgTypeXUDTS
	MsgTypeLUDT
	MsgTypeLUDTS
)

// Message is an interface that defines SCCP messages.
type Message interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	MarshalTo([]byte) error
	MarshalLen() int
	MessageType() MsgType
	MessageTypeName() string
	fmt.Stringer
}

// FormatMessage returns the byte sequence generated from Message by Message Type.
func FormatMessage(m Message) ([]byte, error) {
	b := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// ParseMessage decodes the byte sequence into Message by Message Type.
// Currently this only supports UDT type of message only.
func ParseMessage(b []byte) (Message, error) {
	var m Message
	switch MsgType(b[0]) {
	/* TODO: implement!
	case CR:
	case CC:
	case CREF:
	case RLSD:
	case RLC:
	case DT1:
	case DT2:
	case AK:
	*/
	case MsgTypeUDT:
		m = &UDT{}
	/* TODO: implement!
	case UDTS:
	case ED:
	case EA:
	case RSR:
	case RSC:
	case ERR:
	case IT:
	case XUDT:
	case XUDTS:
	case LUDT:
	case LUDTS:
	*/
	default:
		return nil, UnsupportedTypeError(b[0])
	}

	if err := m.UnmarshalBinary(b); err != nil {
		return nil, errors.Wrap(err, "failed to decode SCCP")
	}
	return m, nil
}
