package logging

import (
	"log"
	"os"
)

type LoggerInterface interface {
	Printf(format string, v ...interface{})

	Fatalf(format string, v ...interface{})
}

type logger struct {
	actualLogger *log.Logger
	verbose      bool
}

func (logger *logger) Printf(format string, v ...interface{}) {
	if logger.verbose {
		logger.actualLogger.Printf(format, v...)
	}
}

func (logger *logger) Fatalf(format string, v ...interface{}) {
	logger.actualLogger.Fatalf(format, v...)
}

var loggerInstance logger

func SetVerbose(verbose bool) {
	loggerInstance.verbose = verbose
}

func GetLogger() LoggerInterface {
	return &loggerInstance
}

func init() {
	actualLogger := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds)

	loggerInstance = logger{
		actualLogger: actualLogger,
		verbose:      true,
	}
}
