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
		CalledPartyAddress:  cdpa,
		CallingPartyAddress: cgpa,
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
	if l <= 5 {
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
	offsetPtr1 := 2 + int(u.ptr1)
	if l < offsetPtr1+1 { // where CdPA starts
		return io.ErrUnexpectedEOF
	}
	u.ptr2 = b[offset+1]
	offsetPtr2 := 3 + int(u.ptr2)
	if l < offsetPtr2+1 { // where CgPA starts
		return io.ErrUnexpectedEOF
	}
	u.ptr3 = b[offset+2]
	offsetPtr3 := 4 + int(u.ptr3)
	if l < offsetPtr3+1 { // where u.Data starts
		return io.ErrUnexpectedEOF
	}

	cdpaEnd := offsetPtr1 + int(b[offsetPtr1]) + 1 // +1 is the data length included from the beginning
	if l < cdpaEnd {                               // where CdPA ends
		return io.ErrUnexpectedEOF
	}
	cgpaEnd := offsetPtr2 + int(b[offsetPtr2]) + 1
	if l < cgpaEnd { // where CgPA ends
		return io.ErrUnexpectedEOF
	}
	dataEnd := offsetPtr3 + int(b[offsetPtr3]) + 1
	if l < dataEnd { // where Data ends
		return io.ErrUnexpectedEOF
	}

	u.CalledPartyAddress, _, err = params.ParseCalledPartyAddress(b[offsetPtr1:cdpaEnd])
	if err != nil {
		return err
	}

	u.CallingPartyAddress, _, err = params.ParseCallingPartyAddress(b[offsetPtr2:cgpaEnd])
	if err != nil {
		return err
	}

	u.Data = &params.Data{}
	if _, err := u.Data.Read(b[offsetPtr3:dataEnd]); err != nil {
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
