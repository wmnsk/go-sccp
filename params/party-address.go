// Copyright 2019-2024 go-sccp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/wmnsk/go-sccp/utils"
)

// PartyAddress is a SCCP parameter that represents a Called/Calling Party Address.
type PartyAddress struct {
	Length             uint8
	Indicator          uint8
	SignalingPointCode uint16
	SubsystemNumber    uint8
	*GlobalTitle
}

// NewAddressIndicator creates a new AddressIndicator, which is meant to be used in
// NewPartyAddress as the first argument.
//
// The last bit, which is "reserved for national use", is always set to 0.
// You can set the bit to 1 by doing `| 0b10000000` to the result of this function.
func NewAddressIndicator(hasPC, hasSSN, routeOnSSN bool, gti GlobalTitleIndicator) uint8 {
	var ai uint8
	if hasPC {
		ai |= 0b00000001
	}
	if hasSSN {
		ai |= 0b00000010
	}
	if routeOnSSN {
		ai |= 0b01000000
	}
	ai |= uint8(gti) << 2

	return ai
}

// NewPartyAddress creates a new PartyAddress including GlobalTitle.
// Deprecated: Use NewPartyAddressTyped instead.
func NewPartyAddress(ai, spc, ssn, tt, np, es, nai int, addr []byte) *PartyAddress {
	var globalTitle *GlobalTitle
	gti := gti(ai)
	if gti == 0 {
		globalTitle = nil
	} else {
		globalTitle = NewGlobalTitle(
			gti,
			TranslationType(tt),
			NumberingPlan(np),
			EncodingScheme(es),
			NatureOfAddressIndicator(nai),
			addr,
		)
	}
	return NewPartyAddressTyped(uint8(ai), uint16(spc), uint8(ssn), globalTitle)
}

// NewPartyAddress creates a new PartyAddress from properly-typed values.
//
// The given SPC and SSN are set to 0 if the corresponding bit is not properly set in the
// AddressIndicator. Use NewAddressIndicator to create a proper AddressIndicator.
func NewPartyAddressTyped(ai uint8, spc uint16, ssn uint8, gt *GlobalTitle) *PartyAddress {
	p := &PartyAddress{
		Indicator:   ai,
		GlobalTitle: gt,
	}

	if p.HasPC() {
		p.SignalingPointCode = spc
	}

	if p.HasSSN() {
		p.SubsystemNumber = ssn
	}

	p.SetLength()
	return p
}

// MarshalBinary returns the byte sequence generated from a PartyAddress instance.
func (p *PartyAddress) MarshalBinary() ([]byte, error) {
	b := make([]byte, p.MarshalLen())
	if err := p.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (p *PartyAddress) MarshalTo(b []byte) error {
	b[0] = p.Length
	b[1] = p.Indicator
	var offset = 2
	if p.HasPC() {
		binary.BigEndian.PutUint16(b[offset:offset+2], p.SignalingPointCode)
		offset += 2
	}
	if p.HasSSN() {
		b[offset] = p.SubsystemNumber
		offset++
	}

	if p.GlobalTitle != nil {
		return p.GlobalTitle.MarshalTo(b[offset:p.MarshalLen()])
	}

	return nil
}

// ParsePartyAddress decodes given byte sequence as a SCCP common header.
func ParsePartyAddress(b []byte) (*PartyAddress, error) {
	p := new(PartyAddress)
	if err := p.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return p, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a SCCP common header.
func (p *PartyAddress) UnmarshalBinary(b []byte) error {
	if len(b) < 2 {
		return io.ErrUnexpectedEOF
	}
	p.Length = b[0]
	l := int(p.Length)
	if l >= len(b) {
		return io.ErrUnexpectedEOF
	}
	p.Indicator = b[1]

	var offset = 2
	if p.HasPC() {
		end := offset + 2
		if end >= len(b) {
			return io.ErrUnexpectedEOF
		}
		p.SignalingPointCode = binary.BigEndian.Uint16(b[offset:end])
		offset = end
	}
	if p.HasSSN() {
		p.SubsystemNumber = b[offset]
		offset++
	}

	gti := p.GTI()
	if gti == 0 {
		return nil
	}

	gt := &GlobalTitle{GTI: gti}
	if err := gt.UnmarshalBinary(b[offset : l+1]); err != nil {
		return err
	}
	p.GlobalTitle = gt

	return nil
}

// MarshalLen returns the serial length.
func (p *PartyAddress) MarshalLen() int {
	l := 2
	if p.HasPC() {
		l += 2
	}
	if p.HasSSN() {
		l++
	}

	if p.GlobalTitle != nil {
		l = l + p.GlobalTitle.MarshalLen()
	}

	return l
}

// SetLength sets the length in Length field.
func (p *PartyAddress) SetLength() {
	p.Length = uint8(p.MarshalLen()) - 1
}

// RouteOnGT reports whether the packet is routed on Global Title or not.
func (p *PartyAddress) RouteOnGT() bool {
	return (int(p.Indicator) >> 6 & 0b1) == 0
}

// RouteOnSSN reports whether the packet is routed on SSN or not.
func (p *PartyAddress) RouteOnSSN() bool {
	return !p.RouteOnGT()
}

// GTI returns GlobalTitleIndicator value retrieved from Indicator.
func (p *PartyAddress) GTI() GlobalTitleIndicator {
	return gti(int(p.Indicator))
}

func gti(ai int) GlobalTitleIndicator {
	return GlobalTitleIndicator(ai >> 2 & 0b1111)
}

// HasSSN reports whether PartyAddress has a Subsystem Number.
func (p *PartyAddress) HasSSN() bool {
	return (int(p.Indicator) >> 1 & 0b1) == 1
}

// HasPC reports whether PartyAddress has a Signaling Point Code.
func (p *PartyAddress) HasPC() bool {
	return (int(p.Indicator) & 0b1) == 1
}

// IsOddDigits reports whether AddressInformation is odd number or not.
func (p *PartyAddress) IsOddDigits() bool {
	return p.EncodingScheme == 1
}

// GTString returns the AddressInformation in human readable string.
func (p *PartyAddress) GTString() string {
	return utils.SwappedBytesToStr(p.AddressInformation, p.IsOddDigits())
}

// String returns the PartyAddress values in human readable format.
func (p *PartyAddress) String() string {
	return fmt.Sprintf("{Length: %d, Indicator: %#08b, SignalingPointCode: %d, SubsystemNumber: %d, GlobalTitle: %v}",
		p.Length, p.Indicator, p.SignalingPointCode, p.SubsystemNumber, p.GlobalTitle,
	)
}
