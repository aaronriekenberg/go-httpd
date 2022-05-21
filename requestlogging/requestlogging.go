package requestlogging

import (
	"io"
	"net/http"
	"os"

	gorillaHandlers "github.com/gorilla/handlers"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/aaronriekenberg/go-httpd/config"
)

type RequestLogger struct {
	writer io.Writer
}

func (requestLogger *RequestLogger) WrapHttpHandler(handler http.Handler) http.Handler {
	if requestLogger == nil {
		return handler
	}

	return gorillaHandlers.CombinedLoggingHandler(requestLogger.writer, handler)
}

func NewRequestLogger(
	requestLoggerConfig *config.RequestLogger,
) *RequestLogger {

	if requestLoggerConfig == nil {
		return nil
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

	return &RequestLogger{
		writer: writer,
	}
}
