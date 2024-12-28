package sccp

import (
	"encoding/hex"
	"fmt"
	"io"

	"github.com/wmnsk/go-sccp/params"
)

type DT1 struct {
	Type                      MsgType
	DestinationLocalReference params.LocalReference
	Segmenting                uint8
	Data                      []byte
}

func ParseDT1(b []byte) (*DT1, error) {
	msg := &DT1{}
	if err := msg.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return msg, nil
}

func (msg *DT1) UnmarshalBinary(b []byte) error {
	l := uint8(len(b))
	if l <= (1 + 3 + 1 + 1) {
		return io.ErrUnexpectedEOF
	}

	msg.Type = MsgType(b[0])
	if err := msg.DestinationLocalReference.Read(b[1:4]); err != nil {
		return err
	}

	msg.Segmenting = b[4]

	if b[5] != 1 { // pointer to var, ae next position
		return io.ErrUnexpectedEOF
	}

	dlen := b[6]
	if l != (dlen + 6 + 1) {
		return io.ErrUnexpectedEOF
	}

	msg.Data = b[7:]
	return nil
}

func (msg *DT1) MarshalBinary() ([]byte, error) {
	b := make([]byte, msg.MarshalLen())
	if err := msg.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

func (msg *DT1) MarshalLen() int {
	return len(msg.Data) + 7
}

func (msg *DT1) MarshalTo(b []byte) error {
	b[0] = uint8(msg.Type)
	msg.DestinationLocalReference.Read(b[1:4])
	b[4] = msg.Segmenting
	b[5] = 1
	b[6] = byte(len(msg.Data))
	copy(b[7:], msg.Data)
	return nil
}

func (msg *DT1) String() string {
	return fmt.Sprintf("{Type: DT1, DataLength: %d, Data: %s}", len(msg.Data), hex.EncodeToString(msg.Data))
}

// MessageType returns the Message Type in int.
func (msg *DT1) MessageType() MsgType {
	return msg.Type
}

func (msg *DT1) MessageTypeName() string {
	return "DT1"
}
