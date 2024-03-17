package slogleveler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"slices"
	"sync"

	"github.com/valyala/bytebufferpool"
)

var (
	ErrLevelStringEmpty      = errors.New("level string representation is empty")
	ErrWriterNotDefined      = errors.New("writer not defined")
	ErrLevelMapsNotDefined   = errors.New("level maps not defined")
	ErrLevelAlreadyPresented = errors.New("level already presented")
)

type Level struct {
	l slog.Level
	s string
}

type Leveler struct {
	levelsBytes map[slog.Level][]byte
	levelsStr   map[slog.Level]string
	w           io.Writer
	mu          sync.Mutex
}

func New(opts ...func(*Leveler) error) (*Leveler, error) {
	r := &Leveler{
		levelsBytes: make(map[slog.Level][]byte),
		levelsStr:   make(map[slog.Level]string),
		w:           io.Discard,
	}
	for _, o := range opts {
		if err := o(r); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func WithWriter(w io.Writer) func(*Leveler) error {
	return func(l *Leveler) error {
		if w == nil {
			return ErrWriterNotDefined
		}

		l.mu.Lock()
		defer l.mu.Unlock()

		l.w = w
		return nil
	}
}

func WithLevel(level Level) func(*Leveler) error {
	return func(l *Leveler) error {
		l.mu.Lock()
		defer l.mu.Unlock()

		if l.levelsBytes == nil || l.levelsStr == nil {
			return ErrLevelMapsNotDefined
		}

		for k, v := range l.levelsStr {
			if k == level.l {
				return fmt.Errorf("check slog level: %v: %w", k, ErrLevelAlreadyPresented)
			}
			if v == level.s {
				return fmt.Errorf("check slog level string: %s: %w", v, ErrLevelAlreadyPresented)
			}
		}

		for k, v := range l.levelsBytes {
			if k == level.l {
				return fmt.Errorf("check slog level: %v: %w", k, ErrLevelAlreadyPresented)
			}
			if slices.Equal(v, []byte(level.s)) {
				return fmt.Errorf("check slog level string: %s: %w", v, ErrLevelAlreadyPresented)
			}
		}

		l.levelsBytes[level.l] = []byte(level.s)
		l.levelsStr[level.l] = level.s
		return nil
	}
}

func NewLevel(l slog.Level, s string) (Level, error) {
	if s == "" {
		return Level{}, ErrLevelStringEmpty
	}
	return Level{
		l: l,
		s: s,
	}, nil
}

func (l *Leveler) Write(p []byte) (int, error) {
	buf := bytebufferpool.Get()
	buf.Set(p)
	for level, levelBytes := range l.levelsBytes {
		buf.Set(bytes.Replace(buf.B, []byte(level.String()), levelBytes, 1))
	}
	return l.w.Write(buf.B)
}

func (l *Leveler) ReplaceLevels(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.LevelKey {
		level := a.Value.Any().(slog.Level)
		levelStr, ok := l.levelsStr[level]
		if !ok {
			return a
		}

		a.Value = slog.StringValue(levelStr)
	}

	return a
}
