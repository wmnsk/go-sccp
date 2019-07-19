// Copyright 2019 go-sccp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package sccp

import (
	"fmt"
)

// UnsupportedTypeError indicates the value in Version field is invalid.
type UnsupportedTypeError byte

// Error returns the type of receiver and some additional message.
func (e UnsupportedTypeError) Error() string {
	return fmt.Sprintf("sccp: got unsupported type %d", e)
}
