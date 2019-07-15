// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gsl

// Interface containing the methods required for underlying writers.
// This interface is implemented by go-reader-writer.
//  - https://github.com/spatialcurrent/go-reader-writer
type Writer interface {
	WriteLine(str string) (int, error)     // write line to underlying writer
	WriteLineSafe(str string) (int, error) // lock underlying writer, write line, and then unlock
	Lock()                                 // lock underlying writer
	Unlock()                               // unlock underlying writer
	Flush() error                          // flush buffer to underlying writer
	FlushSafe() error                      // lock underlying writer, flush buffer, and then unlock
	Close() error                          // lock all the underlying writers, flush their buffers, close all the writers, and then unlock.
}
