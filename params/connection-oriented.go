package params

import (
	"io"
)

const (
	DataTag      uint8 = 0x0F
	CdPtyAddrTag uint8 = 0x03
	CgPtyAddrTag uint8 = 0x04
)

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
