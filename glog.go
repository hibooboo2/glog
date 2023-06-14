package glog

import (
	"fmt"
	"io"
	"math/bits"
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

var defaultLevels = map[uint64]string{
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARNING",
	LevelError: "ERROR",
	LevelFatal: "FATAL",
}

func (l *Logger) RegisterLevel(level uint64, name string) {
	for registeredLevel, registeredName := range l.levels {
		if name == registeredName {
			panic(fmt.Sprintf("level %q already registered as %b cannot register as %b", registeredName, registeredLevel, level))
		}
		if level == registeredLevel {
			panic(fmt.Sprintf("level %b already registered as %q cannot register as %q", registeredLevel, registeredName, name))
		}
	}
	l.levels[level] = name
	l.levelPrefixes = map[uint64]string{}
}

func (l *Logger) UnRegisterLevel(level uint64) {
	delete(l.levels, level)
	l.levelPrefixes = map[uint64]string{}
}

func (l *Logger) NextLevelShouldRegister() uint64 {
	maxLevel := uint64(0)
	for level := range l.levels {
		if level > maxLevel {
			maxLevel = level
		}
	}
	bitLen := bits.Len64(maxLevel)
	if bitLen < 64 {
		return 1 << bitLen
	}
	panic("too many levels try registering a level with more than one bit set")
}

func (l *Logger) RegisterNextLevel(name string) uint64 {
	level := l.NextLevelShouldRegister()
	l.RegisterLevel(level, name)
	return level
}

type Logger struct {
	w             io.Writer
	level         uint64
	prefix        string
	levels        map[uint64]string
	levelPrefixes map[uint64]string
}

func NewLogger(w io.Writer, level uint64) *Logger {
	if level == 0 {
		level = DefaultLevel
	}

	l := &Logger{
		w:             w,
		level:         level,
		levels:        make(map[uint64]string),
		levelPrefixes: make(map[uint64]string),
	}
	for level, name := range defaultLevels {
		l.RegisterLevel(level, name)
	}
	return l
}

func (l *Logger) SetLevel(level uint64) {
	l.level = level
	l.levelPrefixes = map[uint64]string{}
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

func (l *Logger) AtLevel(level uint64, args ...interface{}) {
	l.write(level, args...)
}

func (l *Logger) AtLevelln(level uint64, args ...interface{}) {
	l.writeln(level, args...)
}

func (l *Logger) AtLevelf(level uint64, msg string, args ...interface{}) {
	l.writef(level, msg, args...)
}

func (l *Logger) CustomLogAtLevel(level uint64) func(args ...interface{}) {
	return func(args ...interface{}) {
		l.AtLevel(level, args...)
	}
}

func (l *Logger) CustomLogAtLevelln(level uint64) func(args ...interface{}) {
	return func(args ...interface{}) {
		l.AtLevelln(level, args...)
	}
}

func (l *Logger) CustomLogAtLevelf(level uint64) func(msg string, args ...interface{}) {
	return func(msg string, args ...interface{}) {
		l.AtLevelf(level, msg, args...)
	}
}

func (l *Logger) writeln(level uint64, args ...interface{}) error {
	if level&l.level == 0 {
		return nil
	}

	levelPrefix := l.getFullPrefix(level)
	if levelPrefix == "" {
		return nil
	}

	_, err := fmt.Fprintf(l.w, "%s%s\n", levelPrefix, fmt.Sprintln(args...))
	return err
}

func (l *Logger) writef(level uint64, msg string, args ...interface{}) error {
	if level&l.level == 0 {
		return nil
	}

	levelPrefix := l.getFullPrefix(level)
	if levelPrefix == "" {
		return nil
	}

	_, err := fmt.Fprintf(l.w, "%s%s\n", levelPrefix, fmt.Sprintf(msg, args...))
	return err
}

func (l *Logger) write(level uint64, args ...interface{}) error {
	if level&l.level == 0 {
		return nil
	}

	levelPrefix := l.getFullPrefix(level)
	if levelPrefix == "" {
		return nil
	}

	_, err := fmt.Fprintf(l.w, "%s%s\n", levelPrefix, fmt.Sprint(args...))
	return err
}

func (l *Logger) getFullPrefix(level uint64) string {
	levelPrefix := l.getLevelPrefix(level)
	if levelPrefix == "" {
		return ""
	}

	if l.prefix != "" {
		levelPrefix = fmt.Sprintf("%s[%s]", levelPrefix, l.prefix)
	}
	if levelPrefix != "" {
		levelPrefix = fmt.Sprintf("%s ", levelPrefix)
	}
	return levelPrefix
}

func (l *Logger) getLevelPrefix(messageLevel uint64) string {
	levelPrefix, exists := l.levelPrefixes[messageLevel]
	if exists {
		return levelPrefix
	}

	levelPrefixes := []string{}
	for level, prefix := range l.levels {
		if l.level&level&messageLevel != 0 {
			levelPrefixes = append(levelPrefixes, prefix)
		}
	}

	if len(levelPrefixes) == 0 {
		l.levelPrefixes[messageLevel] = ""
		return ""
	}

	sort.Strings(levelPrefixes)

	levelPrefix = fmt.Sprintf("[%s]", strings.Join(levelPrefixes, "|"))
	l.levelPrefixes[messageLevel] = levelPrefix
	return levelPrefix
}
