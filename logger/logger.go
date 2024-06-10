package logger

import (
	"log"
	"os"
)

type LogLevel int

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	FATAL
)

var logLevelStrings = map[LogLevel]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

var (
	logging  log.Logger
	logLevel LogLevel
)

func InitLogger(level LogLevel) {
	logLevel = level
	logging = *log.New(os.Stdout, "", log.Lmicroseconds|log.Llongfile)
}

func Debug(format string, v ...any) {
	if logLevel == DEBUG {
		logging.Printf("[%s] "+format, append([]any{logLevelStrings[logLevel]}, v...)...)
	}
}
