// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gsl

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
