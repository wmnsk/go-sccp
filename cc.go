package sccp

import (
	"fmt"
	"io"

	"github.com/wmnsk/go-sccp/params"
	"github.com/wmnsk/go-sccp/utils"
)

type CC struct {
	Type                      MsgType
	DestinationLocalReference uint32
	SourceLocalReference      uint32
	params.ProtocolClass

	Opts []*params.Optional

	Data               *params.Optional
	CalledPartyAddress *params.PartyAddress
}

// ParseCC decodes given byte sequence as a SCCP CC.
func ParseCC(b []byte) (*CC, error) {
	msg := &CC{}
	if err := msg.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return msg, nil
}

func (msg *CC) UnmarshalBinary(b []byte) error {
	l := uint8(len(b))

	if l < (1 + 3 + 3 + 1 + 1) {
		return io.ErrUnexpectedEOF
	}

	msg.Type = MsgType(b[0])
	msg.DestinationLocalReference = utils.Uint24To32(b[1:4])
	msg.SourceLocalReference = utils.Uint24To32(b[4:7])

	msg.ProtocolClass = params.ProtocolClass(b[7])

	optr := b[8]

	if optr == 0 {
		return nil
	}
	if optr != 1 {
		return io.ErrUnexpectedEOF
	}

	if err := msg.parseOptional(b[9:]); err != nil {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func (msg *CC) parseOptional(b []byte) error {
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
		case params.CdPtyAddrTag:
			var err error
			msg.CalledPartyAddress, err = params.ParsePartyAddress(b[p : p+2+l])
			if err != nil {
				return err
			}
		}

		msg.Opts = append(msg.Opts, o)
		p += 2 + l
	}

	return nil
}

// MarshalBinary returns the byte sequence generated from a CC instance.
func (msg *CC) MarshalBinary() ([]byte, error) {
	b := make([]byte, msg.MarshalLen())
	if err := msg.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

func (msg *CC) MarshalLen() int {
	if len(msg.Opts) == 0 {
		return 9 // 8 fixed + 0 ptr
	}
	l := 10 // 8 fixed + 1 ptr + last optional
	for _, v := range msg.Opts {
		l += int(v.Len) + 2
	}

	return l
}

// MarshalTo puts the byte sequence in the byte array given as b.
// SCCP is dependent on the Pointers when serializing, which means that it might fail when invalid Pointers are set.
func (msg *CC) MarshalTo(b []byte) error {
	b[0] = uint8(msg.Type)
	copy(b[1:4], utils.Uint32To24(msg.DestinationLocalReference))
	copy(b[4:7], utils.Uint32To24(msg.SourceLocalReference))
	b[7] = byte(msg.ProtocolClass)

	if len(msg.Opts) == 0 {
		return nil
	}

	b[8] = 1
	p := uint8(9)

	for i := 0; i < len(msg.Opts); i++ {
		b[p] = msg.Opts[i].Tag
		b[p+1] = msg.Opts[i].Len
		copy(b[p+2:], msg.Opts[i].Value)

		p += msg.Opts[i].Len + 2
	}
	return nil
}

func (msg *CC) String() string {
	if msg.CalledPartyAddress != nil {
		return fmt.Sprintf("{Type: CC, CalledPartyAddress: %v}", msg.CalledPartyAddress)
	}
	return "{Type: CC}"
}

// MessageType returns the Message Type in int.
func (msg *CC) MessageType() MsgType {
	return msg.Type
}

func (msg *CC) MessageTypeName() string {
	return "CC"
}
