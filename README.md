[![CircleCI](https://circleci.com/gh/spatialcurrent/go-sync-logger/tree/master.svg?style=svg)](https://circleci.com/gh/spatialcurrent/go-sync-logger/tree/master) [![Go Report Card](https://goreportcard.com/badge/spatialcurrent/go-sync-logger)](https://goreportcard.com/report/spatialcurrent/go-sync-logger)  [![GoDoc](https://godoc.org/github.com/spatialcurrent/go-sync-logger?status.svg)](https://godoc.org/github.com/spatialcurrent/go-sync-logger) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/spatialcurrent/go-sync-logger/blob/master/LICENSE)

# go-sync-logger

# Description

**go-sync-logger** (aka GSL) is a simple concurrency-safe library for sharing writers (e.g., stdout, errors.log, etc.) among multiple log levels.  With **GSL** you can have `info` and `error` messages share a log with messages being written in the proper sequence.

 **go-sync-logger** depends on:
- [go-reader-writer](https://github.com/spatialcurrent/go-reader-writer) for writing to arbitrary locations.
- [go-simple-serializer](https://github.com/spatialcurrent/go-simple-serializer) for serializing log messages.

# Usage

**Go**

You can import **go-sync-logger** as a library with:

```go
import (
  "github.com/spatialcurrent/go-sync-logger/gsl"
)
```

To initialize the logger, use `gsl.NewLogger` and pass it the writers, formats, and proper mapping of levels to writers.  For example, the configuration below, has `error` and `warn` share a writer.

```go
... () {
  logger := gsl.NewLogger(
    map[string]int{"info": 0, "error": 1, "warn": 1},
    []grw.ByteWriteCloser{infoWriter, errorWriter},
    []string{infoFormat, errorFormat},
  )
}
```

For a complete example on how to initialize the logger using configuration provided by [viper](https://github.com/spf13/viper) see [viper.md](https://github.com/spatialcurrent/go-sync-logger/tree/master/example/viper.md) in [examples](https://github.com/spatialcurrent/go-sync-logger/tree/master/example).

See [gsl](https://godoc.org/github.com/spatialcurrent/go-sync-logger/gsl) in GoDoc for information on how to use Go API.

# Contributing

[Spatial Current, Inc.](https://spatialcurrent.io) is currently accepting pull requests for this repository.  We'd love to have your contributions!  Please see [Contributing.md](https://github.com/spatialcurrent/go-sync-logger/blob/master/CONTRIBUTING.md) for how to get started.

# License

This work is distributed under the **MIT License**.  See **LICENSE** file.
