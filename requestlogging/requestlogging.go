package requestlogging

import (
	"io"
	"net/http"
	"os"

	gorillaHandlers "github.com/gorilla/handlers"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/aaronriekenberg/go-httpd/config"
)

const writeChannelCapacity = 100

type RequestLogger interface {
	WrapHttpHandler(handler http.Handler) http.Handler
}

type requestLogger struct {
	writeChannel chan []byte
}

func (requestLogger *requestLogger) Write(p []byte) (n int, err error) {
	bufferLength := len(p)
	requestLogger.writeChannel <- p
	return bufferLength, nil
}

func (requestLogger *requestLogger) runAsyncWriter(
	writer io.Writer,
) {
	for {
		buffer := <-requestLogger.writeChannel
		writer.Write(buffer)
	}
}

func (requestLogger *requestLogger) WrapHttpHandler(handler http.Handler) http.Handler {
	if requestLogger == nil {
		return handler
	}

	return gorillaHandlers.CombinedLoggingHandler(requestLogger, handler)
}

func NewRequestLogger(
	requestLoggerConfig *config.RequestLogger,
) RequestLogger {

	if requestLoggerConfig == nil {
		return (*requestLogger)(nil)
	}

	var writer io.Writer

	if requestLoggerConfig.LogToStdout {
		writer = os.Stdout
	} else {
		writer = &lumberjack.Logger{
			Filename:   requestLoggerConfig.RequestLogFile,
			MaxSize:    requestLoggerConfig.MaxSizeMegabytes,
			MaxBackups: requestLoggerConfig.MaxBackups,
		}
	}

	requestLogger := &requestLogger{
		writeChannel: make(chan []byte, writeChannelCapacity),
	}

	go requestLogger.runAsyncWriter(writer)

	return requestLogger
}
