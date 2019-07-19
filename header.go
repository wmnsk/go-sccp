// Copyright 2019 go-sccp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package sccp

import (
	"fmt"
)

// Header is a SCCP common header.
type Header struct {
	Type    uint8
	Payload []byte
}

// NewHeader creates a new Header.
func NewHeader(mtype uint8, payload []byte) *Header {
	return &Header{
		Type:    mtype,
		Payload: payload,
	}
}

// Serialize returns the byte sequence generated from a Header instance.
func (h *Header) Serialize() ([]byte, error) {
	b := make([]byte, h.Len())
	if err := h.SerializeTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// SerializeTo puts the byte sequence in the byte array given as b.
func (h *Header) SerializeTo(b []byte) error {
	b[0] = h.Type
	copy(b[1:h.Len()], h.Payload)
	return nil
}

// DecodeHeader decodes given byte sequence as a SCCP common header.
func DecodeHeader(b []byte) (*Header, error) {
	h := &Header{}
	if err := h.DecodeFromBytes(b); err != nil {
		return nil, err
	}
	return h, nil
}

// DecodeFromBytes sets the values retrieved from byte sequence in a SCCP common header.
func (h *Header) DecodeFromBytes(b []byte) error {
	l := len(b)
	if l < 2 {
		return ErrTooShortToDecode
	}
	h.Type = b[0]
	h.Payload = b[1:l]
	return nil
}

// Len returns the actual length of Header.
func (h *Header) Len() int {
	return 1 + len(h.Payload)
}

// String returns the SCCP common header values in human readable format.
func (h *Header) String() string {
	return fmt.Sprintf("{Type: %d, Payload: %x}",
		h.Type,
		h.Payload,
	)
}
