// Copyright 2019 go-sccp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package sccp

import (
	"errors"
	"fmt"
)

// ErrTooShortToDecode indicates the length of user input is too short to be decoded.
var ErrTooShortToDecode = errors.New("too short to decode")

// ErrTooShortToSerialize indicates the length of user input is too short to be serialized.
var ErrTooShortToSerialize = errors.New("too short to serialize")

// ErrUnsupportedType indicates the value in Version field is invalid.
type ErrUnsupportedType struct {
	Type string
	Msg  string
}

// Error returns the type of receiver and some additional message.
func (e *ErrUnsupportedType) Error() string {
	return fmt.Sprintf("%s: got unsupported type: %s", e.Msg, e.Type)
}
