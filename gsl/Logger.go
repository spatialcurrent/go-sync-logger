package gsl

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/spatialcurrent/go-reader-writer/grw"
	"github.com/spatialcurrent/go-simple-serializer/gss"
)

type Logger struct {
	levels  map[string]int        // level --> position in writers
	writers []grw.ByteWriteCloser // list of writers
	formats []string              // list of formats for each writer
}

func NewLogger(levels map[string]int, writers []grw.ByteWriteCloser, formats []string) *Logger {
	return &Logger{
		levels:  levels,
		writers: writers,
		formats: formats,
	}
}

func (l *Logger) Debug(obj interface{}) error {
	level := "debug"
	position, ok := l.levels[level]
	if !ok {
		return &ErrUnknownLevel{Level: level}
	}
	_, err := WriteLineSafe("debug", obj, l.writers[position], l.formats[position])
	if err != nil {
		return errors.Wrap(err, "error writing debug message")
	}
	return nil
}

func (l *Logger) Info(obj interface{}) error {
	level := "info"
	position, ok := l.levels[level]
	if !ok {
		return &ErrUnknownLevel{Level: level}
	}
	_, err := WriteLineSafe("info", obj, l.writers[position], l.formats[position])
	if err != nil {
		return errors.Wrap(err, "error writing info message")
	}
	return nil
}

func (l *Logger) Warn(obj interface{}) error {
	level := "warn"
	position, ok := l.levels[level]
	if !ok {
		return &ErrUnknownLevel{Level: level}
	}
	_, err := WriteLineSafe("warn", obj, l.writers[position], l.formats[position])
	if err != nil {
		return errors.Wrap(err, "error writing warn message")
	}
	return nil
}

func (l *Logger) Error(obj interface{}) error {
	level := "error"
	position, ok := l.levels[level]
	if !ok {
		return &ErrUnknownLevel{Level: level}
	}
	_, err := WriteLineSafe("error", obj, l.writers[position], l.formats[position])
	if err != nil {
		return errors.Wrap(err, "error writing error message")
	}
	return nil
}

func (l *Logger) Fatal(obj interface{}) {
	for _, w := range l.writers {
		w.Lock()
	}
	for _, w := range l.writers {
		w.Flush() // #nosec
	}
	l.Error(obj) // #nosec
	for _, w := range l.writers {
		w.Flush() // #nosec
	}
	for _, w := range l.writers {
		w.Close() // #nosec
	}
	os.Exit(1)
}

func (l *Logger) Flush() {
	for _, w := range l.writers {
		w.FlushSafe() // #nosec
	}
}

func (l *Logger) Close() {
	for _, w := range l.writers {
		w.FlushSafe() // #nosec
	}
	for _, w := range l.writers {
		w.CloseSafe() // #nosec
	}
}

func (l *Logger) ListenInfo(messages chan interface{}, wg *sync.WaitGroup) {
	go func(messages chan interface{}) {
		for message := range messages {
			l.Info(message) // #nosec
			l.Flush()
		}
		if wg != nil {
			wg.Done()
		}
	}(messages)
}

func (l *Logger) ListenError(messages chan interface{}, wg *sync.WaitGroup) {
	go func(messages chan interface{}) {
		for message := range messages {
			l.Error(message) // #nosec
			l.Flush()
		}
		if wg != nil {
			wg.Done()
		}
	}(messages)
}

func (l *Logger) ListenFatal(c chan error) {
	go func() {
		for err := range c {
			l.Fatal(err)
		}
	}()
}

// FormatObject formats a given object using a given level and format and returns the formatted string and error, if any.
func FormatObject(level string, obj interface{}, format string) (string, error) {
	if err, ok := obj.(error); ok {
		m := map[string]interface{}{
			"level": level,
			"msg":   strings.Replace(err.Error(), "\n", ": ", -1),
			"ts":    time.Now().Format(time.RFC3339),
		}
		return gss.SerializeString(m, format, gss.NoHeader, gss.NoLimit)
	} else if msg, ok := obj.(string); ok {
		m := map[string]interface{}{
			"level": level,
			"msg":   msg,
			"ts":    time.Now().Format(time.RFC3339),
		}
		return gss.SerializeString(m, format, gss.NoHeader, gss.NoLimit)
	} else if m, ok := obj.(map[string]string); ok {
		m["level"] = level
		m["ts"] = time.Now().Format(time.RFC3339)
		return gss.SerializeString(m, format, gss.NoHeader, gss.NoLimit)
	} else if m, ok := obj.(map[string]interface{}); ok {
		m["level"] = level
		m["ts"] = time.Now().Format(time.RFC3339)
		return gss.SerializeString(m, format, gss.NoHeader, gss.NoLimit)
	}
	return gss.SerializeString(obj, format, gss.NoHeader, gss.NoLimit)
}

// WriteLineSafe formats the given object using FormatObject then writes the formatted string with a trailing newline to the matching grw.ByteWriteCloser and returns an error, if any.
// WriteLineSafe also locks the underlying writer for the duration of writing using a sync.Mutex.
//  - https://godoc.org/github.com/spatialcurrent/go-reader-writer/grw#Writer.WriteLineSafe
//  - https://godoc.org/io#Writer
//  - https://godoc.org/sync#Mutex
func WriteLineSafe(level string, obj interface{}, writer grw.ByteWriteCloser, format string) (n int, err error) {
	line, err := FormatObject(level, obj, format)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("error formating object at level %s using format %s", level, format))
	}
	if len(line) > 0 {
		_, err := writer.WriteLineSafe(line)
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("error formating object at level %s using format %s", level, format))
		}
	}
	return 0, nil
}
