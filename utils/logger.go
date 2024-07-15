package utils

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	nested "github.com/antonfisher/nested-logrus-formatter"
	_ "github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func InitLogger() {
	Logger.SetLevel(getLoggerLevel(os.Getenv("LOG_LEVEL")))
	Logger.SetReportCaller(true)
	Logger.SetFormatter(&nested.Formatter{
		// HideKeys: true,
		// FieldsOrder:     []string{"category"},
		TimestampFormat: "2006-01-02 15:04:05",
		CallerFirst:     true,
		CustomCallerFormatter: func(f *runtime.Frame) string {
			cwd, _ := os.Getwd()
			return fmt.Sprintf(" [%s:%d]", strings.Replace(f.File, cwd+"/", "", 1), f.Line)
		},
	})
}

func getLoggerLevel(value string) logrus.Level {
	switch value {
	case "DEBUG":
		return logrus.DebugLevel
	case "TRACE":
		return logrus.TraceLevel
	default:
		return logrus.InfoLevel
	}
}

// import (
// 	"log"
// 	"os"
// )
//
// type LogLevel int
//
// const (
// 	DEBUG = iota
// 	INFO
// 	WARN
// 	ERROR
// 	FATAL
// )
//
// var logLevelStrings = map[LogLevel]string{
// 	DEBUG: "DEBUG",
// 	INFO:  "INFO",
// 	WARN:  "WARN",
// 	ERROR: "ERROR",
// 	FATAL: "FATAL",
// }
//
// var logLevelFlag = map[LogLevel]int{
// 	DEBUG: log.LstdFlags | log.Lmicroseconds | log.Llongfile,
// 	INFO:  log.LstdFlags,
// }
//
// var (
// 	logging  = log.New(os.Stdout, "", log.LstdFlags)
// 	logLevel LogLevel
// )
//
// func InitLogger() {
// 	level := os.Getenv("LOG_LEVEL")
// 	if level == "" {
// 		level = "DEBUG"
// 	}
// 	logLevel = getLogLevelFromString(level)
//
// 	flag := logLevelFlag[INFO]
// 	if level == "DEBUG" {
// 		flag = logLevelFlag[DEBUG]
// 	}
// 	logging.SetFlags(flag)
// }
//
// func getLogLevelFromString(value string) LogLevel {
// 	for k, v := range logLevelStrings {
// 		if v == value {
// 			return k
// 		}
// 	}
// 	return LogLevel(0)
// }

// func Debug(format string, v ...any) {
// 	if logLevel == DEBUG {
// 		logging.Printf("[%s] "+format, append([]any{logLevelStrings[logLevel]}, v...)...)
// 	}
// }
//
// func Info(format string, v ...any) {
// 	if logLevel <= INFO {
// 		logging.Printf("[%s] "+format, append([]any{logLevelStrings[logLevel]}, v...)...)
// 	}
// }
//
// func Warn(format string, v ...any) {
// 	if logLevel <= WARN {
// 		logging.Printf("[%s] "+format, append([]any{logLevelStrings[logLevel]}, v...)...)
// 	}
// }
//
// func Error(format string, v ...any) {
// 	if logLevel <= ERROR {
// 		logging.Printf("[%s] "+format, append([]any{logLevelStrings[logLevel]}, v...)...)
// 	}
// }
//
// func Fatal(format string, v ...any) {
// 	logging.Fatalf("[%s] "+format, append([]any{logLevelStrings[logLevel]}, v...)...)
// }
