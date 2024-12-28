package sccp

import (
	"io"

	"github.com/wmnsk/go-sccp/params"
)

type RLC struct {
	Type                      MsgType
	DestinationLocalReference params.LocalReference
	SourceLocalReference      params.LocalReference
}

func ParseRLC(b []byte) (*RLC, error) {
	msg := &RLC{}
	if err := msg.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return msg, nil
}

func (msg *RLC) UnmarshalBinary(b []byte) error {
	l := uint8(len(b))
	if l != 7 {
		return io.ErrUnexpectedEOF
	}

	msg.Type = MsgType(b[0])
	if err := msg.DestinationLocalReference.Read(b[1:4]); err != nil {
		return err
	}
	return msg.SourceLocalReference.Read(b[4:])
}

func (msg *RLC) MarshalBinary() ([]byte, error) {
	b := make([]byte, msg.MarshalLen())
	if err := msg.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

func (msg *RLC) MarshalLen() int {
	return 7
}

func (msg *RLC) MarshalTo(b []byte) error {
	b[0] = uint8(msg.Type)
	msg.DestinationLocalReference.Read(b[1:4])
	msg.SourceLocalReference.Read(b[4:])
	return nil
}

func (msg *RLC) String() string {
	return "{Type: RLC}"
}

// MessageType returns the Message Type in int.
func (msg *RLC) MessageType() MsgType {
	return msg.Type
}

func (msg *RLC) MessageTypeName() string {
	return "RLC"
}
