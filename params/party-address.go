// Copyright 2019 go-sccp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

import (
	"encoding/binary"

	"github.com/pkg/errors"
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

// GlobalTitle is a GlobalTitle inside the Called/Calling Party Address.
type GlobalTitle struct {
	TranslationType          uint8
	NumberingPlan            int // 1/2 Octet
	EncodingScheme           int // 1/2 Octet
	NatureOfAddressIndicator uint8
	GlobalTitleInfo          []byte
}

// NewPartyAddress creates a new PartyAddress including GlobalTitle.
func NewPartyAddress(gti, spc, ssn, tt, np, es, nai int, gt []byte) *PartyAddress {
	p := &PartyAddress{
		Indicator:          uint8(gti),
		SignalingPointCode: uint16(spc),
		SubsystemNumber:    uint8(ssn),
		GlobalTitle: &GlobalTitle{
			TranslationType:          uint8(tt),
			NumberingPlan:            np,
			EncodingScheme:           es,
			NatureOfAddressIndicator: uint8(nai),
			GlobalTitleInfo:          gt,
		},
	}
	p.Length = uint8(p.MarshalLen() - 1)
	return p
}

// MarshalBinary returns the byte sequence generated from a PartyAddress instance.
func (p *PartyAddress) MarshalBinary() ([]byte, error) {
	b := make([]byte, p.MarshalLen())
	if err := p.MarshalTo(b); err != nil {
		return nil, errors.Wrap(err, "failed to serialize PartyAddress:")
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

	gt := p.GlobalTitle
	switch p.GTI() {
	case 1:
		b[offset] = gt.NatureOfAddressIndicator
		offset++
	case 2:
		b[offset] = gt.TranslationType
		offset++
	case 3:
		b[offset] = gt.TranslationType
		b[offset+1] = uint8(p.NumberingPlan<<4 | p.EncodingScheme)
		offset += 2
	case 4:
		b[offset] = gt.TranslationType
		b[offset+1] = uint8(p.NumberingPlan<<4 | p.EncodingScheme)
		b[offset+2] = p.NatureOfAddressIndicator
		offset += 3
	}

	copy(b[offset:p.MarshalLen()], gt.GlobalTitleInfo)
	return nil
}

// ParsePartyAddress decodes given byte sequence as a SCCP common header.
func ParsePartyAddress(b []byte) (*PartyAddress, error) {
	p := &PartyAddress{
		GlobalTitle: &GlobalTitle{},
	}
	if err := p.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return p, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a SCCP common header.
func (p *PartyAddress) UnmarshalBinary(b []byte) error {
	l := len(b)
	if l < 2 {
		return ErrTooShortToDecode
	}

	p.Length = b[0]
	p.Indicator = b[1]

	var offset = 2
	if p.HasPC() {
		p.SignalingPointCode = binary.BigEndian.Uint16(b[offset : offset+2])
		offset += 2
	}
	if p.HasSSN() {
		p.SubsystemNumber = b[offset]
		offset++
	}

	gt := p.GlobalTitle
	switch p.GTI() {
	case 1:
		gt.NatureOfAddressIndicator = b[offset]
		offset++
	case 2:
		gt.TranslationType = b[offset]
		offset++
	case 3:
		gt.TranslationType = b[offset]
		gt.NumberingPlan = int(b[offset+1]) >> 4 & 0xf
		gt.EncodingScheme = int(b[offset+1]) & 0xf
		offset += 2
	case 4:
		gt.TranslationType = b[3]
		gt.NumberingPlan = int(b[offset+1]) >> 4 & 0xf
		gt.EncodingScheme = int(b[offset+1]) & 0xf
		gt.NatureOfAddressIndicator = b[offset+2]
		offset += 3
	}

	gt.GlobalTitleInfo = b[offset:l]

	return nil
}

// MarshalLen returns the serial length.
func (p *PartyAddress) MarshalLen() int {
	l := 2 + len(p.GlobalTitle.GlobalTitleInfo)
	if p.HasPC() {
		l += 2
	}
	if p.HasSSN() {
		l++
	}
	switch p.GTI() {
	case 1:
		l++
	case 2:
		l++
	case 3:
		l += 2
	case 4:
		l += 3
	}

	return l
}

// SetLength sets the length in Length field.
func (p *PartyAddress) SetLength() {
	l := 1 + len(p.GlobalTitle.GlobalTitleInfo)
	if p.HasPC() {
		l += 2
	}
	if p.HasSSN() {
		l++
	}
	switch p.GTI() {
	case 1:
		l++
	case 2:
		l++
	case 3:
		l += 2
	case 4:
		l += 3
	}

	p.Length = uint8(l)
}

// RouteOnGT reports whether the packet is routed on Global Title or not.
func (p *PartyAddress) RouteOnGT() bool {
	return (int(p.Indicator) >> 6 & 0x1) == 0
}

// GTI returns GlobalTitleIndicator value retrieved from Indicator.
func (p *PartyAddress) GTI() int {
	return (int(p.Indicator) >> 2 & 0xf)
}

// HasSSN reports whether PartyAddress has a Subsystem Number.
func (p *PartyAddress) HasSSN() bool {
	return (int(p.Indicator) >> 1 & 0x1) == 1
}

// HasPC reports whether PartyAddress has a Signaling Point Code.
func (p *PartyAddress) HasPC() bool {
	return (int(p.Indicator) & 0x1) == 1
}

// IsOddDigits reports whether GlobalTitleInfo is odd number or not.
func (p *PartyAddress) IsOddDigits() bool {
	return p.EncodingScheme == 1
}

// GTString returns the GlobalTitleInfo in human readable string.
func (p *PartyAddress) GTString() string {
	return utils.SwappedBytesToStr(p.GlobalTitleInfo, p.IsOddDigits())
}
