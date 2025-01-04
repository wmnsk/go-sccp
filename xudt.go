// Copyright 2019-2024 go-sccp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package sccp

import (
	"fmt"
	"io"

	"github.com/wmnsk/go-sccp/params"
)

// XUDT represents a SCCP Message Extended unitdata (XUDT).
type XUDT struct {
	Type                    MsgType
	ProtocolClass           *params.ProtocolClass
	HopCounter              *params.HopCounter
	CalledPartyAddress      *params.PartyAddress
	CallingPartyAddress     *params.PartyAddress
	Data                    *params.Data
	Segmentation            *params.Segmentation
	Importance              *params.Importance
	EndOfOptionalParameters *params.EndOfOptionalParameters

	ptr1, ptr2, ptr3, ptr4 uint8
}

// NewXUDT creates a new XUDT.
func NewXUDT(pcls int, retOnErr bool, hc uint8, cdpa, cgpa *params.PartyAddress, data []byte, opts ...params.Parameter) *XUDT {
	x := &XUDT{
		Type:                MsgTypeXUDT,
		ProtocolClass:       params.NewProtocolClass(pcls, retOnErr),
		HopCounter:          params.NewHopCounter(hc),
		CalledPartyAddress:  cdpa,
		CallingPartyAddress: cgpa,
		Data:                params.NewData(data),
	}

	x.ptr1 = 4
	x.ptr2 = x.ptr1 + uint8(cdpa.MarshalLen()) - 1
	x.ptr3 = x.ptr2 + uint8(cgpa.MarshalLen()) - 1
	x.ptr4 = 0

	for _, opt := range opts {
		switch opt.Code() {
		case params.PCodeSegmentation:
			x.Segmentation = opt.(*params.Segmentation)
		case params.PCodeImportance:
			x.Importance = opt.(*params.Importance)
		case params.PCodeEndOfOptionalParameters:
			x.EndOfOptionalParameters = opt.(*params.EndOfOptionalParameters)
		default:
			logf("unexpected parameter: %s in NewXUDT", opt.Code())
		}
	}

	if len(opts) > 0 {
		x.ptr4 = x.ptr3 + uint8(x.Data.MarshalLen()) - 1
		// so that users don't have to give EndOfOptionalParameters explicitly
		x.EndOfOptionalParameters = params.NewEndOfOptionalParameters()
	}

	return x
}

// MarshalBinary returns the byte sequence generated from a XUDT instance.
func (x *XUDT) MarshalBinary() ([]byte, error) {
	b := make([]byte, x.MarshalLen())
	if err := x.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
// SCCP is dependent on the Pointers when serializing, which means that it might fail when invalid Pointers are set.
func (x *XUDT) MarshalTo(b []byte) error {
	l := len(b)
	if l < 5 {
		return io.ErrUnexpectedEOF
	}

	b[0] = uint8(x.Type)

	n := 1
	m, err := x.ProtocolClass.Write(b[1:])
	if err != nil {
		return err
	}
	n += m

	m, err = x.HopCounter.Write(b[n:])
	if err != nil {
		return err
	}
	n += m

	b[n] = x.ptr1
	if p := int(x.ptr1); l < p {
		return io.ErrUnexpectedEOF
	}
	b[n+1] = x.ptr2
	if p := int(x.ptr2 + 4); l < p {
		return io.ErrUnexpectedEOF
	}
	b[n+2] = x.ptr3
	if p := int(x.ptr3 + 5); l < p {
		return io.ErrUnexpectedEOF
	}
	b[n+3] = x.ptr4
	if p := int(x.ptr4 + 6); l < p {
		return io.ErrUnexpectedEOF
	}
	n += 4

	cdpaEnd := int(x.ptr2 + 4)
	cgpaEnd := int(x.ptr3 + 5)
	dataEnd := int(x.ptr4 + 6)
	if _, err := x.CalledPartyAddress.Write(b[n:cdpaEnd]); err != nil {
		return err
	}

	if _, err := x.CallingPartyAddress.Write(b[cdpaEnd:cgpaEnd]); err != nil {
		return err
	}

	if _, err := x.Data.Write(b[cgpaEnd:]); err != nil {
		return err
	}

	if x.ptr4 == 0 {
		return nil
	}

	offset := dataEnd
	if param := x.Segmentation; param != nil {
		m, err := param.Write(b[offset:])
		if err != nil {
			return err
		}
		offset += m
	}
	if param := x.Importance; param != nil {
		m, err := param.Write(b[offset:])
		if err != nil {
			return err
		}
		offset += m
	}
	if param := x.EndOfOptionalParameters; param != nil {
		_, err := param.Write(b[offset:])
		if err != nil {
			return err
		}
	}

	return nil
}

// ParseXUDT decodes given byte sequence as a SCCP XUDT.
func ParseXUDT(b []byte) (*XUDT, error) {
	x := &XUDT{}
	if err := x.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return x, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a SCCP XUDT.
func (x *XUDT) UnmarshalBinary(b []byte) error {
	l := len(b)
	if l <= 5 { // where CdPA starts
		return io.ErrUnexpectedEOF
	}

	x.Type = MsgType(b[0])

	offset := 1
	x.ProtocolClass = &params.ProtocolClass{}
	n, err := x.ProtocolClass.Read(b[offset:])
	if err != nil {
		return err
	}
	offset += n

	x.HopCounter = &params.HopCounter{}
	n, err = x.HopCounter.Read(b[offset:])
	if err != nil {
		return err
	}
	offset += n

	x.ptr1 = b[offset]
	if l < int(x.ptr1) {
		return io.ErrUnexpectedEOF
	}
	x.ptr2 = b[offset+1]
	if l < int(x.ptr2+4) { // where CgPA starts
		return io.ErrUnexpectedEOF
	}
	x.ptr3 = b[offset+2]
	if l < int(x.ptr3+5) { // where Data starts
		return io.ErrUnexpectedEOF
	}
	x.ptr4 = b[offset+3]
	if m := int(x.ptr4); (m != 0) && l < m+7 { // where optional parameters start
		return io.ErrUnexpectedEOF
	}

	offset += 4
	cdpaEnd := int(x.ptr2 + 4)
	cgpaEnd := int(x.ptr3 + 5)
	dataEnd := int(x.ptr4 + 6)
	x.CalledPartyAddress, _, err = params.ParseCalledPartyAddress(b[offset:cdpaEnd])
	if err != nil {
		return err
	}

	x.CallingPartyAddress, _, err = params.ParseCallingPartyAddress(b[cdpaEnd:cgpaEnd])
	if err != nil {
		return err
	}

	x.Data, _, err = params.ParseData(b[cgpaEnd:])
	if err != nil {
		return err
	}

	if x.ptr4 == 0 {
		return nil
	}

	opts, _, err := params.ParseOptionalParameters(b[dataEnd:])
	if err != nil {
		return err
	}

	for _, opt := range opts {
		switch opt.Code() {
		case params.PCodeSegmentation:
			x.Segmentation = opt.(*params.Segmentation)
		case params.PCodeImportance:
			x.Importance = opt.(*params.Importance)
		case params.PCodeEndOfOptionalParameters:
			x.EndOfOptionalParameters = opt.(*params.EndOfOptionalParameters)
		}
	}

	return nil
}

// MarshalLen returns the serial length.
func (x *XUDT) MarshalLen() int {
	l := 7 // MsgType + ProtocolClass + HopCounter + Pointers

	// if optional parameters exist
	if x.ptr4 != 0 {
		l += int(x.ptr4) - 1 // length without optional parameters
		if param := x.Segmentation; param != nil {
			l += param.MarshalLen()
		}
		if param := x.Importance; param != nil {
			l += param.MarshalLen()
		}
		if param := x.EndOfOptionalParameters; param != nil {
			l += param.MarshalLen()
		}

		return l
	}

	l += int(x.ptr3) - 2 // length without Data
	if param := x.Data; param != nil {
		l += param.MarshalLen()
	}

	return l
}

// String returns the XUDT values in human readable format.
func (x *XUDT) String() string {
	return fmt.Sprintf("%s: {ProtocolClass: %s, HopCounter: %s, CalledPartyAddress: %v, CallingPartyAddress: %v, Data: %s, Segmentation: %s, Importance: %s}",
		x.Type,
		x.ProtocolClass,
		x.HopCounter,
		x.CalledPartyAddress,
		x.CallingPartyAddress,
		x.Data,
		x.Segmentation,
		x.Importance,
	)
}

// MessageType returns the Message Type in int.
func (x *XUDT) MessageType() MsgType {
	return MsgTypeXUDT
}

// MessageTypeName returns the Message Type in string.
func (x *XUDT) MessageTypeName() string {
	return x.MessageType().String()
}

// CdGT returns the GT in CalledPartyAddress in human readable string.
func (x *XUDT) CdGT() string {
	if x.CalledPartyAddress.GlobalTitle == nil {
		return ""
	}
	return x.CalledPartyAddress.Address()
}

// CgGT returns the GT in CalledPartyAddress in human readable string.
func (x *XUDT) CgGT() string {
	if x.CallingPartyAddress.GlobalTitle == nil {
		return ""
	}
	return x.CallingPartyAddress.Address()
}
