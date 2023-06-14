package glog

import "fmt"

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
