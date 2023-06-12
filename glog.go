package glog

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

const (
	// LevelDebug is the default logging level.
	LevelDebug = 1 << iota
	// LevelInfo is the logging level for informational messages.
	LevelInfo
	// LevelWarn is the logging level for warning messages.
	LevelWarn
	// LevelError is the logging level for error messages.
	LevelError
	// LevelFatal is the logging level for fatal messages.
	LevelFatal

	DefaultLevel = LevelInfo | LevelWarn | LevelError | LevelFatal
)

var levels = map[int]string{
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARNING",
	LevelError: "ERROR",
	LevelFatal: "FATAL",
}

func RegisterLevel(level int, name string) {
	levels[level] = name
}

type Logger struct {
	w      io.Writer
	level  int
	prefix string
}

func NewLogger(w io.Writer, level int) *Logger {
	if level == 0 {
		level = DefaultLevel
	}

	return &Logger{
		w:     w,
		level: level,
	}
}

func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.write(LevelDebug, msg, args...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.write(LevelInfo, msg, args...)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.write(LevelWarn, msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.write(LevelError, msg, args...)
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.write(LevelFatal, msg, args...)
	panic(fmt.Errorf(msg, args...))
}

func (l *Logger) AtLevel(level int, msg string, args ...interface{}) {
	l.write(level, msg, args...)
}

func (l *Logger) write(messageLevel int, msg string, args ...interface{}) error {
	if messageLevel&l.level == 0 {
		return nil
	}
	levelPrefixes := []string{}
	for level, prefix := range levels {
		if l.level&level&messageLevel != 0 {
			levelPrefixes = append(levelPrefixes, prefix)
		}
	}

	sort.Strings(levelPrefixes)

	levelPrefix := strings.Join(levelPrefixes, "|")

	if l.prefix != "" {
		levelPrefix = fmt.Sprintf("%s|%s", levelPrefix, l.prefix)
	}
	_, err := fmt.Fprintf(l.w, "[%s] %s\n", levelPrefix, fmt.Sprintf(msg, args...))
	return err
}
