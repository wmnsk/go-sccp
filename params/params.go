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
	MarshalLen() int
	Code() ParameterNameCode
	fmt.Stringer
}

// ParameterType is a type for Parameter described in the tables in section 4 of Q.713.
type ParameterType uint8

// ParameterType values.
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
func ParseOptionalParameters(b []byte) ([]Parameter, int, error) {
	var params []Parameter
	var offset int
	for len(b) > 0 {
		p, n, err := ParseOptionalParameter(b[offset:])
		if err != nil {
			return nil, offset, err
		}
		params = append(params, p)
		if p.Code() == PCodeEndOfOptionalParameters {
			break
		}
		offset += n
	}
	return params, offset, nil
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

/*
Specific Parameter implementations

Each parameter should implement the following functions/methods:
- Constructor: NewParameterName creates a new ParameterName.
- Parser: ParseParameterName parses the given byte sequence as a ParameterName.
- Read: Read sets the values retrieved from byte sequence in a ParameterName.
- Write: Write serializes the ParameterName parameter and returns it as a byte slice.
- MarshalLen: MarshalLen returns the serial length of ParameterName.
- Code: Code returns the ParameterName in ParameterNameCode.
- Value: Value returns the ParameterName in specific type.
- String: String returns the ParameterName in string.
- Others (helpers): Optionally, each parameter can have helper functions to get specific values.

*/

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

// ParseEndOfOptionalParameters parses the given byte sequence as an EndOfOptionalParameters.
func ParseEndOfOptionalParameters(b []byte) (*EndOfOptionalParameters, int, error) {
	e := &EndOfOptionalParameters{}
	n, err := e.Read(b)
	if err != nil {
		return nil, n, err
	}

	return e, n, nil
}

// Read sets the values retrieved from byte sequence in a EndOfOptionalParameters.
func (e *EndOfOptionalParameters) Read(b []byte) (int, error) {
	n := 1
	if len(b) < n {
		return 0, io.ErrUnexpectedEOF
	}

	e.paramType = PTypeO
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

// MarshalLen returns the serial length of EndOfOptionalParameters.
func (e *EndOfOptionalParameters) MarshalLen() int {
	return e.length
}

// Code returns the EndOfOptionalParameters in ParameterNameCode.
func (e *EndOfOptionalParameters) Code() ParameterNameCode {
	return e.code
}

// Value returns the EndOfOptionalParameters in uint8.
func (e *EndOfOptionalParameters) Value() uint8 {
	return e.value
}

// String returns the EndOfOptionalParameters in string.
func (e *EndOfOptionalParameters) String() string {
	return fmt.Sprintf("{%s (%s): %d}", e.code, e.paramType, e.value)
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
func ParseDestinationLocalReference(b []byte) (*LocalReference, int, error) {
	return parseLocalReference(PCodeDestinationLocalReference, b)
}

// ParseSourceLocalReference parses the given byte sequence as a Source Local Reference.
func ParseSourceLocalReference(b []byte) (*LocalReference, int, error) {
	return parseLocalReference(PCodeSourceLocalReference, b)
}

func parseLocalReference(c ParameterNameCode, b []byte) (*LocalReference, int, error) {
	l := &LocalReference{
		code: c,
	}

	n, err := l.Read(b)
	if err != nil {
		return nil, n, err
	}

	return l, n, nil
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

// MarshalLen returns the serial length of LocalReference.
func (l *LocalReference) MarshalLen() int {
	return l.length
}

// Code returns the LocalReference in ParameterNameCode.
func (l *LocalReference) Code() ParameterNameCode {
	return l.code
}

// Value returns the LocalReference in []byte.
func (l *LocalReference) Value() []byte {
	return l.value
}

// String returns the LocalReference in string.
func (l *LocalReference) String() string {
	if l.code == PCodeDestinationLocalReference || l.code == PCodeSourceLocalReference {
		return fmt.Sprintf("{%s (%s): %d}", l.code, l.paramType, l.Uint32())
	}
	return fmt.Sprintf("{%s (%s): %d}", "(Destination or Source) local reference", l.paramType, l.Uint32())
}

// Uint32 returns the LocalReference in uint32.
func (l *LocalReference) Uint32() uint32 {
	return utils.Uint24To32(l.value)
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
// NewCalled/CallingPartyAddress as the first argument.
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
//
// When you are aware of the type of PartyAddress you are creating, you can use
// NewCalled/CallingPartyAddress to create a PartyAddress with the correct code.
// Otherwise, you can use AsCalled/Calling to set the code after creating a PartyAddress.
func NewPartyAddress(cdcg ParameterNameCode, ai uint8, spc uint16, ssn uint8, gt *GlobalTitle) *PartyAddress {
	if cdcg != PCodeCalledPartyAddress && cdcg != PCodeCallingPartyAddress {
		logf("invalid parameter code: expected %v or %v, got %v", PCodeCalledPartyAddress, PCodeCallingPartyAddress, cdcg)
	}

	p := &PartyAddress{
		paramType:   PTypeV,
		code:        cdcg,
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

// NewPartyAddressOptional creates a new PartyAddress from properly-typed values.
func NewPartyAddressOptional(cdcg ParameterNameCode, ai uint8, spc uint16, ssn uint8, gt *GlobalTitle) *PartyAddress {
	p := NewPartyAddress(cdcg, ai, spc, ssn, gt)
	p.paramType = PTypeO
	return p
}

// NewCalledPartyAddress creates a new PartyAddress for Called Party Address.
func NewCalledPartyAddress(ai uint8, spc uint16, ssn uint8, gt *GlobalTitle) *PartyAddress {
	return NewPartyAddress(PCodeCalledPartyAddress, ai, spc, ssn, gt)
}

// NewCallingPartyAddress creates a new PartyAddress for Calling Party Address.
func NewCallingPartyAddress(ai uint8, spc uint16, ssn uint8, gt *GlobalTitle) *PartyAddress {
	return NewPartyAddress(PCodeCallingPartyAddress, ai, spc, ssn, gt)
}

// NewCalledPartyAddressOptional creates a new PartyAddress for Called Party Address as an optional parameter.
func NewCalledPartyAddressOptional(ai uint8, spc uint16, ssn uint8, gt *GlobalTitle) *PartyAddress {
	return NewPartyAddressOptional(PCodeCalledPartyAddress, ai, spc, ssn, gt)
}

// NewCallingPartyAddressOptional creates a new PartyAddress for Calling Party Address as an optional parameter.
func NewCallingPartyAddressOptional(ai uint8, spc uint16, ssn uint8, gt *GlobalTitle) *PartyAddress {
	return NewPartyAddressOptional(PCodeCallingPartyAddress, ai, spc, ssn, gt)
}

// ParseCalledPartyAddress parses the given byte sequence as a mandatory fixed length
// Called Party Address and returns it as a PartyAddress.
func ParseCalledPartyAddress(b []byte) (*PartyAddress, int, error) {
	return parsePartyAddress(PTypeV, PCodeCalledPartyAddress, b)
}

// ParseCallingPartyAddress parses the given byte sequence as a mandatory fixed length
// Calling Party Address and returns it as a PartyAddress.
func ParseCallingPartyAddress(b []byte) (*PartyAddress, int, error) {
	return parsePartyAddress(PTypeV, PCodeCallingPartyAddress, b)
}

// ParseCalledPartyAddressOptional parses the given byte sequence as an optional
// Called Party Address and returns it as a PartyAddress.
func ParseCalledPartyAddressOptional(b []byte) (*PartyAddress, int, error) {
	return parsePartyAddress(PTypeO, PCodeCalledPartyAddress, b)
}

// ParseCallingPartyAddressOptional parses the given byte sequence as an optional
// Calling Party Address and returns it as a PartyAddress.
func ParseCallingPartyAddressOptional(b []byte) (*PartyAddress, int, error) {
	return parsePartyAddress(PTypeO, PCodeCallingPartyAddress, b)
}

func parsePartyAddress(ptype ParameterType, code ParameterNameCode, b []byte) (*PartyAddress, int, error) {
	p := &PartyAddress{
		paramType: ptype,
		code:      code,
	}

	n, err := p.Read(b)
	if err != nil {
		return nil, n, err
	}

	return p, n, nil
}

// Read sets the values retrieved from byte sequence in a PartyAddress.
func (p *PartyAddress) Read(b []byte) (int, error) {
	if p.paramType == PTypeO {
		return p.readOptional(b)
	}

	// force to read as V if it's not O
	p.paramType = PTypeV
	return p.read(b)
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

// Code returns the PartyAddress in ParameterNameCode.
func (p *PartyAddress) Code() ParameterNameCode {
	return p.code
}

// Value returns the PartyAddress as it is.
func (p *PartyAddress) Value() *PartyAddress {
	return p
}

// String returns the PartyAddress values in human readable format.
func (p *PartyAddress) String() string {
	return fmt.Sprintf("{%s (%s): {length: %d, Indicator: %#08b, SignalingPointCode: %d, SubsystemNumber: %d, GlobalTitle: %v}}",
		p.code, p.paramType, p.length, p.Indicator, p.SignalingPointCode, p.SubsystemNumber, p.GlobalTitle,
	)
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

// SetLength sets the length in length field.
// This should be called after changing the values in PartyAddress.
func (p *PartyAddress) SetLength() {
	p.length = p.MarshalLen() - 1
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

// ParseProtocolClass parses the given byte sequence as a ProtocolClass.
func ParseProtocolClass(b []byte) (*ProtocolClass, int, error) {
	p := &ProtocolClass{}
	n, err := p.Read(b)
	if err != nil {
		return nil, n, err
	}

	return p, n, nil
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

// MarshalLen returns the serial length of ProtocolClass.
func (p *ProtocolClass) MarshalLen() int {
	return p.length
}

// Code returns the ProtocolClass in ParameterNameCode.
func (p *ProtocolClass) Code() ParameterNameCode {
	return p.code
}

// Value returns the ProtocolClass in uint8.
func (p *ProtocolClass) Value() uint8 {
	return uint8(p.value)
}

// String returns the ProtocolClass in string.
func (p *ProtocolClass) String() string {
	return fmt.Sprintf(
		"{%s (%s): {Class: %d, ReturnOnError: %v}}",
		p.code, p.paramType, p.Class(), p.ReturnOnError(),
	)
}

// Class returns the class part from ProtocolClass parameter.
func (p *ProtocolClass) Class() int {
	return int(p.value) & 0xf
}

// ReturnOnError judges if ProtocolClass has "Return Message On Error" option.
func (p *ProtocolClass) ReturnOnError() bool {
	return (int(p.value) >> 7) == 1
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

// ParseSegmentingReassembling parses the given byte sequence as a SegmentingReassembling.
func ParseSegmentingReassembling(b []byte) (*SegmentingReassembling, int, error) {
	s := &SegmentingReassembling{}
	n, err := s.Read(b)
	if err != nil {
		return nil, n, err
	}

	return s, n, nil
}

// Read sets the values retrieved from byte sequence in a SegmentingReassembling.
func (s *SegmentingReassembling) Read(b []byte) (int, error) {
	n := 1
	if len(b) < n {
		return 0, io.ErrUnexpectedEOF
	}

	s.paramType = PTypeF
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

// MarshalLen returns the serial length of SegmentingReassembling.
func (s *SegmentingReassembling) MarshalLen() int {
	return s.length
}

// Code returns the SegmentingReassembling in ParameterNameCode.
func (s *SegmentingReassembling) Code() ParameterNameCode {
	return s.code
}

// Value returns the SegmentingReassembling in uint8.
func (s *SegmentingReassembling) Value() uint8 {
	return s.value
}

// String returns the SegmentingReassembling in string.
func (s *SegmentingReassembling) String() string {
	return fmt.Sprintf("{%s (%s): %d}", s.code, s.paramType, s.value)
}

// MoreData judges if the message has more data.
func (s *SegmentingReassembling) MoreData() bool {
	return s.value&0b1 == 1
}

// ReceiveSequenceNumber represents the Receive Sequence Number.
type ReceiveSequenceNumber struct {
	paramType ParameterType
	code      ParameterNameCode
	length    int
	value     uint8
}

// NewReceiveSequenceNumber creates a new ReceiveSequenceNumber.
// The value is masked out to 0b11111110 since the LSB is spare.
func NewReceiveSequenceNumber(v uint8) *ReceiveSequenceNumber {
	return &ReceiveSequenceNumber{
		paramType: PTypeF,
		code:      PCodeReceiveSequenceNumber,
		length:    1,
		value:     v & 0b11111110,
	}
}

// ParseReceiveSequenceNumber parses the given byte sequence as a ReceiveSequenceNumber.
func ParseReceiveSequenceNumber(b []byte) (*ReceiveSequenceNumber, int, error) {
	r := &ReceiveSequenceNumber{}
	n, err := r.Read(b)
	if err != nil {
		return nil, n, err
	}

	return r, n, nil
}

// Read sets the values retrieved from byte sequence in a ReceiveSequenceNumber.
func (r *ReceiveSequenceNumber) Read(b []byte) (int, error) {
	n := 1
	if len(b) < n {
		return 0, io.ErrUnexpectedEOF
	}

	r.paramType = PTypeF
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

	b[0] = r.value & 0b11111110
	return r.length, nil
}

// MarshalLen returns the serial length of ReceiveSequenceNumber.
func (r *ReceiveSequenceNumber) MarshalLen() int {
	return r.length
}

// Code returns the ReceiveSequenceNumber in ParameterNameCode.
func (r *ReceiveSequenceNumber) Code() ParameterNameCode {
	return r.code
}

// Value returns the ReceiveSequenceNumber in uint8.
func (r *ReceiveSequenceNumber) Value() uint8 {
	return r.value
}

// String returns the ReceiveSequenceNumber in string.
func (r *ReceiveSequenceNumber) String() string {
	return fmt.Sprintf("{%s (%s): %d}", r.code, r.paramType, r.value)
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

// ParseSequencingSegmenting parses the given byte sequence as a SequencingSegmenting.
func ParseSequencingSegmenting(b []byte) (*SequencingSegmenting, int, error) {
	s := &SequencingSegmenting{}
	n, err := s.Read(b)
	if err != nil {
		return nil, n, err
	}

	return s, n, nil
}

// Read sets the values retrieved from byte sequence in a SequencingSegmenting.
func (s *SequencingSegmenting) Read(b []byte) (int, error) {
	n := 2
	if len(b) < n {
		return 0, io.ErrUnexpectedEOF
	}

	s.paramType = PTypeF
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

// MarshalLen returns the serial length of SequencingSegmenting.
func (s *SequencingSegmenting) MarshalLen() int {
	return s.length
}

// Code returns the SequencingSegmenting in ParameterNameCode.
func (s *SequencingSegmenting) Code() ParameterNameCode {
	return s.code
}

// Value returns the SequencingSegmenting as it is.
func (s *SequencingSegmenting) Value() *SequencingSegmenting {
	return s
}

// String returns the SequencingSegmenting in string.
func (s *SequencingSegmenting) String() string {
	return fmt.Sprintf(
		"{%s: {SendSequenceNumber=%d, ReceiveSequenceNumber=%d, MoreData=%t}}",
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

// NewCreditOptional creates a new optional Credit.
func NewCreditOptional(v uint8) *Credit {
	c := NewCredit(v)
	c.paramType = PTypeO
	return c
}

// ParseCredit parses the given byte sequence as a Credit.
func ParseCredit(b []byte) (*Credit, int, error) {
	c := &Credit{}
	n, err := c.Read(b)
	if err != nil {
		return nil, n, err
	}

	return c, n, nil
}

// ParseCreditOptional parses the given byte sequence as an optional Credit.
func ParseCreditOptional(b []byte) (*Credit, int, error) {
	c := &Credit{paramType: PTypeO}
	n, err := c.Read(b)
	if err != nil {
		return nil, n, err
	}

	return c, n, nil
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

	c.paramType = PTypeF
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
	if c.length != n-2 {
		logf("%s: invalid length: expected %d, got %d", PCodeCredit, n-2, c.length)
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

// MarshalLen returns the serial length of Credit.
func (c *Credit) MarshalLen() int {
	if c.paramType == PTypeO {
		return 2 + c.length
	}
	return c.length
}

// Code returns the Credit in ParameterNameCode.
func (c *Credit) Code() ParameterNameCode {
	return c.code
}

// Value returns the Credit in uint8.
func (c *Credit) Value() uint8 {
	return c.value
}

// String returns the Credit in string.
func (c *Credit) String() string {
	return fmt.Sprintf("{%s (%s): %d}", c.code, c.paramType, c.value)
}

// Cause represents a common structure for all Cause types.
type Cause[T ~uint8] struct {
	paramType ParameterType
	code      ParameterNameCode
	length    int
	value     T
}

// NewCause creates a new Cause with the given value and its type.
func NewCause[T ~uint8](value T) *Cause[T] {
	c := &Cause[T]{
		paramType: PTypeF,
		length:    1,
		value:     value,
	}

	switch any(c).(type) {
	case *ReleaseCause:
		c.code = PCodeReleaseCause
	case *ReturnCause:
		c.code = PCodeReturnCause
	case *ResetCause:
		c.code = PCodeResetCause
	case *ErrorCause:
		c.code = PCodeErrorCause
	case *RefusalCause:
		c.code = PCodeRefusalCause
	default:
		logf("invalid Cause type: %T", c)
	}

	return c
}

// Read sets the values retrieved from byte sequence in a Cause.
func (c *Cause[T]) Read(b []byte) (int, error) {
	n := 1
	if len(b) < n {
		return 0, io.ErrUnexpectedEOF
	}

	c.paramType = PTypeF

	switch any(c).(type) {
	case *ReleaseCause:
		c.code = PCodeReleaseCause
	case *ReturnCause:
		c.code = PCodeReturnCause
	case *ResetCause:
		c.code = PCodeResetCause
	case *ErrorCause:
		c.code = PCodeErrorCause
	case *RefusalCause:
		c.code = PCodeRefusalCause
	default:
		return 0, UnsupportedParameterError(b[0])
	}

	c.length = n
	c.value = T(b[0])

	return n, nil
}

// Write serializes the Cause parameter and returns it as a byte slice.
func (c *Cause[T]) Write(b []byte) (int, error) {
	if len(b) < c.length {
		return 0, io.ErrUnexpectedEOF
	}

	b[0] = uint8(c.value)
	return c.length, nil
}

// MarshalLen returns the serial length of Cause.
func (c *Cause[T]) MarshalLen() int {
	return c.length
}

// Code returns the code in the Cause.
func (c *Cause[T]) Code() ParameterNameCode {
	return c.code
}

// Value returns the value in the Cause.
func (c *Cause[T]) Value() T {
	return T(c.value)
}

// String returns the Cause as a string.
func (c *Cause[T]) String() string {
	return fmt.Sprintf("{%s (%s): %v}", c.code, c.paramType, c.value)
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

// ReleaseCause is a specific Cause for ReleaseCause.
type ReleaseCause = Cause[ReleaseCauseValue]

// ParseReleaseCause parses the given byte sequence as a ReleaseCause.
func ParseReleaseCause(b []byte) (*ReleaseCause, int, error) {
	c := &ReleaseCause{}
	n, err := c.Read(b)
	if err != nil {
		return nil, n, err
	}

	return c, n, nil
}

// ReturnCause is a specific Cause for ReturnCause.
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

// ReturnCause is a specific instance of Cause.
type ReturnCause = Cause[ReturnCauseValue]

// ParseReturnCause parses the given byte sequence as a ReturnCause.
func ParseReturnCause(b []byte) (*ReturnCause, int, error) {
	c := &ReturnCause{}
	n, err := c.Read(b)
	if err != nil {
		return nil, n, err
	}

	return c, n, nil
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

// ResetCause is a specific Cause for ResetCause.
type ResetCause = Cause[ResetCauseValue]

// ParseResetCause parses the given byte sequence as a ResetCause.
func ParseResetCause(b []byte) (*ResetCause, int, error) {
	c := &ResetCause{}
	n, err := c.Read(b)
	if err != nil {
		return nil, n, err
	}

	return c, n, nil
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

// ErrorCause is a specific Cause for ErrorCause.
type ErrorCause = Cause[ErrorCauseValue]

// ParseErrorCause parses the given byte sequence as a ErrorCause.
func ParseErrorCause(b []byte) (*ErrorCause, int, error) {
	c := &ErrorCause{}
	n, err := c.Read(b)
	if err != nil {
		return nil, n, err
	}

	return c, n, nil
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

// RefusalCause is a specific Cause for RefusalCause.
type RefusalCause = Cause[RefusalCauseValue]

// ParseRefusalCause parses the given byte sequence as a RefusalCause.
func ParseRefusalCause(b []byte) (*RefusalCause, int, error) {
	c := &RefusalCause{}
	n, err := c.Read(b)
	if err != nil {
		return nil, n, err
	}

	return c, n, nil
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

// NewDataOptional creates a new Data as an optional parameter.
func NewDataOptional(v []byte) *Data {
	d := NewData(v)
	d.paramType = PTypeO
	return d
}

// ParseData parses the given byte sequence as a Data.
func ParseData(b []byte) (*Data, int, error) {
	d := &Data{}
	n, err := d.Read(b)
	if err != nil {
		return nil, n, err
	}

	return d, n, nil
}

// ParseDataOptional parses the given byte sequence as an optional Data.
func ParseDataOptional(b []byte) (*Data, int, error) {
	d := &Data{paramType: PTypeO}
	n, err := d.Read(b)
	if err != nil {
		return nil, n, err
	}

	return d, n, nil
}

// Read sets the values retrieved from byte sequence in a Data.
func (d *Data) Read(b []byte) (int, error) {
	if d.paramType == PTypeO {
		return d.readOptional(b)
	}

	// force to read as V if it's not O
	if d.paramType == PTypeF {
		d.paramType = PTypeV
	}
	return d.read(b)
}

// read sets the values retrieved from byte sequence in a Data.
func (d *Data) read(b []byte) (int, error) {
	n := len(b)
	if n < 1 {
		return 0, io.ErrUnexpectedEOF
	}

	d.code = PCodeData
	d.length = int(b[0])
	if d.length == 0 {
		d.value = nil
		return 1, nil
	}

	if n < d.length+1 {
		return 1, io.ErrUnexpectedEOF
	}

	d.value = b[1 : d.length+1]

	return n, nil
}

func (d *Data) readOptional(b []byte) (int, error) {
	n := len(b)

	d.code = ParameterNameCode(b[0])
	if d.code != PCodeData {
		logf("invalid parameter code: expected %d, got %d", PCodeData, d.code)
	}

	m, err := d.read(b[1:])
	n += m
	if err != nil {
		return n, err
	}

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
	if len(b) < d.length+1 {
		return 0, io.ErrUnexpectedEOF
	}

	b[0] = uint8(d.length)
	if d.length == 0 {
		return 1, nil
	}

	copy(b[1:d.length+1], d.value)
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

// MarshalLen returns the serial length of Data.
func (d *Data) MarshalLen() int {
	if d.paramType == PTypeO {
		return 2 + len(d.value)
	}
	return 1 + len(d.value)
}

// Code returns the Data in ParameterNameCode.
func (d *Data) Code() ParameterNameCode {
	return d.code
}

// Value returns the Data in []byte.
func (d *Data) Value() []byte {
	return d.value
}

// String returns the Data in string.
func (d *Data) String() string {
	return fmt.Sprintf("{%s (%s): %x}", d.code, d.paramType, d.value)
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
		length:            4,
		FirstSegment:      first,
		Class:             cls & 0b1,
		RemainingSegments: rem & 0b1111,
		LocalReference:    lrn & 0x00ffffff,
	}
}

// NewSegmentationOptional creates a new optional Segmentation.
func NewSegmentationOptional(first bool, cls, rem uint8, lrn uint32) *Segmentation {
	return NewSegmentation(first, cls, rem, lrn)
}

// ParseSegmentation parses the given byte sequence as a Segmentation.
func ParseSegmentation(b []byte) (*Segmentation, int, error) {
	s := &Segmentation{}
	n, err := s.Read(b)
	if err != nil {
		return nil, n, err
	}

	return s, n, nil
}

// ParseSegmentationOptional parses the given byte sequence as an optional Segmentation.
func ParseSegmentationOptional(b []byte) (*Segmentation, int, error) {
	return ParseSegmentation(b)
}

// Read sets the values retrieved from byte sequence in a Segmentation.
func (s *Segmentation) Read(b []byte) (int, error) {
	if s.paramType != PTypeO {
		s.paramType = PTypeO
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
	if s.length != n-2 {
		logf("%s: invalid length: expected %d, got %d", PCodeSegmentation, n-2, s.length)
	}

	s.FirstSegment = b[2]>>7&0b1 == 1
	s.Class = b[2] >> 6 & 0b1
	s.RemainingSegments = b[2] & 0b1111
	s.LocalReference = utils.Uint24To32(b[3:6])

	return n, nil
}

// Write serializes the Segmentation parameter and returns it as a byte slice.
func (s *Segmentation) Write(b []byte) (int, error) {
	if s.paramType != PTypeO {
		logf("Segmentation parameter must be optional: %v", s)
	}

	n := s.length + 2
	if len(b) < n {
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

	return n, nil
}

// MarshalLen returns the serial length of Segmentation.
func (s *Segmentation) MarshalLen() int {
	return s.length + 2
}

// Code returns the Segmentation in ParameterNameCode.
func (s *Segmentation) Code() ParameterNameCode {
	return s.code
}

// Value returns the Segmentation as it is.
func (s *Segmentation) Value() *Segmentation {
	return s
}

// String returns the Segmentation in string.
func (s *Segmentation) String() string {
	return fmt.Sprintf(
		"{%s (%s): {FirstSegment=%t, Class=%d, RemainingSegments=%d, LocalReference=%d}}",
		s.code, s.paramType, s.FirstSegment, s.Class, s.RemainingSegments, s.LocalReference,
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

// NewHopCounterOptional creates a new optional HopCounter.
func NewHopCounterOptional(v uint8) *HopCounter {
	h := NewHopCounter(v)
	h.paramType = PTypeO
	return h
}

// ParseHopCounter parses the given byte sequence as a HopCounter.
func ParseHopCounter(b []byte) (*HopCounter, int, error) {
	h := &HopCounter{}
	n, err := h.Read(b)
	if err != nil {
		return nil, n, err
	}

	return h, n, nil
}

// ParseHopCounterOptional parses the given byte sequence as an optional HopCounter.
func ParseHopCounterOptional(b []byte) (*HopCounter, int, error) {
	h := &HopCounter{paramType: PTypeO}
	n, err := h.Read(b)
	if err != nil {
		return nil, n, err
	}

	return h, n, nil
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

	h.paramType = PTypeF
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
	if h.length != n-2 {
		logf("%s: invalid length: expected %d, got %d", PCodeHopCounter, n-2, h.length)
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

// MarshalLen returns the serial length of HopCounter.
func (h *HopCounter) MarshalLen() int {
	if h.paramType == PTypeO {
		return 2 + h.length
	}
	return h.length
}

// Code returns the HopCounter in ParameterNameCode.
func (h *HopCounter) Code() ParameterNameCode {
	return h.code
}

// Value returns the HopCounter in uint8.
func (h *HopCounter) Value() uint8 {
	return h.value
}

// String returns the HopCounter in string.
func (h *HopCounter) String() string {
	return fmt.Sprintf("{%s (%s): %d}", h.code, h.paramType, h.value)
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
		length:    1,
		value:     v & 0b111,
	}
}

// NewImportanceOptional creates a new optional Importance.
func NewImportanceOptional(v uint8) *Importance {
	return NewImportance(v)
}

// ParseImportance parses the given byte sequence as an Importance.
func ParseImportance(b []byte) (*Importance, int, error) {
	i := &Importance{}
	n, err := i.Read(b)
	if err != nil {
		return nil, n, err
	}

	return i, n, nil
}

// ParseImportanceOptional parses the given byte sequence as an optional Importance.
func ParseImportanceOptional(b []byte) (*Importance, int, error) {
	return ParseImportance(b)
}

// Read sets the values retrieved from byte sequence in a Importance.
func (i *Importance) Read(b []byte) (int, error) {
	if i.paramType != PTypeO {
		i.paramType = PTypeO
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
	if i.length != n-2 {
		logf("%s: invalid length: expected %d, got %d", PCodeImportance, n-2, i.length)
	}

	i.value = b[2] & 0b111

	return n, nil
}

// Write serializes the Importance parameter and returns it as a byte slice.
func (i *Importance) Write(b []byte) (int, error) {
	if i.paramType != PTypeO {
		logf("Importance parameter must be optional: %v", i)
	}

	n := i.length + 2
	if len(b) < n {
		return 0, io.ErrUnexpectedEOF
	}

	b[0] = uint8(i.code)
	b[1] = uint8(i.length)
	b[2] = i.value

	return n, nil
}

// MarshalLen returns the serial length of Importance.
func (i *Importance) MarshalLen() int {
	return i.length + 2
}

// Code returns the Importance in ParameterNameCode.
func (i *Importance) Code() ParameterNameCode {
	return i.code
}

// Value returns the Importance in uint8.
func (i *Importance) Value() uint8 {
	return i.value
}

// String returns the Importance in string.
func (i *Importance) String() string {
	return fmt.Sprintf("{%s (%s): %d}", i.code, i.paramType, i.value)
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

// ParseLongData parses the given byte sequence as a LongData.
func ParseLongData(b []byte) (*LongData, int, error) {
	l := &LongData{}
	n, err := l.Read(b)
	if err != nil {
		return nil, n, err
	}

	return l, n, nil
}

// Read sets the values retrieved from byte sequence in a LongData.
func (l *LongData) Read(b []byte) (int, error) {
	n := len(b)

	l.paramType = PTypeV
	l.code = PCodeLongData

	l.length = int(binary.BigEndian.Uint16(b[:2]))
	if n < l.length+2 {
		return n, io.ErrUnexpectedEOF
	}

	l.value = b[2 : l.length+2]
	return n, nil
}

// Write serializes the LongData parameter and returns it as a byte slice.
func (l *LongData) Write(b []byte) (int, error) {
	if len(b) < l.length+2 {
		return 0, io.ErrUnexpectedEOF
	}

	binary.BigEndian.PutUint16(b, uint16(l.length))
	copy(b[2:], l.value)
	return l.length, nil
}

// MarshalLen returns the serial length of LongData.
func (l *LongData) MarshalLen() int {
	return l.length + 2
}

// Code returns the LongData in ParameterNameCode.
func (l *LongData) Code() ParameterNameCode {
	return l.code
}

// Value returns the LongData in []byte.
func (l *LongData) Value() []byte {
	return l.value
}

// String returns the LongData in string.
func (l *LongData) String() string {
	return fmt.Sprintf("{%s (%s): %x}", l.code, l.paramType, l.value)
}
