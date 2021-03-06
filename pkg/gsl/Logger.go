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
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/spatialcurrent/go-simple-serializer/pkg/gss"
)

// Logger contains a slice of writers, a slice of matching formats, and a mapping of levels to writers.
type Logger struct {
	levels          map[string]int // level --> position in writers
	writers         []Writer       // list of writers
	formats         []string       // list of formats for each writer
	LevelField      string         // the key for the level field
	TimeStampField  string         // the key for the timestamp field
	TimeStampFormat string         // the format for the timestamp field
	MessageField    string         // the key for the message field
	AutoFlush       bool           // flush after every message
}

// NewLogger returns a new logger with the given configuration and default field keys.
// Set autoFlush to true to flush the buffer to the underlying writer after every message.
func NewLogger(levels map[string]int, writers []Writer, formats []string, autoFlush bool) *Logger {
	return &Logger{
		levels:          levels,
		writers:         writers,
		formats:         formats,
		TimeStampField:  "ts",
		TimeStampFormat: time.RFC3339,
		LevelField:      "level",
		MessageField:    "msg",
		AutoFlush:       autoFlush,
	}
}

// Debug writes the provided object to the `debug` writer.
// If no `debug` writer exists, then return an ErrUnknownLevel error.
func (l *Logger) Debug(obj interface{}) error {
	level := "debug"
	position, ok := l.levels[level]
	if !ok {
		return &ErrUnknownLevel{Level: level}
	}
	_, err := l.WriteLineSafe(level, obj, l.writers[position], l.formats[position])
	if err != nil {
		return errors.Wrap(err, "error writing debug message")
	}
	return nil
}

// DebugF writes the provided message to the `debug` writer.
// The message is generated using `fmt.Sprintf(format, values...)`.
// If no `debug` writer exists, then return an ErrUnknownLevel error.
func (l *Logger) DebugF(format string, values ...interface{}) error {
	level := "debug"
	position, ok := l.levels[level]
	if !ok {
		return &ErrUnknownLevel{Level: level}
	}
	_, err := l.WriteLineSafe(level, fmt.Sprintf(format, values...), l.writers[position], l.formats[position])
	if err != nil {
		return errors.Wrap(err, "error writing info message")
	}
	return nil
}

// Info writes the provided object to the `info` writer.
// If no `info` writer exists, then return an ErrUnknownLevel error.
func (l *Logger) Info(obj interface{}) error {
	level := "info"
	position, ok := l.levels[level]
	if !ok {
		return &ErrUnknownLevel{Level: level}
	}
	_, err := l.WriteLineSafe(level, obj, l.writers[position], l.formats[position])
	if err != nil {
		return errors.Wrap(err, "error writing info message")
	}
	return nil
}

// InfoF writes the provided message to the `info` writer.
// The message is generated using `fmt.Sprintf(format, values...)`.
// If no `info` writer exists, then return an ErrUnknownLevel error.
func (l *Logger) InfoF(format string, values ...interface{}) error {
	level := "info"
	position, ok := l.levels[level]
	if !ok {
		return &ErrUnknownLevel{Level: level}
	}
	_, err := l.WriteLineSafe(level, fmt.Sprintf(format, values...), l.writers[position], l.formats[position])
	if err != nil {
		return errors.Wrap(err, "error writing info message")
	}
	return nil
}

// Warn writes the provided object to the `warn` writer.
// If no `warn` writer exists, then return an ErrUnknownLevel error.
func (l *Logger) Warn(obj interface{}) error {
	level := "warn"
	position, ok := l.levels[level]
	if !ok {
		return &ErrUnknownLevel{Level: level}
	}
	_, err := l.WriteLineSafe(level, obj, l.writers[position], l.formats[position])
	if err != nil {
		return errors.Wrap(err, "error writing warn message")
	}
	return nil
}

// Error writes the provided object to the `error` writer.
// If no `error` writer exists, then return an ErrUnknownLevel error.
func (l *Logger) Error(obj interface{}) error {
	level := "error"
	position, ok := l.levels[level]
	if !ok {
		return &ErrUnknownLevel{Level: level}
	}
	_, err := l.WriteLineSafe(level, obj, l.writers[position], l.formats[position])
	if err != nil {
		return errors.Wrap(err, "error writing error message")
	}
	return nil
}

// ErrorF writes the provided message to the `error` writer.
// The message is generated using `fmt.Sprintf(format, values...)`.
// If no `error` writer exists, then return an ErrUnknownLevel error.
func (l *Logger) ErrorF(format string, values ...interface{}) error {
	level := "error"
	position, ok := l.levels[level]
	if !ok {
		return &ErrUnknownLevel{Level: level}
	}
	_, err := l.WriteLineSafe(level, fmt.Sprintf(format, values...), l.writers[position], l.formats[position])
	if err != nil {
		return errors.Wrap(err, "error writing info message")
	}
	return nil
}

// Fatal locks all the writers, flushes them, writes the given message to the fatal writer, flushes the writers again, closes the writers, unlocks the writers, and finally exits with code 1.
func (l *Logger) Fatal(obj interface{}) {
	for _, w := range l.writers {
		w.Lock()
	}
	for _, w := range l.writers {
		w.Flush() // #nosec
	}
	if position, ok := l.levels["error"]; ok {
		l.WriteLine("fatal", obj, l.writers[position], l.formats[position]) // #nosec
	}
	for _, w := range l.writers {
		w.Flush() // #nosec
	}
	for _, w := range l.writers {
		w.Close() // #nosec
	}
	for _, w := range l.writers {
		w.Unlock()
	}
	os.Exit(1)
}

// FatalF writes a message to the fatal writer using the provide format and values.
// The message is generated using `fmt.Sprintf(format, values...)`.
func (l *Logger) FatalF(format string, values ...interface{}) {
	l.Fatal(fmt.Sprintf(format, values...))
}

// FlushSafe flushes all the writers using concurrency-safe methods.
func (l *Logger) Flush() error {
	for _, w := range l.writers {
		err := w.FlushSafe()
		if err != nil {
			return err
		}
	}
	return nil
}

// Close locks all the writers, flushes them, closes them, and then unlocks them.
func (l *Logger) Close() {
	for _, w := range l.writers {
		w.Lock()
	}
	for _, w := range l.writers {
		w.Flush() // #nosec
	}
	for _, w := range l.writers {
		w.Close() // #nosec
	}
	for _, w := range l.writers {
		w.Unlock()
	}
}

// ListenError listens for  message on a `chan interface{}` channel and writes them to the info writer.
// If a *sync.WaitGroup is not nil, then it is marked as done once the channel is closed.
func (l *Logger) ListenInfo(messages chan interface{}, wg *sync.WaitGroup) {
	go func(messages chan interface{}) {
		for message := range messages {
			err := l.Info(message)
			if err != nil {
				l.Error(err) // #nosec
			}
			l.Flush()
		}
		if wg != nil {
			wg.Done()
		}
	}(messages)
}

// ListenError listens for  message on a `chan interface{}` channel and writes them to the error writer.
// If a *sync.WaitGroup is not nil, then it is marked as done once the channel is closed.
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

// ListenFatal listens for a message on a `chan interface{}` channel.
// Once a message is received, the logger immediately calls l.Fatal(), which writes the fatal error, flushes the logs, closes the logs, and finally exits with code 1.
func (l *Logger) ListenFatal(messages chan interface{}) {
	go func() {
		for msg := range messages {
			l.Fatal(msg)
		}
	}()
}

// FormatObject formats a given object using a given level and format and returns the formatted bytes and error, if any.
func (l *Logger) FormatObject(level string, obj interface{}, format string) ([]byte, error) {
	if err, ok := obj.(error); ok {
		m := map[string]interface{}{}
		h := make([]interface{}, 0)
		if len(l.LevelField) > 0 {
			m[l.LevelField] = level
			h = append(h, l.LevelField)
		}
		if _, ok := m[l.TimeStampField]; !ok && len(l.TimeStampField) > 0 {
			m[l.TimeStampField] = time.Now().Format(l.TimeStampFormat)
			h = append(h, l.TimeStampField)
		}
		if len(l.MessageField) > 0 {
			m[l.MessageField] = strings.Replace(err.Error(), "\n", ": ", -1)
			h = append(h, l.MessageField)
		}
		return gss.SerializeBytes(&gss.SerializeBytesInput{
			Object:            m,
			Format:            format,
			Header:            h,
			ExpandHeader:      true,
			Limit:             gss.NoLimit,
			KeyValueSeparator: "=",
			LineSeparator:     "\n",
			Sorted:            true,
			Pretty:            false,
		})
	} else if msg, ok := obj.(string); ok {
		m := map[string]interface{}{}
		h := make([]interface{}, 0)
		if len(l.LevelField) > 0 {
			m[l.LevelField] = level
			h = append(h, l.LevelField)
		}
		if len(l.TimeStampField) > 0 {
			m[l.TimeStampField] = time.Now().Format(l.TimeStampFormat)
			h = append(h, l.TimeStampField)
		}
		if len(l.MessageField) > 0 {
			m[l.MessageField] = msg
			h = append(h, l.MessageField)
		}
		return gss.SerializeBytes(&gss.SerializeBytesInput{
			Object:            m,
			Format:            format,
			Header:            h,
			ExpandHeader:      true,
			Limit:             gss.NoLimit,
			KeyValueSeparator: "=",
			LineSeparator:     "\n",
			Sorted:            true,
			Pretty:            false,
		})
	} else if m, ok := obj.(map[string]string); ok {
		h := make([]interface{}, 0)
		if len(l.LevelField) > 0 {
			m[l.LevelField] = level
			h = append(h, l.LevelField)
		}
		if len(l.TimeStampField) > 0 {
			m[l.TimeStampField] = time.Now().Format(l.TimeStampFormat)
			h = append(h, l.TimeStampField)
		}
		return gss.SerializeBytes(&gss.SerializeBytesInput{
			Object:            m,
			Format:            format,
			Header:            h,
			ExpandHeader:      true,
			Limit:             gss.NoLimit,
			KeyValueSeparator: "=",
			LineSeparator:     "\n",
			Sorted:            true,
			Pretty:            false,
		})
	} else if m, ok := obj.(map[string]interface{}); ok {
		if len(l.LevelField) > 0 {
			m[l.LevelField] = level
		}
		if len(l.TimeStampField) > 0 {
			m[l.TimeStampField] = time.Now().Format(l.TimeStampFormat)
		}
		return gss.SerializeBytes(&gss.SerializeBytesInput{
			Object:            m,
			Format:            format,
			Header:            gss.NoHeader,
			ExpandHeader:      true,
			Limit:             gss.NoLimit,
			KeyValueSeparator: "=",
			LineSeparator:     "\n",
			Sorted:            true,
			Pretty:            false,
		})
	}
	return gss.SerializeBytes(&gss.SerializeBytesInput{
		Object:            obj,
		Format:            format,
		Header:            gss.NoHeader,
		ExpandHeader:      true,
		Limit:             gss.NoLimit,
		KeyValueSeparator: "=",
		LineSeparator:     "\n",
		Sorted:            true,
		Pretty:            false,
	})
}

// WriteLine formats the given object using FormatObject then writes the formatted string with a trailing newline to the matching grw.ByteWriteCloser and returns an error, if any.
// WriteLine calls the writer's WriteLine method, which does not lock the underlying writer.
// The writer can already be locked.
//  - https://godoc.org/github.com/spatialcurrent/go-reader-writer/grw#Writer.WriteLine
//  - https://godoc.org/io#Writer
//  - https://godoc.org/sync#Mutex
func (l *Logger) WriteLine(level string, obj interface{}, writer Writer, format string) (n int, err error) {
	line, err := l.FormatObject(level, obj, format)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("error formating object at level %s using format %s", level, format))
	}
	if len(line) > 0 {
		_, err := writer.WriteLine(string(line))
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("error formating object at level %s using format %s", level, format))
		}
	}
	return 0, nil
}

// WriteLineSafe formats the given object using FormatObject then writes the formatted string with a trailing newline to the matching grw.ByteWriteCloser and returns an error, if any.
// WriteLineSafe calls the writer's WriteLineSafe method, which locks the underlying writer for the duration of writing using a sync.Mutex.
//  - https://godoc.org/github.com/spatialcurrent/go-reader-writer/grw#Writer.WriteLineSafe
//  - https://godoc.org/io#Writer
//  - https://godoc.org/sync#Mutex
func (l *Logger) WriteLineSafe(level string, obj interface{}, writer Writer, format string) (n int, err error) {
	line, err := l.FormatObject(level, obj, format)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("error formating object at level %s using format %s", level, format))
	}
	if len(line) > 0 {
		_, err := writer.WriteLineSafe(string(line))
		if err != nil {
			return 0, errors.Wrap(err, fmt.Sprintf("error formating object at level %s using format %s", level, format))
		}
	}
	if l.AutoFlush {
		err := writer.Flush()
		if err != nil {
			return 0, errors.Wrap(err, "error flushing after writing line")
		}
	}
	return 0, nil
}
