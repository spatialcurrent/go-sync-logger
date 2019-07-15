// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package gsl

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/spatialcurrent/go-reader-writer/grw"
	"github.com/spatialcurrent/go-simple-serializer/pkg/tags"
)

const (
	testMessage = "test message"
)

func TestLoggerInfo(t *testing.T) {

	w, b := grw.WriteMemoryBytes()

	levels := map[string]int{"info": 0}
	writers := []Writer{w}
	formats := []string{"json"}
	autoFlush := false

	l := NewLogger(levels, writers, formats, autoFlush)

	err := l.Info(testMessage)
	assert.NoError(t, err)

	err = l.Flush()
	assert.NoError(t, err)

	outBytes := b.Bytes()
	assert.NotEmpty(t, outBytes)

	outObject := map[string]interface{}{}
	err = json.Unmarshal(outBytes, &outObject)
	assert.NoError(t, err)

	assert.Equal(t, "info", outObject["level"])
	assert.Equal(t, testMessage, outObject["msg"])
}

func TestLoggerInfoAutoFlush(t *testing.T) {

	w, b := grw.WriteMemoryBytes()

	levels := map[string]int{"info": 0}
	writers := []Writer{w}
	formats := []string{"json"}
	autoFlush := true

	l := NewLogger(levels, writers, formats, autoFlush)

	err := l.Info(testMessage)
	assert.NoError(t, err)

	outBytes := b.Bytes()
	assert.NotEmpty(t, outBytes)

	outObject := map[string]interface{}{}
	err = json.Unmarshal(outBytes, &outObject)
	assert.NoError(t, err)

	assert.Equal(t, "info", outObject["level"])
	assert.Equal(t, testMessage, outObject["msg"])
}

func TestLoggerInfoAutoFlushTags(t *testing.T) {

	w, b := grw.WriteMemoryBytes()

	levels := map[string]int{"info": 0}
	writers := []Writer{w}
	formats := []string{"tags"}
	autoFlush := true

	l := NewLogger(levels, writers, formats, autoFlush)

	err := l.Info(testMessage)
	assert.NoError(t, err)

	outBytes := b.Bytes()
	assert.NotEmpty(t, outBytes)

	fmt.Println("Yo:", string(outBytes))

	outObject, err := tags.UnmarshalType(outBytes, '=', reflect.TypeOf(map[string]string{}))
	assert.NoError(t, err)
	assert.IsType(t, map[string]string{}, outObject)

	if m, ok := outObject.(map[string]string); ok {
		assert.Equal(t, "info", m["level"])
		assert.Equal(t, testMessage, m["msg"])
	}
}

func TestLoggerInfoAutoFlushTagsMap(t *testing.T) {

	w, b := grw.WriteMemoryBytes()

	levels := map[string]int{"info": 0}
	writers := []Writer{w}
	formats := []string{"tags"}
	autoFlush := true

	l := NewLogger(levels, writers, formats, autoFlush)

	err := l.Info(map[string]string{
		"a": "x",
		"b": "y",
		"c": "z",
		"d": "x",
		"e": "y",
		"f": "z",
	})
	assert.NoError(t, err)

	outBytes := b.Bytes()
	assert.NotEmpty(t, outBytes)

	outObject, err := tags.UnmarshalType(outBytes, '=', reflect.TypeOf(map[string]string{}))
	assert.NoError(t, err)
	assert.IsType(t, map[string]string{}, outObject)

	if m, ok := outObject.(map[string]string); ok {
		assert.Equal(t, "info", m["level"])
		assert.Equal(t, "x", m["a"])
		assert.Equal(t, "y", m["b"])
		assert.Equal(t, "z", m["c"])
	}
}
