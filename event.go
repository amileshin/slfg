package slfg

import (
	"fmt"
	"time"
)

// LoggingEvent — событие логирования, передаваемое в appender'ы.
type LoggingEvent struct {
	Level      Level
	Message    string
	LoggerName string
	Time       time.Time
	MDC        map[string]string
}

// Formatter форматирует событие логирования в байты для записи.
type Formatter func(event *LoggingEvent) []byte

// DefaultFormatter — формат по умолчанию, аналогичный Logback:
// 2006-01-02 15:04:05.000 [INFO] logger.name - сообщение
func DefaultFormatter(event *LoggingEvent) []byte {
	return []byte(fmt.Sprintf("%s [%s] %s - %s\n",
		event.Time.Format("2006-01-02 15:04:05.000"),
		event.Level,
		event.LoggerName,
		event.Message,
	))
}
