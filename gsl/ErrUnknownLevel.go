// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gsl

import (
	"fmt"
)

type ErrUnknownLevel struct {
	Level string
}

func (e *ErrUnknownLevel) Error() string {
	return fmt.Sprintf("unknown level %s", e.Level)
}
