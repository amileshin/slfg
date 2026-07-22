package slfg

import (
	"bytes"
	"strings"
	"testing"
)

func TestDefaultLoggerLevels(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test", WARN, NewWriterAppender(&buf, nil))

	logger.Debug("не должно попасть")
	if buf.Len() != 0 {
		t.Error("Debug не должен писаться при уровне WARN")
	}

	logger.Warn("предупреждение")
	if !strings.Contains(buf.String(), "предупреждение") {
		t.Error("Warn должен писаться при уровне WARN")
	}
	if !strings.Contains(buf.String(), "[WARN]") {
		t.Error("ожидался уровень WARN в выводе")
	}
	if !strings.Contains(buf.String(), "test") {
		t.Error("ожидалось имя логгера в выводе")
	}
}

func TestDefaultLoggerIsEnabled(t *testing.T) {
	logger := NewLogger("test", INFO)

	if logger.IsEnabled(DEBUG) {
		t.Error("DEBUG не должен быть включён при уровне INFO")
	}
	if !logger.IsEnabled(INFO) {
		t.Error("INFO должен быть включён при уровне INFO")
	}
	if !logger.IsEnabled(ERROR) {
		t.Error("ERROR должен быть включён при уровне INFO")
	}
}

func TestDefaultLoggerArgs(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test", TRACE, NewWriterAppender(&buf, nil))

	logger.Info("пользователь", "ivan", 42)
	output := buf.String()
	if !strings.Contains(output, "пользователь ivan 42") {
		t.Errorf("ожидалось 'пользователь ivan 42', получено: %s", output)
	}
}

func TestSetProvider(t *testing.T) {
	var buf bytes.Buffer
	p := &testProvider{logger: NewLogger("custom", DEBUG, NewWriterAppender(&buf, nil))}

	old := provider
	SetProvider(p)
	defer SetProvider(old)

	l := GetLogger("anything")
	l.Debug("тест")
	if !strings.Contains(buf.String(), "тест") {
		t.Error("подменённый provider не работает")
	}
}

type testProvider struct {
	logger Logger
}

func (p *testProvider) Logger(string) Logger { return p.logger }

func TestPanicLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("test", TRACE, NewWriterAppender(&buf, nil))

	defer func() {
		r := recover()
		if r == nil {
			t.Error("Panic должен вызывать panic")
		}
		if !strings.Contains(buf.String(), "упс") {
			t.Error("сообщение должно быть залогировано до panic")
		}
	}()

	logger.Panic("упс")
}

func TestPackageLevelFunctions(t *testing.T) {
	var buf bytes.Buffer
	p := &testProvider{logger: NewLogger("root", TRACE, NewWriterAppender(&buf, nil))}

	old := provider
	SetProvider(p)
	defer SetProvider(old)

	Info("пакетная функция")
	if !strings.Contains(buf.String(), "пакетная функция") {
		t.Error("пакетная функция Info не работает")
	}
}
