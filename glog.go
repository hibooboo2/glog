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

func (l *Logger) SetLevel(level uint64) {
	l.level = level
	l.levelPrefixes = map[uint64]string{}
}

func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
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

	levelName, exists := l.levels[messageLevel]

	levelPrefixes := []string{}

	if exists && l.level&messageLevel != 0 {
		levelPrefixes = append(levelPrefixes, levelName)
	}

	if !exists {
		for level, prefix := range l.levels {
			if l.level&level&messageLevel != 0 {
				levelPrefixes = append(levelPrefixes, prefix)
			}
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
