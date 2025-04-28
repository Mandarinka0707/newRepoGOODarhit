// pkg/logger/logger.go
package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}

type logger struct {
	log   *log.Logger
	level LogLevel
}

func NewWithLevel(level string) Logger {
	var lvl LogLevel
	switch strings.ToLower(level) {
	case "debug":
		lvl = LevelDebug
	case "info":
		lvl = LevelInfo
	case "warn", "warning":
		lvl = LevelWarn
	case "error":
		lvl = LevelError
	case "fatal":
		lvl = LevelFatal
	default:
		lvl = LevelInfo
	}

	return &logger{
		log:   log.New(os.Stdout, "[CHAT] ", log.LstdFlags|log.Lshortfile),
		level: lvl,
	}
}

func (l *logger) Debug(v ...interface{}) {
	if l.level <= LevelDebug {
		l.log.Output(2, "[DEBUG] "+fmt.Sprint(v...))
	}
}

func (l *logger) Debugf(format string, v ...interface{}) {
	if l.level <= LevelDebug {
		l.log.Output(2, "[DEBUG] "+fmt.Sprintf(format, v...))
	}
}

func (l *logger) Info(v ...interface{}) {
	if l.level <= LevelInfo {
		l.log.Output(2, "[INFO] "+fmt.Sprint(v...))
	}
}

func (l *logger) Infof(format string, v ...interface{}) {
	if l.level <= LevelInfo {
		l.log.Output(2, "[INFO] "+fmt.Sprintf(format, v...))
	}
}

func (l *logger) Warn(v ...interface{}) {
	if l.level <= LevelWarn {
		l.log.Output(2, "[WARN] "+fmt.Sprint(v...))
	}
}

func (l *logger) Warnf(format string, v ...interface{}) {
	if l.level <= LevelWarn {
		l.log.Output(2, "[WARN] "+fmt.Sprintf(format, v...))
	}
}

func (l *logger) Error(v ...interface{}) {
	if l.level <= LevelError {
		l.log.Output(2, "[ERROR] "+fmt.Sprint(v...))
	}
}

func (l *logger) Errorf(format string, v ...interface{}) {
	if l.level <= LevelError {
		l.log.Output(2, "[ERROR] "+fmt.Sprintf(format, v...))
	}
}

func (l *logger) Fatal(v ...interface{}) {
	l.log.Output(2, "[FATAL] "+fmt.Sprint(v...))
	os.Exit(1)
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	l.log.Output(2, "[FATAL] "+fmt.Sprintf(format, v...))
	os.Exit(1)
}
