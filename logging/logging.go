package logging

import (
	"io"
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
	printfLogger *log.Logger
	fatalfLogger *log.Logger
}

func (logger *logger) Printf(format string, v ...interface{}) {
	logger.printfLogger.Printf(format, v...)
}

func (logger *logger) Fatalf(format string, v ...interface{}) {
	logger.fatalfLogger.Fatalf(format, v...)
}

var loggerInstance logger

func SetVerbose(verbose bool) {
	if verbose {
		loggerInstance.printfLogger.SetOutput(os.Stdout)
	} else {
		loggerInstance.printfLogger.SetOutput(io.Discard)
	}
}

func GetLogger() Logger {
	return &loggerInstance
}

func init() {
	printfLogger := log.New(os.Stdout, "[debug] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lmsgprefix)
	fatalfLogger := log.New(os.Stderr, "[fatal] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lmsgprefix)

	loggerInstance = logger{
		printfLogger: printfLogger,
		fatalfLogger: fatalfLogger,
	}
}
