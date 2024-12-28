package params

import (
	"fmt"
	"io"
)

const (
	DataTag      uint8 = 0x0F
	CdPtyAddrTag uint8 = 0x03
	CgPtyAddrTag uint8 = 0x04
)

type LocalReference struct {
	Value uint32
} // just 24 bits used
func (lr *LocalReference) Read(b []byte) error {
	if len(b) != 3 {
		return fmt.Errorf("unable to read local reference: given bytes length is invalid")
	}
	b[0] = byte((lr.Value >> 16) & 0xFF)
	b[1] = byte((lr.Value >> 8) & 0xFF)
	b[2] = byte(lr.Value & 0xFF)
	return nil
}
func (lr *LocalReference) Write(b [3]byte) {
	lr.Value = (uint32(b[0])<<16)&0xFF0000 + (uint32(b[1])<<8)&0xFF00 + uint32(b[2])&0xFF
}
func (lr *LocalReference) String() string {
	return fmt.Sprintf("%05X", lr.Value)
}

type Optional struct {
	Tag   uint8
	Len   uint8
	Value []byte
}

func ParseOptional(b []byte) ([]*Optional, error) {
	p := uint8(0)
	opts := make([]*Optional, 0)
	for p < uint8(len(b)) {
		t := b[p]

		if t == 0 {
			return opts, nil
		}
		if (p + 1) >= uint8(len(b)) {
			return nil, io.ErrUnexpectedEOF
		}

		l := b[p+1]
		if (p + 1 + l) >= uint8(len(b)) {
			return nil, io.ErrUnexpectedEOF
		}

		o := &Optional{
			Tag:   t,
			Len:   l,
			Value: b[p+2 : p+2+l],
		}

		opts = append(opts, o)
		p += 2 + l

	}

	return opts, nil
}
