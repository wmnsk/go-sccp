package sccp

import (
	"encoding/hex"
	"fmt"
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

	Opts []*params.Optional // all others

	// just pointers, not used for Marshal-ing,  I kust really need these two
	// similar objects are expected to be found in Opts
	Data                *params.Optional
	CallingPartyAddress *params.PartyAddress

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

// MarshalBinary returns the byte sequence generated from a UDT instance.
func (msg *CR) MarshalBinary() ([]byte, error) {
	b := make([]byte, msg.MarshalLen())
	if err := msg.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

func (msg *CR) MarshalLen() int {
	l := 5 + 2 + 1 // fixed + ptrs + last optional
	for _, v := range msg.Opts {
		l += int(v.Len) + 2
	}
	l += int(msg.CalledPartyAddress.Length) + 1

	return l
}

func (msg *CR) MarshalTo(b []byte) error {
	b[0] = uint8(msg.Type)
	msg.SourceLocalReference.Read(b[1:4])
	b[4] = byte(msg.ProtocolClass)
	b[5] = 2
	b[6] = msg.CalledPartyAddress.Length + 2
	if err := msg.CalledPartyAddress.MarshalTo(b[7 : 7+int(msg.CalledPartyAddress.Length)+1]); err != nil {
		return err
	}
	p := 6 + msg.CalledPartyAddress.Length + 1 + 1
	for i := 0; i < len(msg.Opts); i++ {
		b[p] = msg.Opts[i].Tag
		b[p+1] = msg.Opts[i].Len
		copy(b[p+2:], msg.Opts[i].Value)

		p += msg.Opts[i].Len + 2
	}
	return nil
}

func (msg *CR) String() string {
	s := fmt.Sprintf("{Type: CR, CalledPartyAddress: %v", msg.CalledPartyAddress)
	if msg.CallingPartyAddress != nil {
		s += fmt.Sprintf(", CallingPartyAddress: %v", msg.CallingPartyAddress)
	}
	if msg.Data != nil {
		s += fmt.Sprintf(", DataLength: %d, Data: %s", msg.Data.Len, hex.EncodeToString(msg.Data.Value))
	}

	return s + "}"
}

// MessageType returns the Message Type in int.
func (msg *CR) MessageType() MsgType {
	return msg.Type
}

func (msg *CR) MessageTypeName() string {
	return "CR"
}
