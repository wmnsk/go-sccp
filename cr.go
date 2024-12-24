package sccp

import (
	"io"

	"github.com/wmnsk/go-sccp/params"
)

/*
Message type code 2.1 F 1
Source local reference 3.3 F 3
Protocol class 3.6 F 1
Called party address 3.4 V 3 minimum
Credit 3.10 O 3
Calling party address 3.5 O 4 minimum
Data 3.16 O 3-130
Hop counter 3.18 O 3
Importance 3.19 O 3
End of optional parameters 3.1 O 1
*/

type CR struct {
	Type                 MsgType
	SourceLocalReference params.LocalReference
	params.ProtocolClass
	CalledPartyAddress *params.PartyAddress

	Data                *params.Optional // because I do really need it
	CallingPartyAddress *params.PartyAddress

	Opts []*params.Optional // all others

	mptr uint8
	optr uint8
}

func ParseCR(b []byte) (*CR, error) {
	msg := &CR{}
	if err := msg.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return msg, nil
}

func (msg *CR) UnmarshalBinary(b []byte) error {
	l := uint8(len(b))
	if l <= (1 + 3 + 1 + 2 /*ptrs*/ + 3) { // where CdPA starts
		return io.ErrUnexpectedEOF
	}

	msg.Type = MsgType(b[0])
	msg.SourceLocalReference.Read(b[1:4])
	msg.ProtocolClass = params.ProtocolClass(b[4])

	msg.mptr = b[5]
	if l < (5 + msg.mptr + 2) {
		return io.ErrUnexpectedEOF
	}
	msg.optr = b[6]
	if l < (6 + msg.optr + 1) {
		return io.ErrUnexpectedEOF
	}

	var err error
	if msg.CalledPartyAddress, err = params.ParsePartyAddress(b[5+msg.mptr : 6+msg.optr]); err != nil {
		return err
	}
	return msg.parseOptional(b[6+msg.optr:])
}

func (msg *CR) parseOptional(b []byte) error {
	// fmt.Println(hex.EncodeToString(b))
	p := uint8(0)
	for p < uint8(len(b)) {
		t := b[p]

		if t == 0 {
			return nil
		}
		if (p + 1) >= uint8(len(b)) {
			return io.ErrUnexpectedEOF
		}

		l := b[p+1]
		if (p + 1 + l) >= uint8(len(b)) {
			return io.ErrUnexpectedEOF
		}

		o := &params.Optional{
			Tag:   t,
			Len:   l,
			Value: b[p+2 : p+2+l],
		}

		switch t {
		case params.DataTag:
			msg.Data = o
		case params.CgPtyAddrTag:
			var err error
			msg.CallingPartyAddress, err = params.ParsePartyAddress(b[p : p+2+l])
			if err != nil {
				return err
			}
		}

		msg.Opts = append(msg.Opts, o)
		p += 2 + l

	}

	return nil
}
