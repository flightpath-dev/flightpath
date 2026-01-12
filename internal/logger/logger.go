package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

// LogLevel represents the logging level
type LogLevel int

// LogFormat represents the logging output format
type LogFormat int

const (
	// LogFormatText uses short human-readable text format (default)
	LogFormatText LogFormat = iota
	// LogFormatJSON uses structured JSON format
	LogFormatJSON
)

const (
	// LogLevelDebug enables debug and above log levels
	LogLevelDebug LogLevel = iota
	// LogLevelInfo enables info and above log levels (default)
	LogLevelInfo
	// LogLevelWarn enables warn and above log levels
	LogLevelWarn
	// LogLevelError enables error log level only
	LogLevelError
)

// toSlogLevel converts LogLevel to slog.Level
func (l LogLevel) toSlogLevel() slog.Level {
	switch l {
	case LogLevelDebug:
		return slog.LevelDebug
	case LogLevelInfo:
		return slog.LevelInfo
	case LogLevelWarn:
		return slog.LevelWarn
	case LogLevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// Logger wraps slog.Logger for consistent usage
type Logger struct {
	*slog.Logger
}

// shortTextHandler is a custom handler that produces short, human-readable log lines
// Format: LEVEL [component] message [key=value ...]
type shortTextHandler struct {
	w         io.Writer
	opts      slog.HandlerOptions
	level     slog.Level
	attrs     []slog.Attr
	groupName string
}

func newShortTextHandler(w io.Writer, opts *slog.HandlerOptions) *shortTextHandler {
	return &shortTextHandler{
		w:     w,
		opts:  *opts,
		level: opts.Level.Level(),
		attrs: []slog.Attr{},
	}
}

func (h *shortTextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *shortTextHandler) Handle(ctx context.Context, r slog.Record) error {
	var parts []string

	// Add level
	levelStr := strings.ToUpper(r.Level.String())
	parts = append(parts, levelStr)

	// Collect all attributes (from handler and record)
	allAttrs := make([]slog.Attr, 0, len(h.attrs)+10)
	allAttrs = append(allAttrs, h.attrs...)
	r.Attrs(func(a slog.Attr) bool {
		allAttrs = append(allAttrs, a)
		return true
	})

	// Extract component from attributes
	var component string
	var otherAttrs []string

	for _, a := range allAttrs {
		if a.Key == "component" {
			if a.Value.Kind() == slog.KindString {
				component = a.Value.String()
			}
		} else {
			// Format other attributes as key=value
			attrStr := formatAttr(a)
			if attrStr != "" {
				otherAttrs = append(otherAttrs, attrStr)
			}
		}
	}

	// Add component in brackets if present
	if component != "" {
		parts = append(parts, fmt.Sprintf("[%s]", component))
	}

	// Add message
	if r.Message != "" {
		parts = append(parts, r.Message)
	}

	// Add other attributes
	parts = append(parts, otherAttrs...)

	// Join and write
	line := strings.Join(parts, " ") + "\n"
	_, err := h.w.Write([]byte(line))
	return err
}

func formatAttr(a slog.Attr) string {
	switch a.Value.Kind() {
	case slog.KindString:
		return fmt.Sprintf("%s=%q", a.Key, a.Value.String())
	case slog.KindInt64:
		return fmt.Sprintf("%s=%d", a.Key, a.Value.Int64())
	case slog.KindUint64:
		return fmt.Sprintf("%s=%d", a.Key, a.Value.Uint64())
	case slog.KindFloat64:
		return fmt.Sprintf("%s=%g", a.Key, a.Value.Float64())
	case slog.KindBool:
		return fmt.Sprintf("%s=%t", a.Key, a.Value.Bool())
	case slog.KindTime:
		return fmt.Sprintf("%s=%s", a.Key, a.Value.Time().Format("2006-01-02T15:04:05Z07:00"))
	case slog.KindAny:
		return fmt.Sprintf("%s=%v", a.Key, a.Value.Any())
	default:
		return fmt.Sprintf("%s=%v", a.Key, a.Value)
	}
}

func (h *shortTextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Create a new handler with additional attributes
	newHandler := &shortTextHandler{
		w:         h.w,
		opts:      h.opts,
		level:     h.level,
		attrs:     append(h.attrs, attrs...),
		groupName: h.groupName,
	}
	return newHandler
}

func (h *shortTextHandler) WithGroup(name string) slog.Handler {
	// Groups are not used in our short format, but we need to track the name
	newHandler := &shortTextHandler{
		w:         h.w,
		opts:      h.opts,
		level:     h.level,
		attrs:     h.attrs,
		groupName: name,
	}
	return newHandler
}

// New creates a new structured logger with the specified level and format.
// Log level and format are required parameters.
func New(level LogLevel, format LogFormat) *Logger {
	opts := &slog.HandlerOptions{
		Level: level.toSlogLevel(),
	}

	var handler slog.Handler
	switch format {
	case LogFormatJSON:
		// Use slog's built-in JSON handler
		handler = slog.NewJSONHandler(os.Stderr, opts)
	case LogFormatText:
		fallthrough
	default:
		// Use custom short text handler
		handler = newShortTextHandler(os.Stderr, opts)
	}

	return &Logger{
		Logger: slog.New(handler),
	}
}

// WithPrefix creates a new logger with a prefix attribute
func (l *Logger) WithPrefix(prefix string) *Logger {
	return &Logger{
		Logger: l.Logger.With(slog.String("component", prefix)),
	}
}
