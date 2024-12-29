// Copyright 2019-2024 go-sccp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
package params

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
