package glog

import (
	"fmt"
	"os"
)

var DefaultLogger = NewLogger(os.Stdout, DefaultLevel)

func SetLevel(level int) {
	DefaultLogger.level = level
}

func SetPrefix(prefix string) {
	DefaultLogger.prefix = prefix
}

func Print(args ...interface{}) {
	DefaultLogger.Info(args...)
}

func Println(args ...interface{}) {
	DefaultLogger.Infoln(args...)
}

func Printf(msg string, args ...interface{}) {
	DefaultLogger.Infof(msg, args...)
}

func Debug(args ...interface{}) {
	DefaultLogger.write(LevelDebug, args...)
}

func Info(args ...interface{}) {
	DefaultLogger.write(LevelInfo, args...)
}

func Warn(args ...interface{}) {
	DefaultLogger.write(LevelWarn, args...)
}

func Error(args ...interface{}) {
	DefaultLogger.write(LevelError, args...)
}

func Fatal(args ...interface{}) {
	DefaultLogger.write(LevelFatal, args...)
	panic(fmt.Errorf(fmt.Sprint(args...)))
}

func Debugln(args ...interface{}) {
	DefaultLogger.writeln(LevelDebug, args...)
}

func Infoln(args ...interface{}) {
	DefaultLogger.writeln(LevelInfo, args...)
}

func Warnln(args ...interface{}) {
	DefaultLogger.writeln(LevelWarn, args...)
}

func Errorln(args ...interface{}) {
	DefaultLogger.writeln(LevelError, args...)
}

func Fatalln(args ...interface{}) {
	DefaultLogger.writeln(LevelFatal, args...)
	panic(fmt.Errorf(fmt.Sprint(args...)))
}

func Debugf(msg string, args ...interface{}) {
	DefaultLogger.writef(LevelDebug, msg, args...)
}

func Infof(msg string, args ...interface{}) {
	DefaultLogger.writef(LevelInfo, msg, args...)
}

func Warnf(msg string, args ...interface{}) {
	DefaultLogger.writef(LevelWarn, msg, args...)
}

func Errorf(msg string, args ...interface{}) {
	DefaultLogger.writef(LevelError, msg, args...)
}

func Fatalf(msg string, args ...interface{}) {
	DefaultLogger.writef(LevelFatal, msg, args...)
	panic(fmt.Errorf(fmt.Sprint(args...)))
}

func AtLevel(level int, args ...interface{}) {
	DefaultLogger.write(level, args...)
}

func AtLevelln(level int, args ...interface{}) {
	DefaultLogger.writeln(level, args...)
}

func AtLevelf(level int, msg string, args ...interface{}) {
	DefaultLogger.writef(level, msg, args...)
}
