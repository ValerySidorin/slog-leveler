package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	slogleveler "github.com/ValerySidorin/slog-leveler"
)

const (
	LogLevelTrace = slog.LevelDebug - 1
	LogLevelFatal = slog.LevelError + 1
)

var (
	traceLevel, _ = slogleveler.NewLevel(LogLevelTrace, "TRACE")
	fatalLevel, _ = slogleveler.NewLevel(LogLevelFatal, "FATAL")
)

func main() {
	logWithReplacer()
	logWithReplaceLevels()
}

func logWithReplacer() {
	leveler, err := slogleveler.New(
		slogleveler.WithWriter(os.Stdout),
		slogleveler.WithLevel(traceLevel),
		slogleveler.WithLevel(fatalLevel))
	if err != nil {
		log.Fatal(err)
	}

	slogger := slog.New(slog.NewTextHandler(leveler,
		&slog.HandlerOptions{
			Level: LogLevelTrace,
		}))

	slogger.Log(context.Background(), LogLevelTrace, "DEBUG-1 log 1", "key", "value")
	slogger.Log(context.Background(), LogLevelFatal, "ERROR+1 log 1", "key", "value")
}

func logWithReplaceLevels() {
	leveler, err := slogleveler.New(
		slogleveler.WithLevel(traceLevel),
		slogleveler.WithLevel(fatalLevel))
	if err != nil {
		log.Fatal(err)
	}

	slogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: leveler.ReplaceLevels,
		Level:       LogLevelTrace,
	}))
	slogger.Log(context.Background(), LogLevelTrace, "DEBUG-1 log 2", "key", "value")
	slogger.Log(context.Background(), LogLevelFatal, "ERROR+1 log 2", "key", "value")
}
