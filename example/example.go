package main

import (
	"context"
	"log/slog"
	"os"

	slogleveler "github.com/ValerySidorin/slog-leveler"
	"github.com/lmittmann/tint"
)

const (
	LogLevelTrace = slog.LevelDebug - 1
)

func main() {
	slogleveler.AddLevel(LogLevelTrace, "TRACE")
	slogger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{
		Level:       LogLevelTrace,
		ReplaceAttr: slogleveler.ReplaceLevels,
		AddSource:   true,
	}))

	slogger.Log(context.Background(), LogLevelTrace, "trace log", "key", "value")
}
