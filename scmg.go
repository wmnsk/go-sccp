// Copyright 2019-2024 go-sccp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package sccp

import (
	"encoding/binary"
	"fmt"
	"io"
)

// SCMGType is type of SCMG message.
type SCMGType uint8

// Table 23/Q.713
const (
	_           SCMGType = iota
	SCMGTypeSSA          // SSA
	SCMGTypeSSP          // SSP
	SCMGTypeSST          // SST
	SCMGTypeSOR          // SOR
	SCMGTypeSOG          // SOG
	SCMGTypeSSC          // SSC
)

// SCMG represents a SCCP Management message (SCMG).
// Chapter 5.3/Q.713
type SCMG struct {
	Type                           SCMGType
	AffectedSSN                    uint8
	AffectedPC                     uint16
	SubsystemMultiplicityIndicator uint8
	SCCPCongestionLevel            uint8
}

// NewSCMG creates a new SCMG.
func NewSCMG(typ SCMGType, assn uint8, apc uint16, smi uint8, scl uint8) *SCMG {
	return &SCMG{
		Type:                           typ,
		AffectedSSN:                    assn,
		AffectedPC:                     apc,
		SubsystemMultiplicityIndicator: smi,
		SCCPCongestionLevel:            scl,
	}
}

// MarshalBinary returns the byte sequence generated from a SCMG instance.
func (s *SCMG) MarshalBinary() ([]byte, error) {
	b := make([]byte, s.MarshalLen())
	if err := s.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (s *SCMG) MarshalTo(b []byte) error {
	l := len(b)

	if l < s.MarshalLen() {
		return io.ErrUnexpectedEOF
	}

	b[0] = uint8(s.Type)
	b[1] = s.AffectedSSN
	binary.LittleEndian.PutUint16(b[2:4], s.AffectedPC)
	b[4] = s.SubsystemMultiplicityIndicator
	if s.Type == SCMGTypeSSC {
		b[5] = s.SCCPCongestionLevel
	}

	return nil
}

// ParseSCMG decodes given byte sequence as a SCMG.
func ParseSCMG(b []byte) (*SCMG, error) {
	s := &SCMG{}
	if err := s.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return s, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a SCMG.
func (s *SCMG) UnmarshalBinary(b []byte) error {
	l := len(b)
	if l < 5 {
		return io.ErrUnexpectedEOF
	}

	s.Type = SCMGType(b[0])
	s.AffectedSSN = b[1]
	s.AffectedPC = binary.LittleEndian.Uint16(b[2:4])
	s.SubsystemMultiplicityIndicator = b[4]

	if s.Type == SCMGTypeSSC {
		if l < 6 {
			return io.ErrUnexpectedEOF
		}
		s.SCCPCongestionLevel = b[5]
	}

	return nil
}

// MarshalLen returns the serial length.
func (s *SCMG) MarshalLen() int {
	// Table 24/Q.713 – SCMG messages
	l := 5

	// Table 25/Q.713 – SSC
	if s.Type == SCMGTypeSSC {
		l += 1
	}

	return l
}

// String returns the SCMG values in human readable format.
func (s *SCMG) String() string {
	return fmt.Sprintf("{Type: %d, AffectedSSN: %v, AffectedPC: %v, SubsystemMultiplicityIndicator: %d, SCCPCongestionLevel: %d}",
		s.Type,
		s.AffectedSSN,
		s.AffectedPC,
		s.SubsystemMultiplicityIndicator,
		s.SCCPCongestionLevel,
	)
}

// MessageType returns the Message Type in int.
func (s *SCMG) MessageType() SCMGType {
	return s.Type
}

// MessageTypeName returns the Message Type in string.
func (s *SCMG) MessageTypeName() string {
	return s.Type.String()
}
