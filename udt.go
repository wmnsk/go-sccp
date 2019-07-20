// Copyright 2019 go-sccp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package sccp

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/wmnsk/go-sccp/params"
)

// UDT is SCCP Message Unit Data(UDT)
type UDT struct {
	Type uint8
	params.ProtocolClass
	Ptr1, Ptr2, Ptr3    uint8
	CalledPartyAddress  *params.PartyAddress
	CallingPartyAddress *params.PartyAddress
	DataLength          uint8
	Data                []byte
}

// NewUDT creates a new UDT.
func NewUDT(pcls int, mhandle bool, cdpa, cgpa *params.PartyAddress, data []byte) *UDT {
	u := &UDT{
		Type: 0x09,
		ProtocolClass: params.NewProtocolClass(
			pcls,
			mhandle,
		),
		Ptr1:                3,
		CalledPartyAddress:  cdpa,
		CallingPartyAddress: cgpa,
		Data:                data,
	}
	u.Ptr2 = u.Ptr1 + uint8(cdpa.Length)
	u.Ptr3 = u.Ptr2 + uint8(cgpa.Length)
	u.SetLength()

	return u
}

// MarshalBinary returns the byte sequence generated from a UDT instance.
func (u *UDT) MarshalBinary() ([]byte, error) {
	b := make([]byte, u.MarshalLen())
	if err := u.MarshalTo(b); err != nil {
		return nil, errors.Wrap(err, "failed to serialize UDT:")
	}

	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
// SCCP is dependent on the Pointers when serializing, which means that it might fail when invalid Pointers are set.
func (u *UDT) MarshalTo(b []byte) error {
	if len(b) < 5 {
		return io.ErrUnexpectedEOF
	}

	b[0] = u.Type
	b[1] = uint8(u.ProtocolClass)
	b[2] = u.Ptr1
	b[3] = u.Ptr2
	b[4] = u.Ptr3
	if err := u.CalledPartyAddress.MarshalTo(b[5:int(u.Ptr2+3)]); err != nil {
		return err
	}
	if err := u.CallingPartyAddress.MarshalTo(b[int(u.Ptr2+3):int(u.Ptr3+4)]); err != nil {
		return err
	}
	b[u.Ptr3+4] = u.DataLength
	copy(b[int(u.Ptr3+5):u.MarshalLen()], u.Data)

	return nil
}

// ParseUDT decodes given byte sequence as a SCCP UDT.
func ParseUDT(b []byte) (*UDT, error) {
	u := &UDT{}
	if err := u.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return u, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a SCCP UDT.
func (u *UDT) UnmarshalBinary(b []byte) error {
	l := len(b)
	if l < 4 {
		return io.ErrUnexpectedEOF
	}

	u.Type = b[0]
	u.ProtocolClass = params.ProtocolClass(b[1])
	u.Ptr1 = b[2]
	u.Ptr2 = b[3]
	u.Ptr3 = b[4]

	var err error
	u.CalledPartyAddress, err = params.ParsePartyAddress(b[5:int(u.Ptr2+3)])
	if err != nil {
		return errors.Wrap(err, "failed to decode CalledPartyAddress:")
	}

	u.CallingPartyAddress, err = params.ParsePartyAddress(b[int(u.Ptr2+3):int(u.Ptr3+4)])
	if err != nil {
		return errors.Wrap(err, "failed to decode CallingPartyAddress:")
	}

	u.DataLength = b[int(u.Ptr3+4)]
	u.Data = b[int(u.Ptr3+5):l]

	return nil
}

// MarshalLen returns the serial length.
func (u *UDT) MarshalLen() int {
	return 6 + u.CalledPartyAddress.MarshalLen() + u.CallingPartyAddress.MarshalLen() + len(u.Data)
}

// SetLength sets the length in Length field.
func (u *UDT) SetLength() {
	u.DataLength = uint8(len(u.Data))
}

// String returns the UDT values in human readable format.
func (u *UDT) String() string {
	return fmt.Sprintf("{Type: %d, CalledPartyAddress: %v, CallingPartyAddress: %v, DataLength: %d, Data: %x}",
		u.Type,
		u.CalledPartyAddress,
		u.CallingPartyAddress,
		u.DataLength,
		u.Data,
	)
}

// MessageType returns the Message Type in int.
func (u *UDT) MessageType() uint8 {
	return MsgTypeUDT
}

// MessageTypeName returns the Message Type in string.
func (u *UDT) MessageTypeName() string {
	return "UDT"
}

// CdGT returns the GT in CalledPartyAddress in human readable string.
func (u *UDT) CdGT() string {
	return u.CalledPartyAddress.GTString()
}

// CgGT returns the GT in CalledPartyAddress in human readable string.
func (u *UDT) CgGT() string {
	return u.CallingPartyAddress.GTString()
}
