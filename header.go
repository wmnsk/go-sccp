// Copyright 2019-2023 go-sccp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package sccp

import (
	"fmt"
	"io"
)

// Header is a SCCP common header.
type Header struct {
	Type    MsgType
	Payload []byte
}

// NewHeader creates a new Header.
func NewHeader(mtype MsgType, payload []byte) *Header {
	return &Header{
		Type:    mtype,
		Payload: payload,
	}
}

// MarshalBinary returns the byte sequence generated from a Header instance.
func (h *Header) MarshalBinary() ([]byte, error) {
	b := make([]byte, h.MarshalLen())
	if err := h.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (h *Header) MarshalTo(b []byte) error {
	b[0] = uint8(h.Type)
	copy(b[1:h.MarshalLen()], h.Payload)
	return nil
}

// ParseHeader decodes given byte sequence as a SCCP common header.
func ParseHeader(b []byte) (*Header, error) {
	h := &Header{}
	if err := h.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return h, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a SCCP common header.
func (h *Header) UnmarshalBinary(b []byte) error {
	if len(b) < 2 {
		return io.ErrUnexpectedEOF
	}
	h.Type = MsgType(b[0])
	h.Payload = b[1:]
	return nil
}

// MarshalLen returns the serial length.
func (h *Header) MarshalLen() int {
	return 1 + len(h.Payload)
}

// String returns the SCCP common header values in human readable format.
func (h *Header) String() string {
	return fmt.Sprintf("{Type: %d, Payload: %x}",
		h.Type,
		h.Payload,
	)
}
