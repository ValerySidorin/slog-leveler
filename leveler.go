package slogleveler

import (
	"log/slog"
)

var levels = make(map[slog.Level]string)

func AddLevel(level slog.Level, levelStr string) {
	levels[level] = levelStr
}

func ReplaceLevels(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey {
		level := a.Value.Any().(slog.Level)
		levelStr, ok := levels[level]
		if !ok {
			return a
		}

		a.Value = slog.StringValue(levelStr)
	}

	return a
}
