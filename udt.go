// Copyright 2019-2024 go-sccp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package sccp

import (
	"fmt"
	"io"

	"github.com/wmnsk/go-sccp/params"
)

// UDT represents a SCCP Message Unit Data (UDT).
type UDT struct {
	Type                MsgType
	ProtocolClass       *params.ProtocolClass
	CalledPartyAddress  *params.PartyAddress
	CallingPartyAddress *params.PartyAddress
	Data                *params.Data

	ptr1, ptr2, ptr3 uint8
}

// NewUDT creates a new UDT.
func NewUDT(pcls int, retOnErr bool, cdpa, cgpa *params.PartyAddress, data []byte) *UDT {
	u := &UDT{
		Type:                MsgTypeUDT,
		ProtocolClass:       params.NewProtocolClass(pcls, retOnErr),
		CalledPartyAddress:  cdpa.AsCalled(),
		CallingPartyAddress: cgpa.AsCalling(),
		Data:                params.NewData(data),
	}

	u.ptr1 = 3
	u.ptr2 = u.ptr1 + uint8(cdpa.MarshalLen()) - 1
	u.ptr3 = u.ptr2 + uint8(cgpa.MarshalLen()) - 1

	return u
}

// MarshalBinary returns the byte sequence generated from a UDT instance.
func (u *UDT) MarshalBinary() ([]byte, error) {
	b := make([]byte, u.MarshalLen())
	if err := u.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
// SCCP is dependent on the Pointers when serializing, which means that it might fail when invalid Pointers are set.
func (u *UDT) MarshalTo(b []byte) error {
	l := len(b)
	if l < 5 {
		return io.ErrUnexpectedEOF
	}

	b[0] = uint8(u.Type)

	n := 1
	m, err := u.ProtocolClass.Write(b[1:])
	if err != nil {
		return err
	}
	n += m

	b[n] = u.ptr1
	if p := int(u.ptr1); l < p {
		return io.ErrUnexpectedEOF
	}
	b[n+1] = u.ptr2
	if p := int(u.ptr2 + 3); l < p {
		return io.ErrUnexpectedEOF
	}
	b[n+2] = u.ptr3
	if p := int(u.ptr3 + 5); l < p {
		return io.ErrUnexpectedEOF
	}
	n += 3

	cdpaEnd := int(u.ptr2 + 3)
	cgpaEnd := int(u.ptr3 + 4)
	if _, err := u.CalledPartyAddress.Write(b[n:cdpaEnd]); err != nil {
		return err
	}

	if _, err := u.CallingPartyAddress.Write(b[cdpaEnd:cgpaEnd]); err != nil {
		return err
	}

	if _, err := u.Data.Write(b[cgpaEnd:]); err != nil {
		return err
	}

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
	if l <= 5 { // where CdPA starts
		return io.ErrUnexpectedEOF
	}

	u.Type = MsgType(b[0])

	offset := 1
	u.ProtocolClass = &params.ProtocolClass{}
	n, err := u.ProtocolClass.Read(b[offset:])
	if err != nil {
		return err
	}
	offset += n

	u.ptr1 = b[offset]
	if l < int(u.ptr1) {
		return io.ErrUnexpectedEOF
	}
	u.ptr2 = b[offset+1]
	if l < int(u.ptr2+3) { // where CgPA starts
		return io.ErrUnexpectedEOF
	}
	u.ptr3 = b[offset+2]
	if l < int(u.ptr3+5) { // where u.Data starts
		return io.ErrUnexpectedEOF
	}

	offset += 3
	cdpaEnd := int(u.ptr2 + 3)
	cgpaEnd := int(u.ptr3 + 4)
	u.CalledPartyAddress, err = params.ParseCalledPartyAddress(b[offset:cdpaEnd])
	if err != nil {
		return err
	}

	u.CallingPartyAddress, err = params.ParseCallingPartyAddress(b[cdpaEnd:cgpaEnd])
	if err != nil {
		return err
	}

	u.Data = &params.Data{}
	if _, err := u.Data.Read(b[cgpaEnd:]); err != nil {
		return err
	}

	return nil
}

// MarshalLen returns the serial length.
func (u *UDT) MarshalLen() int {
	l := 5 // MsgType, ProtocolClass, pointers

	l += int(u.ptr3) - 1 // length without Data
	if param := u.Data; param != nil {
		l += param.MarshalLen()
	}

	return l
}

// String returns the UDT values in human readable format.
func (u *UDT) String() string {
	return fmt.Sprintf("%s: {ProtocolClass: %s, CalledPartyAddress: %v, CallingPartyAddress: %v, Data: %s}",
		u.Type,
		u.ProtocolClass,
		u.CalledPartyAddress,
		u.CallingPartyAddress,
		u.Data,
	)
}

// MessageType returns the Message Type in int.
func (u *UDT) MessageType() MsgType {
	return MsgTypeUDT
}

// MessageTypeName returns the Message Type in string.
func (u *UDT) MessageTypeName() string {
	return u.MessageType().String()
}

// CdGT returns the GT in CalledPartyAddress in human readable string.
func (u *UDT) CdGT() string {
	if u.CalledPartyAddress.GlobalTitle == nil {
		return ""
	}
	return u.CalledPartyAddress.Address()
}

// CgGT returns the GT in CalledPartyAddress in human readable string.
func (u *UDT) CgGT() string {
	if u.CallingPartyAddress.GlobalTitle == nil {
		return ""
	}
	return u.CallingPartyAddress.Address()
}
