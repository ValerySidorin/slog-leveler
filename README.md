# slog: log level toolset

Allows to add custom user-defined log levels to [slog](https://pkg.go.dev/log/slog).

```
go get github.com/ValerySidorin/slog-leveler
```

**Another slog helpers:**
- [shslog](https://github.com/ValerySidorin/shclog): `hclog.Logger` adapter for `*slog.Logger`

## 💡 Usage

### ReplaceLevels function
This is the preferred method. Use it, when your handler options allow you to define custom ReplaceAttr function.

```go
import (
	"context"
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

leveler, _ := slogleveler.New(
	slogleveler.WithLevel(traceLevel),
	slogleveler.WithLevel(fatalLevel))

slogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	ReplaceAttr: leveler.ReplaceLevels,
	Level:       LogLevelTrace,
}))
slogger.Log(context.Background(), LogLevelTrace, "DEBUG-1 log 2", "key", "value")
slogger.Log(context.Background(), LogLevelFatal, "ERROR+1 log 2", "key", "value")
```

### Log level rewriter
Leveler can rewrite log level and pass entry to next io.Writer.
```go
import (
	"context"
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

leveler, _ := slogleveler.New(
	slogleveler.WithWriter(os.Stdout),
	slogleveler.WithLevel(traceLevel),
	slogleveler.WithLevel(fatalLevel))

slogger := slog.New(slog.NewTextHandler(leveler,
	&slog.HandlerOptions{
		Level: LogLevelTrace,
	}))

slogger.Log(context.Background(), LogLevelTrace, "DEBUG-1 log 1", "key", "value")
slogger.Log(context.Background(), LogLevelFatal, "ERROR+1 log 1", "key", "value")
```