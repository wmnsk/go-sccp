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

// UnsupportedParameterError indicates the value in Version field is invalid.
type UnsupportedParameterError uint8

// Error returns the type of receiver and some additional message.
func (e UnsupportedParameterError) Error() string {
	return fmt.Sprintf("sccp: got unsupported type %d", e)
}

// Parameter is an interface that all SCCP parameters have to implement.
type Parameter interface {
	io.ReadWriter
	Code() ParameterNameCode
	String() string
}

// ParameterType is a type for Parameter described in the tables in section 4 of Q.713.
type ParameterType uint8

const (
	// F: mandatory fixed length parameter
	PTypeF ParameterType = 0 // F
	// V: mandatory variable length parameter
	PTypeV ParameterType = 1 // V
	// O: optional parameter of fixed or variable length
	PTypeO ParameterType = 2 // O
)

// ParameterNameCode is a type of Parameter Name Code defined in Q.713 Table 2.
type ParameterNameCode uint8

// ParameterNameCode values.
const (
	// O
	PCodeEndOfOptionalParameters ParameterNameCode = 0b00000000 // End of optional parameters
	// F
	PCodeDestinationLocalReference ParameterNameCode = 0b00000001 // Destination local reference
	// F
	PCodeSourceLocalReference ParameterNameCode = 0b00000010 // Source local reference
	// V, O
	PCodeCalledPartyAddress ParameterNameCode = 0b00000011 // Called party address
	// V, O
	PCodeCallingPartyAddress ParameterNameCode = 0b00000100 // Calling party address
	// F
	PCodeProtocolClass ParameterNameCode = 0b00000101 // Protocol class
	// F
	PCodeSegmentingReassembling ParameterNameCode = 0b00000110 // Segmenting/reassembling
	// F
	PCodeReceiveSequenceNumber ParameterNameCode = 0b00000111 // Receive sequence number
	// F
	PCodeSequencingSegmenting ParameterNameCode = 0b00001000 // Sequencing/segmenting
	// F, O
	PCodeCredit ParameterNameCode = 0b00001001 // Credit
	// F
	PCodeReleaseCause ParameterNameCode = 0b00001010 // Release cause
	// F
	PCodeReturnCause ParameterNameCode = 0b00001011 // Return cause
	// F
	PCodeResetCause ParameterNameCode = 0b00001100 // Reset cause
	// F
	PCodeErrorCause ParameterNameCode = 0b00001101 // Error cause
	// F
	PCodeRefusalCause ParameterNameCode = 0b00001110 // Refusal cause
	// V, O
	PCodeData ParameterNameCode = 0b00001111 // Data
	// O
	PCodeSegmentation ParameterNameCode = 0b00010000 // Segmentation
	// F, O
	PCodeHopCounter ParameterNameCode = 0b00010001 // Hop Counter
	// O
	PCodeImportance ParameterNameCode = 0b00010010 // Importance
	// V
	PCodeLongData ParameterNameCode = 0b00010011 // Long data
)

// ParseOptionalParameters parses optional parameters from the given byte sequence.
func ParseOptionalParameters(b []byte) ([]Parameter, error) {
	var params []Parameter
	var offset int
	for len(b) > 0 {
		p, n, err := ParseOptionalParameter(b[offset:])
		if err != nil {
			return nil, err
		}
		params = append(params, p)
		if p.Code() == PCodeEndOfOptionalParameters {
			break
		}
		offset += n
	}
	return params, nil
}

// ParseOptionalParameter parses a single optional parameter from the given byte sequence.
func ParseOptionalParameter(b []byte) (Parameter, int, error) {
	if len(b) < 1 {
		return nil, 0, io.ErrUnexpectedEOF
	}

	var p Parameter
	switch ParameterNameCode(b[0]) {
	case PCodeEndOfOptionalParameters:
		p = &EndOfOptionalParameters{paramType: PTypeO}
	case PCodeCalledPartyAddress:
		p = &PartyAddress{
			paramType: PTypeO,
			code:      PCodeCalledPartyAddress,
		}
	case PCodeCallingPartyAddress:
		p = &PartyAddress{
			paramType: PTypeO,
			code:      PCodeCallingPartyAddress,
		}
	case PCodeCredit:
		p = &Credit{paramType: PTypeO}
	case PCodeData:
		p = &Data{paramType: PTypeO}
	case PCodeSegmentation:
		p = &Segmentation{paramType: PTypeO}
	case PCodeHopCounter:
		p = &HopCounter{paramType: PTypeO}
	case PCodeImportance:
		p = &Importance{paramType: PTypeO}
	default:
		return nil, 0, UnsupportedParameterError(b[0])
	}

	n, err := p.Read(b)
	if err != nil {
		return nil, n, err
	}
	return p, n, nil
}

// EndOfOptionalParameters represents the End Of Optional Parameters.
type EndOfOptionalParameters struct {
	paramType ParameterType
	code      ParameterNameCode
	length    int
	value     uint8
}

// NewEndOfOptionalParameters creates a new EndOfOptionalParameters.
func NewEndOfOptionalParameters() *EndOfOptionalParameters {
	return &EndOfOptionalParameters{
		paramType: PTypeO,
		code:      PCodeEndOfOptionalParameters,
		length:    1,
		value:     0,
	}
}

// Read sets the values retrieved from byte sequence in a EndOfOptionalParameters.
func (e *EndOfOptionalParameters) Read(b []byte) (int, error) {
	n := 1
	if len(b) < n {
		return 0, io.ErrUnexpectedEOF
	}

	e.code = PCodeEndOfOptionalParameters
	e.length = n
	e.value = b[0]

	if e.value != 0 {
		logf("invalid parameter %s=%d: must be 0", e.code, e.value)
	}

	return n, nil
}

// Write serializes the EndOfOptionalParameters parameter and returns it as a byte slice.
func (e *EndOfOptionalParameters) Write(b []byte) (int, error) {
	if len(b) < e.length {
		return 0, io.ErrUnexpectedEOF
	}

	b[0] = e.value
	return e.length, nil
}

// Value returns the EndOfOptionalParameters in uint8.
func (e *EndOfOptionalParameters) Value() uint8 {
	return e.value
}

// Code returns the EndOfOptionalParameters in ParameterNameCode.
func (e *EndOfOptionalParameters) Code() ParameterNameCode {
	return e.code
}

// String returns the EndOfOptionalParameters in string.
func (e *EndOfOptionalParameters) String() string {
	return fmt.Sprintf("%s: %d", e.code, e.value)
}

// LocalReference represents the Destination/Source Local Reference.
type LocalReference struct {
	paramType ParameterType
	code      ParameterNameCode
	length    int
	value     []byte
}

// NewLocalReference creates a new LocalReference.
//
// LocalReference parameter is three octets long. The fourth
// octet is masked out.
func NewLocalReference(c ParameterNameCode, v uint32) *LocalReference {
	return &LocalReference{
		paramType: PTypeF,
		code:      c,
		length:    3,
		value:     utils.Uint32To24(v),
	}
}

// NewDestinationLocalReference creates a new LocalReference for Destination.
func NewDestinationLocalReference(v uint32) *LocalReference {
	return NewLocalReference(PCodeDestinationLocalReference, v)
}

// NewSourceLocalReference creates a new LocalReference for Source.
func NewSourceLocalReference(v uint32) *LocalReference {
	return NewLocalReference(PCodeSourceLocalReference, v)
}

// ParseDestinationLocalReference parses the given byte sequence as a Destination Local Reference.
func ParseDestinationLocalReference(b []byte) (*LocalReference, error) {
	return parseLocalReference(PCodeDestinationLocalReference, b)
}

// ParseSourceLocalReference parses the given byte sequence as a Source Local Reference.
func ParseSourceLocalReference(b []byte) (*LocalReference, error) {
	return parseLocalReference(PCodeSourceLocalReference, b)
}

func parseLocalReference(c ParameterNameCode, b []byte) (*LocalReference, error) {
	l := &LocalReference{
		code: c,
	}

	if _, err := l.Read(b); err != nil {
		return nil, err
	}

	return l, nil
}

// Read sets the values retrieved from byte sequence in a LocalReference.
func (l *LocalReference) Read(b []byte) (int, error) {
	n := 3
	if len(b) < n {
		return 0, io.ErrUnexpectedEOF
	}

	// code must be set by the calle, or use ParseDestination/SourceLocalReference.
	l.length = n
	l.value = b[:n]
	return n, nil
}

// Write serializes the LocalReference parameter and returns it as a byte slice.
func (l *LocalReference) Write(b []byte) (int, error) {
	if len(b) < l.length {
		return 0, io.ErrUnexpectedEOF
	}

	copy(b, l.value)
	return l.length, nil
}

// Value returns the LocalReference in []byte.
func (l *LocalReference) Value() []byte {
	return l.value
}

// Uint32 returns the LocalReference in uint32.
func (l *LocalReference) Uint32() uint32 {
	return utils.Uint24To32(l.value)
}

// Code returns the LocalReference in ParameterNameCode.
func (l *LocalReference) Code() ParameterNameCode {
	return l.code
}

// String returns the LocalReference in string.
func (l *LocalReference) String() string {
	if l.code == PCodeDestinationLocalReference || l.code == PCodeSourceLocalReference {
		return fmt.Sprintf("%s: %d", l.code, l.Uint32())
	}
	return fmt.Sprintf("%s: %x", "(Destination or Source) local reference", l.value)
}

// PartyAddress is a SCCP parameter that represents a Called/Calling Party Address.
type PartyAddress struct {
	paramType ParameterType
	code      ParameterNameCode
	length    int

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

// NewPartyAddress creates a new PartyAddress from properly-typed values.
//
// The given SPC and SSN are set to 0 if the corresponding bit is not properly set in the
// AddressIndicator. Use NewAddressIndicator to create a proper AddressIndicator.
func NewPartyAddress(ai uint8, spc uint16, ssn uint8, gt *GlobalTitle) *PartyAddress {
	p := &PartyAddress{
		paramType:   PTypeV,
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

// NewCalledPartyAddress creates a new PartyAddress for Called Party Address.
func NewCalledPartyAddress(ai uint8, spc uint16, ssn uint8, gt *GlobalTitle) *PartyAddress {
	p := NewPartyAddress(ai, spc, ssn, gt)
	p.code = PCodeCalledPartyAddress
	return p
}

// NewCallingPartyAddress creates a new PartyAddress for Calling Party Address.
func NewCallingPartyAddress(ai uint8, spc uint16, ssn uint8, gt *GlobalTitle) *PartyAddress {
	p := NewPartyAddress(ai, spc, ssn, gt)
	p.code = PCodeCallingPartyAddress
	return p
}

// AsOptional makes the PartyAddress an optional parameter.
func (p *PartyAddress) AsOptional() *PartyAddress {
	if p.paramType != PTypeO {
		p.paramType = PTypeO
		p.length += 2
	}
	return p
}

// ParseCalledPartyAddress parses the given byte sequence as a mandatory fixed length
// Called Party Address and returns it as a PartyAddress.
func ParseCalledPartyAddress(b []byte) (*PartyAddress, error) {
	return parsePartyAddress(PTypeV, PCodeCalledPartyAddress, b)
}

// ParseCallingPartyAddress parses the given byte sequence as a mandatory fixed length
// Calling Party Address and returns it as a PartyAddress.
func ParseCallingPartyAddress(b []byte) (*PartyAddress, error) {
	return parsePartyAddress(PTypeV, PCodeCallingPartyAddress, b)
}

// ParseCalledPartyAddressOptional parses the given byte sequence as an optional
// Called Party Address and returns it as a PartyAddress.
func ParseCalledPartyAddressOptional(b []byte) (*PartyAddress, error) {
	return parsePartyAddress(PTypeO, PCodeCalledPartyAddress, b)
}

// ParseCallingPartyAddressOptional parses the given byte sequence as an optional
// Calling Party Address and returns it as a PartyAddress.
func ParseCallingPartyAddressOptional(b []byte) (*PartyAddress, error) {
	return parsePartyAddress(PTypeO, PCodeCallingPartyAddress, b)
}

func parsePartyAddress(ptype ParameterType, code ParameterNameCode, b []byte) (*PartyAddress, error) {
	p := &PartyAddress{
		paramType: ptype,
		code:      code,
	}

	if _, err := p.Read(b); err != nil {
		return nil, err
	}

	return p, nil
}

// Read sets the values retrieved from byte sequence in a PartyAddress.
func (p *PartyAddress) Read(b []byte) (int, error) {
	if p.paramType == PTypeV {
		return p.read(b)
	}
	return p.readOptional(b)
}

func (p *PartyAddress) read(b []byte) (int, error) {
	var n = 2
	if len(b) < n {
		return 0, io.ErrUnexpectedEOF
	}

	p.length = int(b[0])
	p.Indicator = b[1]

	if int(p.length) != len(b)-1 {
		return n, io.ErrUnexpectedEOF
	}

	if p.HasPC() {
		end := n + 2
		if end >= len(b) {
			return n, io.ErrUnexpectedEOF
		}
		p.SignalingPointCode = binary.BigEndian.Uint16(b[n:end])
		n = end
	}

	if p.HasSSN() {
		p.SubsystemNumber = b[n]
		n++
	}

	gti := p.GTI()
	if gti == 0 {
		return n, nil
	}

	p.GlobalTitle = &GlobalTitle{GTI: gti}
	m, err := p.GlobalTitle.Read(b[n : int(p.length)+1])
	if err != nil {
		return n + m, err
	}
	n += m

	return n, nil
}

func (p *PartyAddress) readOptional(b []byte) (int, error) {
	n := 3
	if len(b) < n {
		return 0, io.ErrUnexpectedEOF
	}

	p.code = ParameterNameCode(b[0])
	if p.code != PCodeCalledPartyAddress && p.code != PCodeCallingPartyAddress {
		logf(
			"invalid parameter code: expected %d or %d, got %d",
			PCodeCalledPartyAddress, PCodeCallingPartyAddress, p.code,
		)
	}

	return p.read(b[1:])
}

// Write serializes the PartyAddress parameter and returns it as a byte slice.
func (p *PartyAddress) Write(b []byte) (int, error) {
	if p.paramType == PTypeV {
		return p.write(b)
	}
	return p.writeOptional(b)
}

func (p *PartyAddress) write(b []byte) (int, error) {
	if len(b) < p.MarshalLen() {
		return 0, io.ErrUnexpectedEOF
	}

	b[0] = uint8(p.length)
	b[1] = p.Indicator

	var n = 2
	if p.HasPC() {
		binary.BigEndian.PutUint16(b[n:n+2], p.SignalingPointCode)
		n += 2
	}

	if p.HasSSN() {
		b[n] = p.SubsystemNumber
		n++
	}

	if p.GlobalTitle != nil {
		m, err := p.GlobalTitle.Write(b[n : n+p.GlobalTitle.MarshalLen()])
		if err != nil {
			return n + m, err
		}
		n += m
	}

	return n, nil
}

func (p *PartyAddress) writeOptional(b []byte) (int, error) {
	if len(b) < p.MarshalLen() {
		return 0, io.ErrUnexpectedEOF
	}

	b[0] = uint8(p.code)
	n, err := p.write(b[1:])
	if err != nil {
		return n + 1, err
	}

	return n + 1, nil
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
	if _, err := p.Write(b); err != nil {
		return err
	}
	return nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a SCCP common header.
func (p *PartyAddress) UnmarshalBinary(b []byte) error {
	if _, err := p.Read(b); err != nil {
		return err
	}
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

// SetLength sets the length in length field.
func (p *PartyAddress) SetLength() {
	p.length = p.MarshalLen() - 1
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

// Length returns the length of PartyAddress.
func (p *PartyAddress) Length() int {
	return p.length
}

// Code returns the PartyAddress in ParameterNameCode.
func (p *PartyAddress) Code() ParameterNameCode {
	return p.code
}

// String returns the PartyAddress values in human readable format.
func (p *PartyAddress) String() string {
	return fmt.Sprintf("%s (%s): {length: %d, Indicator: %#08b, SignalingPointCode: %d, SubsystemNumber: %d, GlobalTitle: %v}",
		p.code, p.paramType, p.length, p.Indicator, p.SignalingPointCode, p.SubsystemNumber, p.GlobalTitle,
	)
}

// ProtocolClass is a Protocol Class SCCP parameter.
type ProtocolClass struct {
	paramType ParameterType
	code      ParameterNameCode
	length    int
	value     uint8
}

// NewProtocolClass creates a new ProtocolClass.
func NewProtocolClass(cls int, returnOnError bool) *ProtocolClass {
	p := &ProtocolClass{
		paramType: PTypeF,
		code:      PCodeProtocolClass,
		length:    1,
		value:     uint8(cls),
	}

	if returnOnError {
		p.value = uint8(cls | 0x80)
	}

	return p
}

// Read sets the values retrieved from byte sequence in a ProtocolClass.
func (p *ProtocolClass) Read(b []byte) (int, error) {
	n := 1
	if len(b) < n {
		return 0, io.ErrUnexpectedEOF
	}

	p.code = PCodeProtocolClass
	p.length = n
	p.value = b[0]

	return n, nil
}

// Write serializes the ProtocolClass parameter and returns it as a byte slice.
func (p *ProtocolClass) Write(b []byte) (int, error) {
	if len(b) < p.length {
		return 0, io.ErrUnexpectedEOF
	}

	b[0] = p.value
	return p.length, nil
}

// Value returns the ProtocolClass in uint8.
func (p *ProtocolClass) Value() uint8 {
	return uint8(p.value)
}

// Class returns the class part from ProtocolClass parameter.
func (p *ProtocolClass) Class() int {
	return int(p.value) & 0xf
}

// ReturnOnError judges if ProtocolClass has "Return Message On Error" option.
func (p *ProtocolClass) ReturnOnError() bool {
	return (int(p.value) >> 7) == 1
}

// Code returns the ProtocolClass in ParameterNameCode.
func (p *ProtocolClass) Code() ParameterNameCode {
	return p.code
}

// String returns the ProtocolClass in string.
func (p *ProtocolClass) String() string {
	return fmt.Sprintf("%s: {Class: %d, ReturnOnError: %v}", p.code, p.Class(), p.ReturnOnError())
}

// SegmentingReassembling represents the Segmenting/Reassembling.
type SegmentingReassembling struct {
	paramType ParameterType
	code      ParameterNameCode
	length    int
	value     uint8
}

// NewSegmentingReassembling creates a new SegmentingReassembling.
func NewSegmentingReassembling(moreData bool) *SegmentingReassembling {
	v := uint8(0)
	if moreData {
		v = 1
	}

	return &SegmentingReassembling{
		paramType: PTypeF,
		code:      PCodeSegmentingReassembling,
		length:    1,
		value:     v,
	}
}

// Read sets the values retrieved from byte sequence in a SegmentingReassembling.
func (s *SegmentingReassembling) Read(b []byte) (int, error) {
	n := 1
	if len(b) < n {
		return 0, io.ErrUnexpectedEOF
	}

	s.code = PCodeSegmentingReassembling
	s.length = n
	s.value = b[0]

	return n, nil
}

// Write serializes the SegmentingReassembling parameter and returns it as a byte slice.
func (s *SegmentingReassembling) Write(b []byte) (int, error) {
	if len(b) < s.length {
		return 0, io.ErrUnexpectedEOF
	}

	b[0] = s.value
	return s.length, nil
}

// Value returns the SegmentingReassembling in uint8.
func (s *SegmentingReassembling) Value() uint8 {
	return s.value
}

// MoreData judges if the message has more data.
func (s *SegmentingReassembling) MoreData() bool {
	return s.value == 1
}

// Code returns the SegmentingReassembling in ParameterNameCode.
func (s *SegmentingReassembling) Code() ParameterNameCode {
	return s.code
}

// String returns the SegmentingReassembling in string.
func (s *SegmentingReassembling) String() string {
	return fmt.Sprintf("%s: %d", s.code, s.value)
}

// ReceiveSequenceNumber represents the Receive Sequence Number.
type ReceiveSequenceNumber struct {
	paramType ParameterType
	code      ParameterNameCode
	length    int
	value     uint8
}

// NewReceiveSequenceNumber creates a new ReceiveSequenceNumber.
func NewReceiveSequenceNumber(v uint8) *ReceiveSequenceNumber {
	return &ReceiveSequenceNumber{
		paramType: PTypeF,
		code:      PCodeReceiveSequenceNumber,
		length:    1,
		value:     v & 0b01111111,
	}
}

// Read sets the values retrieved from byte sequence in a ReceiveSequenceNumber.
func (r *ReceiveSequenceNumber) Read(b []byte) (int, error) {
	n := 1
	if len(b) < n {
		return 0, io.ErrUnexpectedEOF
	}

	r.code = PCodeReceiveSequenceNumber
	r.length = n
	r.value = b[0] & 0b11111110

	return n, nil
}

// Write serializes the ReceiveSequenceNumber parameter and returns it as a byte slice.
func (r *ReceiveSequenceNumber) Write(b []byte) (int, error) {
	if len(b) < r.length {
		return 0, io.ErrUnexpectedEOF
	}

	b[0] = r.value
	return r.length, nil
}

// Value returns the ReceiveSequenceNumber in uint8.
func (r *ReceiveSequenceNumber) Value() uint8 {
	return r.value
}

// Code returns the ReceiveSequenceNumber in ParameterNameCode.
func (r *ReceiveSequenceNumber) Code() ParameterNameCode {
	return r.code
}

// String returns the ReceiveSequenceNumber in string.
func (r *ReceiveSequenceNumber) String() string {
	return fmt.Sprintf("%s: %d", r.code, r.value)
}

// SequencingSegmenting represents the Sequencing/Segmenting.
type SequencingSegmenting struct {
	paramType             ParameterType
	code                  ParameterNameCode
	length                int
	SendSequenceNumber    uint8
	ReceiveSequenceNumber uint8
	MoreData              bool
}

// NewSequencingSegmenting creates a new SequencingSegmenting.
func NewSequencingSegmenting(snd, rcv uint8, moreData bool) *SequencingSegmenting {
	return &SequencingSegmenting{
		paramType:             PTypeF,
		code:                  PCodeSequencingSegmenting,
		length:                2,
		SendSequenceNumber:    snd & 0b01111111,
		ReceiveSequenceNumber: rcv & 0b01111111,
		MoreData:              moreData,
	}
}

// Read sets the values retrieved from byte sequence in a SequencingSegmenting.
func (s *SequencingSegmenting) Read(b []byte) (int, error) {
	n := 2
	if len(b) < n {
		return 0, io.ErrUnexpectedEOF
	}

	s.code = PCodeSequencingSegmenting
	s.length = n

	s.SendSequenceNumber = b[0] & 0b11111110
	s.ReceiveSequenceNumber = b[1] & 0b11111110
	s.MoreData = b[1]&0b00000001 == 1

	return n, nil
}

// Write serializes the SequencingSegmenting parameter and returns it as a byte slice.
func (s *SequencingSegmenting) Write(b []byte) (int, error) {
	if len(b) < s.length {
		return 0, io.ErrUnexpectedEOF
	}

	b[0] = s.SendSequenceNumber
	b[1] = s.ReceiveSequenceNumber

	if s.MoreData {
		b[1] |= 0b00000001
	}

	return s.length, nil
}

// Value returns the SequencingSegmenting in *SequencingSegmenting.
// This is here to implement the Parameter interface, but it won't be used.
func (s *SequencingSegmenting) Value() *SequencingSegmenting {
	return s
}

// Code returns the SequencingSegmenting in ParameterNameCode.
func (s *SequencingSegmenting) Code() ParameterNameCode {
	return s.code
}

// String returns the SequencingSegmenting in string.
func (s *SequencingSegmenting) String() string {
	return fmt.Sprintf(
		"%s: {SendSequenceNumber=%d, ReceiveSequenceNumber=%d, MoreData=%t}",
		s.code, s.SendSequenceNumber, s.ReceiveSequenceNumber, s.MoreData,
	)
}

// Credit represents the Credit.
type Credit struct {
	paramType ParameterType
	code      ParameterNameCode
	length    int
	value     uint8
}

// NewCredit creates a new Credit.
func NewCredit(v uint8) *Credit {
	return &Credit{
		paramType: PTypeF,
		code:      PCodeCredit,
		length:    1,
		value:     v,
	}
}

// AsOptional makes the Credit an optional parameter.
func (c *Credit) AsOptional() *Credit {
	c.paramType = PTypeO
	c.length = 3
	return c
}

// Read sets the values retrieved from byte sequence in a Credit.
func (c *Credit) Read(b []byte) (int, error) {
	if c.paramType == PTypeO {
		return c.readOptional(b)
	}
	return c.read(b)
}

// read sets the values retrieved from byte sequence in a Credit.
func (c *Credit) read(b []byte) (int, error) {
	n := 1
	if len(b) < n {
		return 0, io.ErrUnexpectedEOF
	}

	c.code = PCodeCredit
	c.length = n
	c.value = b[0]

	return n, nil
}

func (c *Credit) readOptional(b []byte) (int, error) {
	n := 3
	if len(b) < n {
		return 0, nil
	}

	c.code = ParameterNameCode(b[0])
	if c.code != PCodeCredit {
		logf("invalid parameter code: expected %d, got %d", PCodeCredit, c.code)
	}

	c.length = int(b[1])
	if c.length != n {
		logf("%s: invalid length: expected %d, got %d", PCodeCredit, n, c.length)
	}

	c.value = b[2]
	return n, nil
}

// Write serializes the Credit parameter and returns it as a byte slice.
func (c *Credit) Write(b []byte) (int, error) {
	if c.paramType == PTypeO {
		return c.writeOptional(b)
	}
	return c.write(b)
}

// write serializes the Credit parameter and returns it as a byte slice.
func (c *Credit) write(b []byte) (int, error) {
	if len(b) < c.length {
		return 0, io.ErrUnexpectedEOF
	}

	b[0] = c.value

	return c.length, nil
}

func (c *Credit) writeOptional(b []byte) (int, error) {
	if len(b) < c.length {
		return 0, io.ErrUnexpectedEOF
	}

	b[0] = uint8(c.code)
	b[1] = uint8(c.length)
	b[2] = c.value

	return c.length, nil
}

// Value returns the Credit in uint8.
func (c *Credit) Value() uint8 {
	return c.value
}

// Code returns the Credit in ParameterNameCode.
func (c *Credit) Code() ParameterNameCode {
	return c.code
}

// String returns the Credit in string.
func (c *Credit) String() string {
	return fmt.Sprintf("%s: %d", c.code, c.value)
}

// ReleaseCause represents the Release Cause.
type ReleaseCause struct {
	paramType ParameterType
	code      ParameterNameCode
	length    int
	value     ReleaseCauseValue
}

// ReleaseCauseValue is a type for ReleaseCause.
type ReleaseCauseValue uint8

// ReleaseCauseValue values.
const (
	ReleaseCauseEndUserOriginated                  ReleaseCauseValue = 0b00000000 // end user originated
	ReleaseCauseEndUserCongestion                  ReleaseCauseValue = 0b00000001 // end user congestion
	ReleaseCauseEndUserFailure                     ReleaseCauseValue = 0b00000010 // end user failure
	ReleaseCauseSCCPUserOriginated                 ReleaseCauseValue = 0b00000011 // SCCP user originated
	ReleaseCauseRemoteProcedureError               ReleaseCauseValue = 0b00000100 // remote procedure error
	ReleaseCauseInconsistentConnectionData         ReleaseCauseValue = 0b00000101 // inconsistent connection data
	ReleaseCauseAccessFailure                      ReleaseCauseValue = 0b00000110 // access failure
	ReleaseCauseAccessCongestion                   ReleaseCauseValue = 0b00000111 // access congestion
	ReleaseCauseSubsystemFailure                   ReleaseCauseValue = 0b00001000 // subsystem failure
	ReleaseCauseSubsystemCongestion                ReleaseCauseValue = 0b00001001 // subsystem congestion
	ReleaseCauseMTPFailure                         ReleaseCauseValue = 0b00001010 // MTP failure
	ReleaseCauseNetworkCongestion                  ReleaseCauseValue = 0b00001011 // network congestion
	ReleaseCauseExpirationOfResetTimer             ReleaseCauseValue = 0b00001100 // expiration of reset timer
	ReleaseCauseExpirationOfReceiveInactivityTimer ReleaseCauseValue = 0b00001101 // expiration of receive inactivity timer
	_                                              ReleaseCauseValue = 0b00001110 // reserved
	ReleaseCauseUnqualified                        ReleaseCauseValue = 0b00001111 // unqualified
	ReleaseCauseSCCPFailure                        ReleaseCauseValue = 0b00010000 // SCCP failure
)

// NewReleaseCause creates a new ReleaseCause.
func NewReleaseCause(v ReleaseCauseValue) *ReleaseCause {
	return &ReleaseCause{
		paramType: PTypeF,
		code:      PCodeReleaseCause,
		length:    1,
		value:     v,
	}
}

// Read sets the values retrieved from byte sequence in a ReleaseCause.
func (r *ReleaseCause) Read(b []byte) (int, error) {
	n := 1
	if len(b) < n {
		return 0, io.ErrUnexpectedEOF
	}

	r.code = PCodeReleaseCause
	r.length = n
	r.value = ReleaseCauseValue(b[0])

	return n, nil
}

// Write serializes the ReleaseCause parameter and returns it as a byte slice.
func (r *ReleaseCause) Write(b []byte) (int, error) {
	if len(b) < r.length {
		return 0, io.ErrUnexpectedEOF
	}

	b[0] = uint8(r.value)
	return r.length, nil
}

// Value returns the ReleaseCause in ReleaseCauseValue.
func (r *ReleaseCause) Value() ReleaseCauseValue {
	return r.value
}

// Code returns the ReleaseCause in ParameterNameCode.
func (r *ReleaseCause) Code() ParameterNameCode {
	return r.code
}

// String returns the ReleaseCause in string.
func (r *ReleaseCause) String() string {
	return fmt.Sprintf("%s: %s", r.code, r.value)
}

// ReturnCause represents the Return Cause.
type ReturnCause struct {
	paramType ParameterType
	code      ParameterNameCode
	length    int
	value     ReturnCauseValue
}

// ReturnCauseValue is a type for ReturnCause.
type ReturnCauseValue uint8

// ReturnCauseValue values.
const (
	ReturnCauseNoTranslationForAnAddressOfSuchNature ReturnCauseValue = 0b00000000 // no translation for an address of such nature
	ReturnCauseNoTranslationForThisSpecificAddress   ReturnCauseValue = 0b00000001 // no translation for this specific address
	ReturnCauseSubsystemCongestion                   ReturnCauseValue = 0b00000010 // subsystem congestion
	ReturnCauseSubsystemFailure                      ReturnCauseValue = 0b00000011 // subsystem failure
	ReturnCauseUnequippedUser                        ReturnCauseValue = 0b00000100 // unequipped user
	ReturnCauseMTPFailure                            ReturnCauseValue = 0b00000101 // MTP failure
	ReturnCauseNetworkCongestion                     ReturnCauseValue = 0b00000110 // network congestion
	ReturnCauseUnqualified                           ReturnCauseValue = 0b00000111 // unqualified
	ReturnCauseErrorInMessageTransport               ReturnCauseValue = 0b00001000 // error in message transport
	ReturnCauseErrorInLocalProcessing                ReturnCauseValue = 0b00001001 // error in local processing
	ReturnCauseDestinationCannotPerformReassembly    ReturnCauseValue = 0b00001010 // destination cannot perform reassembly
	ReturnCauseSCCPFailure                           ReturnCauseValue = 0b00001011 // SCCP failure
	ReturnCauseHopCounterViolation                   ReturnCauseValue = 0b00001100 // hop counter violation
	ReturnCauseSegmentationNotSupported              ReturnCauseValue = 0b00001101 // segmentation not supported
	ReturnCauseSegmentationFailure                   ReturnCauseValue = 0b00001110 // segmentation failure
)

// NewReturnCause creates a new ReturnCause.
func NewReturnCause(v ReturnCauseValue) *ReturnCause {
	return &ReturnCause{
		paramType: PTypeF,
		code:      PCodeReturnCause,
		length:    1,
		value:     v,
	}
}

// Read sets the values retrieved from byte sequence in a ReturnCause.
func (r *ReturnCause) Read(b []byte) (int, error) {
	n := 1
	if len(b) < n {
		return 0, io.ErrUnexpectedEOF
	}

	r.code = PCodeReturnCause
	r.length = n
	r.value = ReturnCauseValue(b[0])

	return n, nil
}

// Write serializes the ReturnCause parameter and returns it as a byte slice.
func (r *ReturnCause) Write(b []byte) (int, error) {
	if len(b) < r.length {
		return 0, io.ErrUnexpectedEOF
	}

	b[0] = uint8(r.value)
	return r.length, nil
}

// Value returns the ReturnCause in ReturnCauseValue.
func (r *ReturnCause) Value() ReturnCauseValue {
	return r.value
}

// Code returns the ReturnCause in ParameterNameCode.
func (r *ReturnCause) Code() ParameterNameCode {
	return r.code
}

// String returns the ReturnCause in string.
func (r *ReturnCause) String() string {
	return fmt.Sprintf("%s: %s", r.code, r.value)
}

// ResetCause represents the Reset Cause.
type ResetCause struct {
	paramType ParameterType
	code      ParameterNameCode
	length    int
	value     ResetCauseValue
}

// ResetCauseValue is a type for ResetCause.
type ResetCauseValue uint8

// ResetCauseValue values.
const (
	ResetCauseEndUserOriginated                                                    ResetCauseValue = 0b00000000 // end user originated
	ResetCauseSCCPUserOriginated                                                   ResetCauseValue = 0b00000001 // SCCP user originated
	ResetCauseMessageOutOfOrderIncorrectSendSequenceNumber                         ResetCauseValue = 0b00000010 // message out of order - incorrect P(S)
	ResetCauseMessageOutOfOrderIncorrectReceiveSequenceNumber                      ResetCauseValue = 0b00000011 // message out of order - incorrect P(R)
	ResetCauseRemoteProcedureErrorMessageOutOfWindow                               ResetCauseValue = 0b00000100 // remote procedure error - message out of window
	ResetCauseRemoteProcedureErrorIncorrectSendSequenceNumberAfterReinitialization ResetCauseValue = 0b00000101 // remote procedure error - incorrect P(S) after (re)initialization
	ResetCauseRemoteProcedureErrorGeneral                                          ResetCauseValue = 0b00000110 // remote procedure error - general
	ResetCauseRemoteEndUserOperational                                             ResetCauseValue = 0b00000111 // remote end user operational
	ResetCauseNetworkOperational                                                   ResetCauseValue = 0b00001000 // network operational
	ResetCauseAccessOperational                                                    ResetCauseValue = 0b00001001 // access operational
	ResetCauseNetworkCongestion                                                    ResetCauseValue = 0b00001010 // network congestion
	_                                                                              ResetCauseValue = 0b00001011
	ResetCauseUnqualified                                                          ResetCauseValue = 0b00001100 // unqualified
)

// NewResetCause creates a new ResetCause.
func NewResetCause(v ResetCauseValue) *ResetCause {
	return &ResetCause{
		paramType: PTypeF,
		code:      PCodeResetCause,
		length:    1,
		value:     v,
	}
}

// Read sets the values retrieved from byte sequence in a ResetCause.
func (r *ResetCause) Read(b []byte) (int, error) {
	n := 1
	if len(b) < n {
		return 0, io.ErrUnexpectedEOF
	}

	r.code = PCodeResetCause
	r.length = n
	r.value = ResetCauseValue(b[0])

	return n, nil
}

// Write serializes the ResetCause parameter and returns it as a byte slice.
func (r *ResetCause) Write(b []byte) (int, error) {
	if len(b) < r.length {
		return 0, io.ErrUnexpectedEOF
	}

	b[0] = uint8(r.value)
	return r.length, nil
}

// Value returns the ResetCause in ResetCauseValue.
func (r *ResetCause) Value() ResetCauseValue {
	return r.value
}

// Code returns the ResetCause in ParameterNameCode.
func (r *ResetCause) Code() ParameterNameCode {
	return r.code
}

// String returns the ResetCause in string.
func (r *ResetCause) String() string {
	return fmt.Sprintf("%s: %s", r.code, r.value)
}

// ErrorCause represents the Error Cause.
type ErrorCause struct {
	paramType ParameterType
	code      ParameterNameCode
	length    int
	value     ErrorCauseValue
}

// ErrorCauseValue is a type for ErrorCause.
type ErrorCauseValue uint8

// ErrorCauseValue values.
const (
	ErrorCauseLocalReferenceNumberMismatchUnassignedDestinationLRN ErrorCauseValue = 0b00000000 // local reference number (LRN) mismatch - unassigned destination LRN
	ErrorCauseLocalReferenceNumberMismatchInconsistentSourceLRN    ErrorCauseValue = 0b00000001 // local reference number (LRN) mismatch - inconsistent source LRN
	ErrorCausePointCodeMismatch                                    ErrorCauseValue = 0b00000010 // point code mismatch
	ErrorCauseServiceClassMismatch                                 ErrorCauseValue = 0b00000011 // service class mismatch
	ErrorCauseUnqualified                                          ErrorCauseValue = 0b00000100 // unqualified
)

// NewErrorCause creates a new ErrorCause.
func NewErrorCause(v ErrorCauseValue) *ErrorCause {
	return &ErrorCause{
		paramType: PTypeF,
		code:      PCodeErrorCause,
		length:    1,
		value:     v,
	}
}

// Read sets the values retrieved from byte sequence in a ErrorCause.
func (e *ErrorCause) Read(b []byte) (int, error) {
	n := 1
	if len(b) < n {
		return 0, io.ErrUnexpectedEOF
	}

	e.code = PCodeErrorCause
	e.length = n
	e.value = ErrorCauseValue(b[0])

	return n, nil
}

// Write serializes the ErrorCause parameter and returns it as a byte slice.
func (e *ErrorCause) Write(b []byte) (int, error) {
	if len(b) < e.length {
		return 0, io.ErrUnexpectedEOF
	}

	b[0] = uint8(e.value)
	return e.length, nil
}

// Value returns the ErrorCause in ErrorCauseValue.
func (e *ErrorCause) Value() ErrorCauseValue {
	return e.value
}

// Code returns the ErrorCause in ParameterNameCode.
func (e *ErrorCause) Code() ParameterNameCode {
	return e.code
}

// String returns the ErrorCause in string.
func (e *ErrorCause) String() string {
	return fmt.Sprintf("%s: %s", e.code, e.value)
}

// RefusalCause represents the Refusal Cause.
type RefusalCause struct {
	paramType ParameterType
	code      ParameterNameCode
	length    int
	value     RefusalCauseValue
}

// RefusalCauseValue is a type for RefusalCause.
type RefusalCauseValue uint8

// RefusalCauseValue values.
const (
	RefusalCauseEndUserOriginated                           RefusalCauseValue = 0b00000000 // end user originated
	RefusalCauseEndUserCongestion                           RefusalCauseValue = 0b00000001 // end user congestion
	RefusalCauseEndUserFailure                              RefusalCauseValue = 0b00000010 // end user failure
	RefusalCauseSCCPUserOriginated                          RefusalCauseValue = 0b00000011 // SCCP user originated
	RefusalCauseDestinationAddressUnknown                   RefusalCauseValue = 0b00000100 // destination address unknown
	RefusalCauseDestinationInaccessible                     RefusalCauseValue = 0b00000101 // destination inaccessible
	RefusalCauseNetworkResourceQoSNotAvailableNonTransient  RefusalCauseValue = 0b00000110 // network resource - QoS not available/non-transient
	RefusalCauseNetworkResourceQoSNotAvailableTransient     RefusalCauseValue = 0b00000111 // network resource - QoS not available/transient
	RefusalCauseAccessFailure                               RefusalCauseValue = 0b00001000 // access failure
	RefusalCauseAccessCongestion                            RefusalCauseValue = 0b00001001 // access congestion
	RefusalCauseSubsystemFailure                            RefusalCauseValue = 0b00001010 // subsystem failure
	RefusalCauseSubsystemCongestion                         RefusalCauseValue = 0b00001011 // subsystem congestion
	RefusalCauseExpirationOfTheConnectionEstablishmentTimer RefusalCauseValue = 0b00001100 // expiration of the connection establishment timer
	RefusalCauseIncompatibleUserData                        RefusalCauseValue = 0b00001101 // incompatible user data
	_                                                       RefusalCauseValue = 0b00001110 // reserved
	RefusalCauseUnqualified                                 RefusalCauseValue = 0b00001111 // unqualified
	RefusalCauseHopCounterViolation                         RefusalCauseValue = 0b00010000 // hop counter violation
	RefusalCauseSCCPFailure                                 RefusalCauseValue = 0b00010001 // SCCP failure
	RefusalCauseNoTranslationForAnAddressOfSuchNature       RefusalCauseValue = 0b00010010 // no translation for an address of such nature
	RefusalCauseUnequippedUser                              RefusalCauseValue = 0b00010011 // unequipped user
)

// NewRefusalCause creates a new RefusalCause.
func NewRefusalCause(v RefusalCauseValue) *RefusalCause {
	return &RefusalCause{
		paramType: PTypeF,
		code:      PCodeRefusalCause,
		length:    1,
		value:     v,
	}
}

// Read sets the values retrieved from byte sequence in a RefusalCause.
func (r *RefusalCause) Read(b []byte) (int, error) {
	n := 1
	if len(b) < n {
		return 0, io.ErrUnexpectedEOF
	}

	r.code = PCodeRefusalCause
	r.length = n
	r.value = RefusalCauseValue(b[0])

	return n, nil
}

// Write serializes the RefusalCause parameter and returns it as a byte slice.
func (r *RefusalCause) Write(b []byte) (int, error) {
	if len(b) < r.length {
		return 0, io.ErrUnexpectedEOF
	}

	b[0] = uint8(r.value)
	return r.length, nil
}

// Value returns the RefusalCause in RefusalCauseValue.
func (r *RefusalCause) Value() RefusalCauseValue {
	return r.value
}

// Code returns the RefusalCause in ParameterNameCode.
func (r *RefusalCause) Code() ParameterNameCode {
	return r.code
}

// String returns the RefusalCause in string.
func (r *RefusalCause) String() string {
	return fmt.Sprintf("%s: %s", r.code, r.value)
}

// Data represents the Data.
type Data struct {
	paramType ParameterType
	code      ParameterNameCode
	length    int
	value     []byte
}

// NewData creates a new Data.
func NewData(v []byte) *Data {
	return &Data{
		paramType: PTypeV,
		code:      PCodeData,
		length:    len(v),
		value:     v,
	}
}

// AsOptional makes the Data an optional parameter.
func (d *Data) AsOptional() *Data {
	if d.paramType != PTypeO {
		d.paramType = PTypeO
		d.length += 2
	}
	return d
}

// Read sets the values retrieved from byte sequence in a Data.
func (d *Data) Read(b []byte) (int, error) {
	if d.paramType == PTypeO {
		return d.readOptional(b)
	}
	return d.read(b)
}

// read sets the values retrieved from byte sequence in a Data.
func (d *Data) read(b []byte) (int, error) {
	n := len(b)

	d.code = PCodeData
	d.length = n
	d.value = b[:n]

	return n, nil
}

func (d *Data) readOptional(b []byte) (int, error) {
	n := len(b)

	d.code = ParameterNameCode(b[0])
	if d.code != PCodeData {
		logf("invalid parameter code: expected %d, got %d", PCodeData, d.code)
	}

	d.length = int(b[1])
	if d.length != n {
		logf("%s: invalid length: expected %d, got %d", PCodeData, n, d.length)
	}

	d.value = b[2:]

	return n, nil
}

// Write serializes the Data parameter and returns it as a byte slice.
func (d *Data) Write(b []byte) (int, error) {
	if d.paramType == PTypeO {
		return d.writeOptional(b)
	}
	return d.write(b)
}

// write serializes the Data parameter and returns it as a byte slice.
func (d *Data) write(b []byte) (int, error) {
	if len(b) < d.length {
		return 0, io.ErrUnexpectedEOF
	}

	copy(b, d.value)
	return d.length, nil
}

func (d *Data) writeOptional(b []byte) (int, error) {
	if len(b) < d.length {
		return 0, io.ErrUnexpectedEOF
	}

	b[0] = uint8(d.code)
	b[1] = uint8(d.length)
	copy(b[2:], d.value)
	return d.length, nil
}

// Value returns the Data in []byte.
func (d *Data) Value() []byte {
	return d.value
}

// Code returns the Data in ParameterNameCode.
func (d *Data) Code() ParameterNameCode {
	return d.code
}

// String returns the Data in string.
func (d *Data) String() string {
	return fmt.Sprintf("%s: %x", d.code, d.value)
}

// Segmentation represents the Segmentation.
type Segmentation struct {
	paramType         ParameterType
	code              ParameterNameCode
	length            int
	FirstSegment      bool
	Class             uint8
	RemainingSegments uint8
	LocalReference    uint32 // 3-octet
}

// NewSegmentation creates a new Segmentation.
func NewSegmentation(first bool, cls, rem uint8, lrn uint32) *Segmentation {
	return &Segmentation{
		paramType:         PTypeO,
		code:              PCodeSegmentation,
		length:            6,
		FirstSegment:      first,
		Class:             cls & 0b1,
		RemainingSegments: rem & 0b1111,
		LocalReference:    lrn & 0x00ffffff,
	}
}

// AsOptional makes the Segmentation an optional parameter.
func (s *Segmentation) AsOptional() *Segmentation {
	s.paramType = PTypeO
	s.length = 6
	return s
}

// Read sets the values retrieved from byte sequence in a Segmentation.
func (s *Segmentation) Read(b []byte) (int, error) {
	if s.paramType != PTypeO {
		logf("Segmentation parameter must be optional: %v", s)
	}

	n := 6
	if len(b) < n {
		return 0, io.ErrUnexpectedEOF
	}

	s.code = ParameterNameCode(b[0])
	if s.code != PCodeSegmentation {
		logf("invalid parameter code: expected %d, got %d", PCodeSegmentation, s.code)
	}

	s.length = int(b[1])
	if s.length != n {
		logf("%s: invalid length: expected %d, got %d", PCodeSegmentation, n, s.length)
	}

	s.FirstSegment = b[2]&0b10000000 == 0b10000000
	s.Class = b[2] & 0b01000000
	s.RemainingSegments = b[2] & 0b1111
	s.LocalReference = utils.Uint24To32(b[3:])

	return n, nil
}

// Write serializes the Segmentation parameter and returns it as a byte slice.
func (s *Segmentation) Write(b []byte) (int, error) {
	if s.paramType != PTypeO {
		logf("Segmentation parameter must be optional: %v", s)
	}

	if len(b) < s.length {
		return 0, io.ErrUnexpectedEOF
	}

	b[0] = uint8(s.code)
	b[1] = uint8(s.length)

	if s.FirstSegment {
		b[2] |= 0b10000000
	}

	b[2] |= s.Class & 0b1 << 6
	b[2] |= s.RemainingSegments & 0b111

	copy(b[3:], utils.Uint32To24(s.LocalReference))

	return s.length, nil
}

// Value returns the Segmentation in *Segmentation.
// This is here to implement the Parameter interface, but it won't be used.
func (s *Segmentation) Value() *Segmentation {
	return s
}

// Code returns the Segmentation in ParameterNameCode.
func (s *Segmentation) Code() ParameterNameCode {
	return s.code
}

// String returns the Segmentation in string.
func (s *Segmentation) String() string {
	return fmt.Sprintf(
		"%s: {FirstSegment=%t, Class=%d, RemainingSegments=%d, LocalReference=%d}",
		s.code, s.FirstSegment, s.Class, s.RemainingSegments, s.LocalReference,
	)
}

// HopCounter represents the Hop Counter.
type HopCounter struct {
	paramType ParameterType
	code      ParameterNameCode
	length    int
	value     uint8
}

// NewHopCounter creates a new HopCounter.
func NewHopCounter(v uint8) *HopCounter {
	return &HopCounter{
		paramType: PTypeF,
		code:      PCodeHopCounter,
		length:    1,
		value:     v,
	}
}

// AsOptional makes the HopCounter an optional parameter.
func (h *HopCounter) AsOptional() *HopCounter {
	h.paramType = PTypeO
	h.length = 3
	return h
}

// Read sets the values retrieved from byte sequence in a HopCounter.
func (h *HopCounter) Read(b []byte) (int, error) {
	if h.paramType == PTypeO {
		return h.readOptional(b)
	}
	return h.read(b)
}

// read sets the values retrieved from byte sequence in a HopCounter.
func (h *HopCounter) read(b []byte) (int, error) {
	n := 1
	if len(b) < n {
		return 0, io.ErrUnexpectedEOF
	}

	h.code = PCodeHopCounter
	h.length = n
	h.value = b[0]

	return n, nil
}

func (h *HopCounter) readOptional(b []byte) (int, error) {
	n := 3
	if len(b) < n {
		return 0, io.ErrUnexpectedEOF
	}

	h.code = ParameterNameCode(b[0])
	if h.code != PCodeHopCounter {
		logf("invalid parameter code: expected %d, got %d", PCodeHopCounter, h.code)
	}
	h.length = int(b[1])
	if h.length != n {
		logf("%s: invalid length: expected %d, got %d", PCodeHopCounter, n, h.length)
	}
	h.value = b[2]

	return n, nil
}

// Write serializes the HopCounter parameter and returns it as a byte slice.
func (h *HopCounter) Write(b []byte) (int, error) {
	if h.paramType == PTypeO {
		return h.writeOptional(b)
	}
	return h.write(b)
}

// write serializes the HopCounter parameter and returns it as a byte slice.
func (h *HopCounter) write(b []byte) (int, error) {
	if len(b) < h.length {
		return 0, io.ErrUnexpectedEOF
	}

	b[0] = h.value
	return h.length, nil
}

func (h *HopCounter) writeOptional(b []byte) (int, error) {
	if len(b) < h.length {
		return 0, io.ErrUnexpectedEOF
	}

	b[0] = uint8(h.code)
	b[1] = uint8(h.length)
	b[2] = h.value

	return h.length, nil
}

// Value returns the HopCounter in uint8.
func (h *HopCounter) Value() uint8 {
	return h.value
}

// Code returns the HopCounter in ParameterNameCode.
func (h *HopCounter) Code() ParameterNameCode {
	return h.code
}

// String returns the HopCounter in string.
func (h *HopCounter) String() string {
	return fmt.Sprintf("%s: %d", h.code, h.value)
}

// Importance represents the Importance.
type Importance struct {
	paramType ParameterType
	code      ParameterNameCode
	length    int
	value     uint8
}

// NewImportance creates a new Importance.
func NewImportance(v uint8) *Importance {
	return &Importance{
		paramType: PTypeO,
		code:      PCodeImportance,
		length:    3,
		value:     v & 0b111,
	}
}

// AsOptional makes the Importance an optional parameter.
func (i *Importance) AsOptional() *Importance {
	i.paramType = PTypeO
	i.length = 3
	return i
}

// Read sets the values retrieved from byte sequence in a Importance.
func (i *Importance) Read(b []byte) (int, error) {
	if i.paramType != PTypeO {
		logf("Importance parameter must be optional: %v", i)
	}

	n := 3
	if len(b) < n {
		return 0, io.ErrUnexpectedEOF
	}

	i.code = ParameterNameCode(b[0])
	if i.code != PCodeImportance {
		logf("invalid parameter code: expected %d, got %d", PCodeImportance, i.code)
	}

	i.length = int(b[1])
	if i.length != n {
		logf("%s: invalid length: expected %d, got %d", PCodeImportance, n, i.length)
	}

	i.value = b[2] & 0b111

	return n, nil
}

// Write serializes the Importance parameter and returns it as a byte slice.
func (i *Importance) Write(b []byte) (int, error) {
	if i.paramType != PTypeO {
		logf("Importance parameter must be optional: %v", i)
	}

	if len(b) < i.length {
		return 0, io.ErrUnexpectedEOF
	}

	b[0] = uint8(i.code)
	b[1] = uint8(i.length)
	b[2] = i.value

	return i.length, nil
}

// Value returns the Importance in uint8.
func (i *Importance) Value() uint8 {
	return i.value
}

// Code returns the Importance in ParameterNameCode.
func (i *Importance) Code() ParameterNameCode {
	return i.code
}

// String returns the Importance in string.
func (i *Importance) String() string {
	return fmt.Sprintf("%s: %d", i.code, i.value)
}

// LongData represents the Long Data.
type LongData struct {
	paramType ParameterType
	code      ParameterNameCode
	length    int
	value     []byte
}

// NewLongData creates a new LongData.
func NewLongData(v []byte) *LongData {
	return &LongData{
		paramType: PTypeV,
		code:      PCodeLongData,
		length:    len(v),
		value:     v,
	}
}

// Read sets the values retrieved from byte sequence in a LongData.
func (l *LongData) Read(b []byte) (int, error) {
	n := len(b)
	l.code = PCodeLongData
	l.length = n
	l.value = b[:n]
	return n, nil
}

// Write serializes the LongData parameter and returns it as a byte slice.
func (l *LongData) Write(b []byte) (int, error) {
	if len(b) < l.length {
		return 0, io.ErrUnexpectedEOF
	}

	copy(b, l.value)
	return l.length, nil
}

// Value returns the LongData in []byte.
func (l *LongData) Value() []byte {
	return l.value
}

// Code returns the LongData in ParameterNameCode.
func (l *LongData) Code() ParameterNameCode {
	return l.code
}

// String returns the LongData in string.
func (l *LongData) String() string {
	return fmt.Sprintf("%s: %x", l.code, l.value)
}
