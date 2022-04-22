package requestlogger

import (
	"io"

	"github.com/aaronriekenberg/go-httpd/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

type RequestLogger struct {
	Writer io.WriteCloser
}

func CreateRequestLogger(
	requestLoggerConfig *config.RequestLogger,
) *RequestLogger {

	if requestLoggerConfig == nil {
		return nil
	}

	return &RequestLogger{
		Writer: &lumberjack.Logger{
			Filename:   requestLoggerConfig.RequestLogFile,
			MaxSize:    requestLoggerConfig.MaxSizeMegabytes,
			MaxBackups: requestLoggerConfig.MaxBackups,
		},
	}
}
