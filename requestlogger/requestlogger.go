package requestlogger

import (
	"io"
	"net/http"

	gorillaHandlers "github.com/gorilla/handlers"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/aaronriekenberg/go-httpd/config"
)

type RequestLogger struct {
	writer io.WriteCloser
}

func (RequestLogger *RequestLogger) WrapHttpHandler(handler http.Handler) http.Handler {
	return gorillaHandlers.CombinedLoggingHandler(RequestLogger.writer, handler)
}

func CreateRequestLogger(
	requestLoggerConfig *config.RequestLogger,
) *RequestLogger {

	if requestLoggerConfig == nil {
		return nil
	}

	return &RequestLogger{
		writer: &lumberjack.Logger{
			Filename:   requestLoggerConfig.RequestLogFile,
			MaxSize:    requestLoggerConfig.MaxSizeMegabytes,
			MaxBackups: requestLoggerConfig.MaxBackups,
		},
	}
}
