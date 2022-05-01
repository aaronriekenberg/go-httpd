package logging

import (
	"io"
	"log"
	"os"
)

type VerboseLogger interface {
	// Printf is silent if verbose = false.
	Printf(format string, v ...interface{})

	// Enable or disable verbose logging.
	SetVerboseEnabled(verboseEnabled bool)
}

type FatalLogger interface {
	// Fatalf always calls log.Fatalf to log and exit.
	Fatalf(format string, v ...interface{})
}

type Logger interface {
	VerboseLogger
	FatalLogger
}

type verboseLogger struct {
	*log.Logger
}

func (verboseLogger *verboseLogger) SetVerboseEnabled(verboseEnabled bool) {
	if verboseEnabled {
		verboseLogger.SetOutput(os.Stdout)
	} else {
		verboseLogger.SetOutput(io.Discard)
	}
}

type logger struct {
	VerboseLogger
	FatalLogger
}

var loggerInstance logger

func GetLogger() Logger {
	return &loggerInstance
}

func init() {
	actualVerboseLogger := log.New(os.Stdout, "[verbose] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lmsgprefix)
	actualFatalLogger := log.New(os.Stderr, "[fatal] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lmsgprefix)

	loggerInstance = logger{
		VerboseLogger: &verboseLogger{
			Logger: actualVerboseLogger,
		},
		FatalLogger: actualFatalLogger,
	}
}
