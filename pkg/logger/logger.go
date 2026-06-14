package logger

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	level  string
	debug  *log.Logger
	info   *log.Logger
	warn   *log.Logger
	errors *log.Logger
}

func New(level string) *Logger {
	flags := log.Ldate | log.Ltime | log.LUTC
	return &Logger{
		level:  level,
		debug:  log.New(os.Stdout, "[DEBUG] ", flags),
		info:   log.New(os.Stdout, "[INFO]  ", flags),
		warn:   log.New(os.Stderr, "[WARN]  ", flags),
		errors: log.New(os.Stderr, "[ERROR] ", flags),
	}
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.level == "debug" {
		l.debug.Output(2, fmt.Sprintf(format, args...))
	}
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.info.Output(2, fmt.Sprintf(format, args...))
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.warn.Output(2, fmt.Sprintf(format, args...))
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.errors.Output(2, fmt.Sprintf(format, args...))
	os.Exit(1)
}

func (l *Logger) Info(msg string)  { l.info.Output(2, msg) }
func (l *Logger) Warn(msg string)  { l.warn.Output(2, msg) }
func (l *Logger) Fatal(msg string) { l.errors.Output(2, msg); os.Exit(1) }
