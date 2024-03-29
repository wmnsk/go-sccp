// Copyright 2019-2024 go-sccp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package params

// ProtocolClass is a Protocol Class SCCP parameter.
type ProtocolClass uint8

// NewProtocolClass creates a new ProtocolClass.
func NewProtocolClass(cls int, opts bool) ProtocolClass {
	if opts {
		return ProtocolClass(cls | 0x80)
	}
	return ProtocolClass(cls)
}

// Class returns the class part from ProtocolClass parameter.
func (p ProtocolClass) Class() int {
	return int(p) & 0xf
}

// ReturnOnError judges if ProtocolClass has "Return Message On Error" option.
func (p ProtocolClass) ReturnOnError() bool {
	return (int(p) >> 7) == 0
}
