[![CircleCI](https://circleci.com/gh/spatialcurrent/go-sync-logger/tree/master.svg?style=svg)](https://circleci.com/gh/spatialcurrent/go-sync-logger/tree/master) [![Go Report Card](https://goreportcard.com/badge/spatialcurrent/go-sync-logger)](https://goreportcard.com/report/spatialcurrent/go-sync-logger)  [![GoDoc](https://godoc.org/github.com/spatialcurrent/go-sync-logger?status.svg)](https://godoc.org/github.com/spatialcurrent/go-sync-logger) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/spatialcurrent/go-sync-logger/blob/master/LICENSE)

# go-sync-logger

# Description

**go-sync-logger** (aka GRW) is a simple library for safely using


reading/writing of resources.  GRW can read from `bzip2`, `gzip`, `snappy`, and `zip` resources and write to `gzip` and `snappy` resources.

Using cross compilers, this library can also be called by other languages.  This library is cross compiled into a Shared Object file (`*.so`).  The Shared Object file can be called by `C`, `C++`, and `Python` on Linux machines.  See the examples folder for patterns that you can use.  This library is also compiled to pure `JavaScript` using [GopherJS](https://github.com/gopherjs/gopherjs).

# Usage

**Go**

You can import **go-sync-logger** as a library with:

```go
import (
  "github.com/spatialcurrent/go-sync-logger/gsl"
  "github.com/spatialcurrent/go-reader-writer/grw"
)
...

... () {
  errorWriter, err := grw.WriteToResource(errorDestination, errorCompression, true, s3_client)
	if err != nil {
		fmt.Println(errors.Wrap(err, "error creating error writer"))
		os.Exit(1)
	}

	infoWriter, err := grw.WriteToResource(infoDestination, infoCompression, true, s3_client)
	if err != nil {
		errorWriter.WriteError(errors.Wrap(err, "error creating log writer")) // #nosec
		errorWriter.Close()                                                   // #nosec
		os.Exit(1)
	}

	logger := gsl.NewLogger(
		map[string]int{"info": 0, "error": 1},
		[]grw.ByteWriteCloser{infoWriter, errorWriter},
		[]string{infoFormat, errorFormat},
	)
}
```

See [gsl](https://godoc.org/github.com/spatialcurrent/go-sync-logger/gsl) in GoDoc for information on how to use Go API.

# Contributing

[Spatial Current, Inc.](https://spatialcurrent.io) is currently accepting pull requests for this repository.  We'd love to have your contributions!  Please see [Contributing.md](https://github.com/spatialcurrent/go-sync-logger/blob/master/CONTRIBUTING.md) for how to get started.

# License

This work is distributed under the **MIT License**.  See **LICENSE** file.
