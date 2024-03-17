package slogleveler

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	LogLevelTrace = slog.LevelDebug - 1
	LogLevelFatal = slog.LevelError + 1
)

var (
	traceLevel, _ = NewLevel(LogLevelTrace, "TRACE")
	fatalLevel, _ = NewLevel(LogLevelFatal, "FATAL")
)

func TestLeveler(t *testing.T) {
	t.Run("rewrite log entry success", func(t *testing.T) {
		output := &bytes.Buffer{}

		leveler, err := New(
			WithWriter(output),
			WithLevel(traceLevel),
			WithLevel(fatalLevel))

		assert.Nil(t, err)

		slogger := slog.New(slog.NewTextHandler(leveler,
			&slog.HandlerOptions{
				Level: LogLevelTrace,
			}))

		slogger.Log(context.Background(), LogLevelTrace, "DEBUG-1 log 1", "key1", "value1")
		slogger.Log(context.Background(), LogLevelFatal, "ERROR+1 log 2", "key2", "value2")

		o := output.String()

		require.Contains(t, o, `level=TRACE msg="DEBUG-1 log 1" key1=value1`)
		require.Contains(t, o, `level=FATAL msg="ERROR+1 log 2" key2=value2`)
	})

	t.Run("rewrite log entry writer not defined", func(t *testing.T) {
		output := &bytes.Buffer{}

		leveler, err := New(
			WithLevel(traceLevel),
			WithLevel(fatalLevel))

		assert.Nil(t, err)

		slogger := slog.New(slog.NewTextHandler(leveler,
			&slog.HandlerOptions{
				Level: LogLevelTrace,
			}))

		slogger.Log(context.Background(), LogLevelTrace, "DEBUG-1 log 1", "key1", "value1")
		slogger.Log(context.Background(), LogLevelFatal, "ERROR+1 log 2", "key2", "value2")

		o := output.String()

		require.NotContains(t, o, `level=TRACE msg="DEBUG-1 log 1" key1=value1`)
		require.NotContains(t, o, `level=FATAL msg="ERROR+1 log 2" key2=value2`)
	})

	t.Run("replace attrs", func(t *testing.T) {
		output := &bytes.Buffer{}

		leveler, err := New(
			WithLevel(traceLevel),
			WithLevel(fatalLevel))

		assert.Nil(t, err)

		slogger := slog.New(slog.NewTextHandler(output,
			&slog.HandlerOptions{
				Level:       LogLevelTrace,
				ReplaceAttr: leveler.ReplaceLevels,
			}))

		slogger.Log(context.Background(), LogLevelTrace, "DEBUG-1 log 1", "key1", "value1")
		slogger.Log(context.Background(), LogLevelFatal, "ERROR+1 log 2", "key2", "value2")

		o := output.String()

		require.Contains(t, o, `level=TRACE msg="DEBUG-1 log 1" key1=value1`)
		require.Contains(t, o, `level=FATAL msg="ERROR+1 log 2" key2=value2`)
	})

	t.Run("leveler init errors", func(t *testing.T) {
		_, err := New(
			WithLevel(traceLevel),
			WithLevel(traceLevel))

		assert.NotNil(t, err)

		traceLevel2, _ := NewLevel(traceLevel.l, "TRACE2")

		_, err = New(
			WithLevel(traceLevel),
			WithLevel(traceLevel2),
		)

		assert.NotNil(t, err)

		traceLevel3, _ := NewLevel(-100, "TRACE")

		_, err = New(
			WithLevel(traceLevel),
			WithLevel(traceLevel3),
		)

		assert.NotNil(t, err)
	})

	t.Run("level", func(t *testing.T) {
		_, err := NewLevel(-100, "")
		assert.NotNil(t, err)

		validLevel, err := NewLevel(-100, "VERBOSE")
		assert.Nil(t, err)
		assert.Equal(t, slog.Level(-100), validLevel.l)
		assert.Equal(t, "VERBOSE", validLevel.s)
	})
}
