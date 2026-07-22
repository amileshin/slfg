package slfg

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// nopLogger — логгер-заглушка, ничего не пишет.
type nopLogger struct{}

func (nopLogger) Trace(string, ...any) {}
func (nopLogger) Debug(string, ...any) {}
func (nopLogger) Info(string, ...any)  {}
func (nopLogger) Warn(string, ...any)  {}
func (nopLogger) Error(string, ...any) {}
func (nopLogger) Panic(string, ...any) {}
func (nopLogger) Fatal(string, ...any) {}
func (nopLogger) IsEnabled(Level) bool { return false }

// Logger — фасад логирования. Реализация может быть подменена через SetProvider.
type Logger interface {
	Trace(msg string, args ...any)
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Panic(msg string, args ...any)
	Fatal(msg string, args ...any)
	IsEnabled(level Level) bool
}

// Provider создаёт экземпляры Logger по имени.
type Provider interface {
	Logger(name string) Logger
}

var (
	providerMu sync.RWMutex
	provider   Provider = &defaultProvider{}
)

// SetProvider подменяет реализацию логгера (фасад).
func SetProvider(p Provider) {
	providerMu.Lock()
	defer providerMu.Unlock()
	provider = p
}

// GetLogger возвращает логгер с указанным именем через текущий Provider.
func GetLogger(name string) Logger {
	providerMu.RLock()
	defer providerMu.RUnlock()
	return provider.Logger(name)
}

// --- defaultProvider ---

type defaultProvider struct{}

func (p *defaultProvider) Logger(name string) Logger {
	return NewLogger(name, INFO, NewWriterAppender(os.Stderr, nil))
}

// --- defaultLogger ---

type defaultLogger struct {
	name      string
	level     Level
	appenders []Appender
}

// NewLogger создаёт логгер с заданным именем, уровнем и appender'ами.
func NewLogger(name string, level Level, appenders ...Appender) Logger {
	return &defaultLogger{
		name:      name,
		level:     level,
		appenders: appenders,
	}
}

func (l *defaultLogger) Trace(msg string, args ...any) { l.log(TRACE, msg, args...) }
func (l *defaultLogger) Debug(msg string, args ...any) { l.log(DEBUG, msg, args...) }
func (l *defaultLogger) Info(msg string, args ...any)  { l.log(INFO, msg, args...) }
func (l *defaultLogger) Warn(msg string, args ...any)  { l.log(WARN, msg, args...) }
func (l *defaultLogger) Error(msg string, args ...any) { l.log(ERROR, msg, args...) }

func (l *defaultLogger) Panic(msg string, args ...any) {
	l.log(PANIC, msg, args...)
	panic(formatMessage(msg, args...))
}

func (l *defaultLogger) Fatal(msg string, args ...any) {
	l.log(FATAL, msg, args...)
	os.Exit(1)
}

func (l *defaultLogger) IsEnabled(level Level) bool {
	return level >= l.level
}

func (l *defaultLogger) log(level Level, msg string, args ...any) {
	if level < l.level {
		return
	}
	event := &LoggingEvent{
		Level:      level,
		Message:    formatMessage(msg, args...),
		LoggerName: l.name,
		Time:       time.Now(),
	}
	for _, a := range l.appenders {
		a.Append(event)
	}
}

func formatMessage(msg string, args ...any) string {
	if len(args) == 0 {
		return msg
	}
	parts := make([]string, len(args))
	for i, a := range args {
		parts[i] = fmt.Sprint(a)
	}
	return msg + " " + strings.Join(parts, " ")
}

// --- Пакетные функции (делегирование к root-логгеру) ---

func Trace(msg string, args ...any) { GetLogger("root").Trace(msg, args...) }
func Debug(msg string, args ...any) { GetLogger("root").Debug(msg, args...) }
func Info(msg string, args ...any)  { GetLogger("root").Info(msg, args...) }
func Warn(msg string, args ...any)  { GetLogger("root").Warn(msg, args...) }
func Error(msg string, args ...any) { GetLogger("root").Error(msg, args...) }
