package params

import (
	"fmt"
	"io"

	"github.com/wmnsk/go-sccp/utils"
)

// GlobalTitle is a GlobalTitle inside the Called/Calling Party Address.
type GlobalTitle struct {
	// GTI is included in the Address Indicator which is not a part of
	// Global Title itself, but necessary to encode/decode it properly.
	GTI GlobalTitleIndicator
	TranslationType
	NumberingPlan
	EncodingScheme
	NatureOfAddressIndicator
	AddressInformation []byte
}

// GlobalTitleIndicator is a type of Global Title Indicator.
// See Q.713 3.4.1 for more details.
type GlobalTitleIndicator uint8

const (
	GTINAIOnly   GlobalTitleIndicator = 0b0001
	GTITTOnly    GlobalTitleIndicator = 0b0010
	GTITTNPES    GlobalTitleIndicator = 0b0011
	GTITTNPESNAI GlobalTitleIndicator = 0b0100
)

// NatureOfAddressIndicator is a type of Nature of Address Indicator.
type NatureOfAddressIndicator uint8

// NatureOfAddressIndicator values.
const (
	NAIUnknown                   NatureOfAddressIndicator = 0b00000000
	NAISubscriberNumber          NatureOfAddressIndicator = 0b00000001
	_                            NatureOfAddressIndicator = 0b00000010 // Reserved for national use
	NAINationalSignificantNumber NatureOfAddressIndicator = 0b00000011
	NAIInternationalNumber       NatureOfAddressIndicator = 0b00000100
)

// Even returns the NatureOfAddressIndicator with the last bit set to 0.
func (nai NatureOfAddressIndicator) Even() NatureOfAddressIndicator {
	return nai & 0b01111111
}

// Odd returns the NatureOfAddressIndicator with the last bit set to 1.
func (nai NatureOfAddressIndicator) Odd() NatureOfAddressIndicator {
	return nai | 0b10000000
}

// TranslationType is a type of Translation Type.
type TranslationType uint8

// NumberingPlan is a type of Numbering Plan.
type NumberingPlan uint8

// NumberingPlan values.
const (
	NPUnknown        NumberingPlan = 0b0000
	NPISDNTelephony  NumberingPlan = 0b0001
	NPGeneric        NumberingPlan = 0b0010
	NPData           NumberingPlan = 0b0011
	NPTelex          NumberingPlan = 0b0100
	NPMaritimeMobile NumberingPlan = 0b0101
	NPLandMobile     NumberingPlan = 0b0110
	NPISDNMobile     NumberingPlan = 0b0111
	NPPrivate        NumberingPlan = 0b1110
)

// EncodingScheme is a type of Encoding Scheme.
type EncodingScheme uint8

// EncodingScheme values.
const (
	ESUnknown          EncodingScheme = 0b0000
	ESBCDOdd           EncodingScheme = 0b0001
	ESBCDEven          EncodingScheme = 0b0010
	ESNationalSpecific EncodingScheme = 0b0011
)

// NewGlobalTitle creates a new GlobalTitle.
//
// The first argument is a Global Title Indicator, which is included in the Address Indicator
// in the parent PartyAddress.
func NewGlobalTitle(
	gti GlobalTitleIndicator,
	tt TranslationType,
	np NumberingPlan,
	es EncodingScheme,
	nai NatureOfAddressIndicator,
	addr []byte,
) *GlobalTitle {
	gt := &GlobalTitle{GTI: gti}

	switch gti {
	case GTINAIOnly:
		gt.NatureOfAddressIndicator = nai
	case GTITTOnly:
		gt.TranslationType = tt
	case GTITTNPES:
		gt.TranslationType = tt
		gt.NumberingPlan = np
		gt.EncodingScheme = es
	case GTITTNPESNAI:
		gt.TranslationType = tt
		gt.NumberingPlan = np
		gt.EncodingScheme = es
		gt.NatureOfAddressIndicator = nai
	}

	gt.AddressInformation = addr
	return gt
}

// MarshalBinary returns the byte sequence generated from a GlobalTitle.
func (g *GlobalTitle) MarshalBinary() []byte {
	b := make([]byte, g.MarshalLen())
	if err := g.MarshalTo(b); err != nil {
		panic(err)
	}

	return b
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (g *GlobalTitle) MarshalTo(b []byte) error {
	if len(b) < g.MarshalLen() {
		return io.ErrUnexpectedEOF
	}
	offset := 0
	switch g.GTI {
	case GTINAIOnly:
		b[offset] = uint8(g.NatureOfAddressIndicator)
		offset++
	case GTITTOnly:
		b[offset] = uint8(g.TranslationType)
		offset++
	case GTITTNPES:
		b[offset] = uint8(g.TranslationType)
		b[offset+1] = uint8(g.NumberingPlan)<<4 | uint8(g.EncodingScheme)
		offset += 2
	case GTITTNPESNAI:
		b[offset] = uint8(g.TranslationType)
		b[offset+1] = uint8(g.NumberingPlan)<<4 | uint8(g.EncodingScheme)
		b[offset+2] = uint8(g.NatureOfAddressIndicator)
		offset += 3
	}

	copy(b[offset:g.MarshalLen()], g.AddressInformation)
	return nil
}

// ParseGlobalTitle decodes given byte sequence as a GlobalTitle.
// The given byte sequence should not include the excess bytes for the parent PartyAddress.
// otherwise, AddressInformation will include them.
func ParseGlobalTitle(gti GlobalTitleIndicator, b []byte) (*GlobalTitle, error) {
	g := &GlobalTitle{GTI: gti}
	if err := g.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return g, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a GlobalTitle.
// The given byte sequence should not include the excess bytes for the parent PartyAddress.
// otherwise, AddressInformation will include them.
func (g *GlobalTitle) UnmarshalBinary(b []byte) error {
	if len(b) < g.lenByGTI() {
		return io.ErrUnexpectedEOF
	}

	offset := 0
	switch g.GTI {
	case GTINAIOnly:
		g.NatureOfAddressIndicator = NatureOfAddressIndicator(b[offset])
		offset++
	case GTITTOnly:
		g.TranslationType = TranslationType(b[offset])
		offset++
	case GTITTNPES:
		g.TranslationType = TranslationType(b[offset])
		g.NumberingPlan = NumberingPlan(b[offset+1] >> 4)
		g.EncodingScheme = EncodingScheme(b[offset+1] & 0x0F)
		offset += 2
	case GTITTNPESNAI:
		g.TranslationType = TranslationType(b[offset])
		g.NumberingPlan = NumberingPlan(b[offset+1] >> 4)
		g.EncodingScheme = EncodingScheme(b[offset+1] & 0x0F)
		g.NatureOfAddressIndicator = NatureOfAddressIndicator(b[offset+2])
		offset += 3
	}

	g.AddressInformation = b[offset:]
	return nil
}

// MarshalLen returns the serial length of a GlobalTitle.
func (g *GlobalTitle) MarshalLen() int {
	return g.lenByGTI()
}

func (g *GlobalTitle) lenByGTI() int {
	var l int
	switch g.GTI {
	case GTINAIOnly:
		l += 1
	case GTITTOnly:
		l += 1
	case GTITTNPES:
		l += 2
	case GTITTNPESNAI:
		l += 3
	}

	if g.AddressInformation != nil {
		l += len(g.AddressInformation)
	}

	return l
}

// IsOddDigits reports whether AddressInformation is odd number or not.
func (g *GlobalTitle) IsOddDigits() bool {
	return g.EncodingScheme == ESBCDOdd
}

// String returns the GlobalTitle in a human-readable format.
func (g *GlobalTitle) String() string {
	return fmt.Sprintf("{GTI: %d, TransationType: %d, NumberingPlan: %d, EncodingScheme: %d, NatureOfAddressIndicator: %d, AddressInformation: %#x}",
		g.GTI, g.TranslationType, g.NumberingPlan, g.EncodingScheme, g.NatureOfAddressIndicator, g.AddressInformation,
	)
}

// Address returns the AddressInformation in a human-friendly string.
func (g *GlobalTitle) Address() string {
	return utils.BCDDecode(g.IsOddDigits(), g.AddressInformation)
}
