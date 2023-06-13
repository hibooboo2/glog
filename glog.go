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
	for registeredLevel, registeredName := range levels {
		if name == registeredName {
			panic(fmt.Sprintf("level %q already registered as %b cannot register as %b", registeredName, registeredLevel, level))
		}
		if level == registeredLevel {
			panic(fmt.Sprintf("level %b already registered as %q cannot register as %q", registeredLevel, registeredName, name))
		}
	}
	levels[level] = name
}

func UnRegisterLevel(level int) {
	delete(levels, level)
}

func MaxRegisteredLevel() int {
	maxLevel := 0
	for level := range levels {
		if level > maxLevel {
			maxLevel = level
		}
	}
	return maxLevel
}

func RegisterNextLevel(name string) int {
	level := MaxRegisteredLevel() >> 1
	RegisterLevel(level, name)
	return level
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

func (l *Logger) SetLevel(level int) {
	l.level = level
}

func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

func (l *Logger) Print(args ...interface{}) {
	l.Info(args...)
}

func (l *Logger) Println(args ...interface{}) {
	l.Infoln(args...)
}

func (l *Logger) Printf(msg string, args ...interface{}) {
	l.Infof(msg, args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.write(LevelDebug, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.write(LevelInfo, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.write(LevelWarn, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.write(LevelError, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.write(LevelFatal, args...)
	panic(fmt.Errorf(fmt.Sprint(args...)))
}

func (l *Logger) Debugln(args ...interface{}) {
	l.writeln(LevelDebug, args...)
}

func (l *Logger) Infoln(args ...interface{}) {
	l.writeln(LevelInfo, args...)
}

func (l *Logger) Warnln(args ...interface{}) {
	l.writeln(LevelWarn, args...)
}

func (l *Logger) Errorln(args ...interface{}) {
	l.writeln(LevelError, args...)
}

func (l *Logger) Fatalln(args ...interface{}) {
	l.writeln(LevelFatal, args...)
	panic(fmt.Errorf(fmt.Sprint(args...)))
}

func (l *Logger) Debugf(msg string, args ...interface{}) {
	l.writef(LevelDebug, msg, args...)
}

func (l *Logger) Infof(msg string, args ...interface{}) {
	l.writef(LevelInfo, msg, args...)
}

func (l *Logger) Warnf(msg string, args ...interface{}) {
	l.writef(LevelWarn, msg, args...)
}

func (l *Logger) Errorf(msg string, args ...interface{}) {
	l.writef(LevelError, msg, args...)
}

func (l *Logger) Fatalf(msg string, args ...interface{}) {
	l.writef(LevelFatal, msg, args...)
	panic(fmt.Errorf(fmt.Sprint(args...)))
}

func (l *Logger) AtLevel(level int, args ...interface{}) {
	l.write(level, args...)
}

func (l *Logger) AtLevelln(level int, args ...interface{}) {
	l.writeln(level, args...)
}

func (l *Logger) AtLevelf(level int, msg string, args ...interface{}) {
	l.writef(level, msg, args...)
}

func (l *Logger) writeln(messageLevel int, args ...interface{}) error {
	if messageLevel&l.level == 0 {
		return nil
	}
	levelPrefixes := []string{}
	for level, prefix := range levels {
		if l.level&level&messageLevel != 0 {
			levelPrefixes = append(levelPrefixes, prefix)
		}
	}

	if len(levelPrefixes) == 0 {
		return nil
	}

	sort.Strings(levelPrefixes)

	levelPrefix := strings.Join(levelPrefixes, "|")

	if l.prefix != "" {
		levelPrefix = fmt.Sprintf("%s|%s", levelPrefix, l.prefix)
	}
	_, err := fmt.Fprintf(l.w, "[%s] %s\n", levelPrefix, fmt.Sprintln(args...))
	return err
}

func (l *Logger) writef(messageLevel int, msg string, args ...interface{}) error {
	if messageLevel&l.level == 0 {
		return nil
	}
	levelPrefixes := []string{}
	for level, prefix := range levels {
		if l.level&level&messageLevel != 0 {
			levelPrefixes = append(levelPrefixes, prefix)
		}
	}

	if len(levelPrefixes) == 0 {
		return nil
	}

	sort.Strings(levelPrefixes)

	levelPrefix := strings.Join(levelPrefixes, "|")

	if l.prefix != "" {
		levelPrefix = fmt.Sprintf("%s|%s", levelPrefix, l.prefix)
	}
	_, err := fmt.Fprintf(l.w, "[%s] %s\n", levelPrefix, fmt.Sprintf(msg, args...))
	return err
}

func (l *Logger) write(messageLevel int, args ...interface{}) error {
	if messageLevel&l.level == 0 {
		return nil
	}
	levelPrefixes := []string{}
	for level, prefix := range levels {
		if l.level&level&messageLevel != 0 {
			levelPrefixes = append(levelPrefixes, prefix)
		}
	}

	if len(levelPrefixes) == 0 {
		return nil
	}

	sort.Strings(levelPrefixes)

	levelPrefix := strings.Join(levelPrefixes, "|")

	if l.prefix != "" {
		levelPrefix = fmt.Sprintf("%s|%s", levelPrefix, l.prefix)
	}
	_, err := fmt.Fprintf(l.w, "[%s] %s\n", levelPrefix, fmt.Sprint(args...))
	return err
}
