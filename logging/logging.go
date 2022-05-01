package logging

import (
	"log"
	"os"
)

type Logger interface {
	// Printf is silent if verbose = false.
	Printf(format string, v ...interface{})

	// Fatalf always calls log.Fatalf to log and exit.
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

func GetLogger() Logger {
	return &loggerInstance
}

func init() {
	actualLogger := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds)

	loggerInstance = logger{
		actualLogger: actualLogger,
		verbose:      true,
	}
}
