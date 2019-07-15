// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gsl

import (
	"fmt"
	"os"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/spatialcurrent/go-reader-writer/grw"
)

// CreateApplicationLoggerInput holds the input for the CreateApplicationLogger function..
type CreateApplicationLoggerInput struct {
	ErrorDestination string
	ErrorCompression string
	ErrorFormat      string
	InfoDestination  string
	InfoCompression  string
	InfoFormat       string
	Verbose          bool
}

// CreateApplicationLogger creates a new *Logger given the fields in *CreateApplicationLoggerInput.
// This function creates a logger that intuitively works as you would expect an application logger to work.
// The logger shares a single grw.ByteWriteCloser if error and info messages are going to the same location.
// If verbose mode is on, warn messages are sent to the error log and debug messages are sent to the info log.
// If there is an error during creation then the program prints the error and exits with exit code 1.
func CreateApplicationLogger(input *CreateApplicationLoggerInput) *Logger {

	errorWriter, err := grw.WriteToResource(input.ErrorDestination, input.ErrorCompression, true, nil)
	if err != nil {
		fmt.Println(errors.Wrap(err, "error creating error writer"))
		os.Exit(1)
	}

	levels := map[string]int{"error": 0, "fatal": 0}
	writers := []Writer{errorWriter}
	formats := []string{input.ErrorFormat}

	if input.Verbose {
		levels["warn"] = 0
	}

	if len(input.InfoDestination) > 0 && input.InfoDestination != "/dev/null" && input.InfoDestination != "null" {
		if input.InfoDestination == input.ErrorDestination {
			if input.InfoFormat != input.ErrorFormat {
				_, err := errorWriter.WriteError(fmt.Errorf("info-format ( %s ) and error-format ( %s ) must match when they share a destination", input.InfoFormat, input.ErrorFormat)) // #nosec
				if err != nil {
					fmt.Println(err.Error())
				}
				errorWriter.Close() // #nosec
				os.Exit(1)
			}
			if input.InfoCompression != input.ErrorCompression {
				_, err := errorWriter.WriteError(fmt.Errorf("info-compression ( %s ) and error-compression ( %s ) must match when they share a destination", input.InfoCompression, input.ErrorCompression)) // #nosec
				if err != nil {
					fmt.Println(err.Error())
				}
				errorWriter.Close() // #nosec
				os.Exit(1)
			}

			levels["info"] = 0
			if input.Verbose {
				levels["debug"] = 0
			}
		} else {
			infoWriter, err := grw.WriteToResource(input.InfoDestination, input.InfoCompression, true, nil)
			if err != nil {
				errorWriter.WriteError(errors.Wrap(err, "error creating log writer")) // #nosec
				errorWriter.Close()                                                   // #nosec
				os.Exit(1)
			}

			levels["info"] = 1
			writers = append(writers, infoWriter)
			formats = append(formats, input.InfoFormat)

			if input.Verbose {
				levels["debug"] = 1
			}
		}
	}

	logger := NewLogger(levels, writers, formats, true)
	return logger
}
