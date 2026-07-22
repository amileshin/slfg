package appender

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"slfg"
)

func testEvent(level slfg.Level, msg string) *slfg.LoggingEvent {
	return &slfg.LoggingEvent{
		Level:      level,
		Message:    msg,
		LoggerName: "test",
		Time:       time.Date(2026, 1, 15, 10, 30, 0, 0, time.UTC),
	}
}

func TestConsoleAppenderStdout(t *testing.T) {
	a := NewConsole("System.out", nil)
	if a.writer != os.Stdout {
		t.Error("System.out должен направлять в stdout")
	}
}

func TestConsoleAppenderStderr(t *testing.T) {
	a := NewConsole("System.err", nil)
	if a.writer != os.Stderr {
		t.Error("System.err должен направлять в stderr")
	}

	a2 := NewConsole("", nil)
	if a2.writer != os.Stderr {
		t.Error("по умолчанию должен быть stderr")
	}
}

func TestConsoleAppenderOutput(t *testing.T) {
	var buf bytes.Buffer
	a := NewConsole("stdout", nil)
	a.writer = &buf

	a.Append(testEvent(slfg.INFO, "привет"))
	output := buf.String()
	if !strings.Contains(output, "[INFO]") {
		t.Errorf("ожидался [INFO], получено: %s", output)
	}
	if !strings.Contains(output, "привет") {
		t.Errorf("ожидалось сообщение, получено: %s", output)
	}
}

func TestConsoleAppenderWithFilter(t *testing.T) {
	var buf bytes.Buffer
	a := NewConsole("stdout", nil)
	a.writer = &buf

	chain := slfg.NewFilterChain(&denyAll{})
	a.SetFilterChain(chain)

	a.Append(testEvent(slfg.INFO, "не должно пройти"))
	if buf.Len() != 0 {
		t.Error("фильтр DENY должен блокировать запись")
	}
}

func TestFileAppender(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.log")

	a, err := NewFile(path, true, nil)
	if err != nil {
		t.Fatalf("NewFile: %v", err)
	}

	a.Append(testEvent(slfg.WARN, "запись в файл"))
	a.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if !strings.Contains(string(data), "запись в файл") {
		t.Errorf("ожидалось сообщение в файле, получено: %s", string(data))
	}
}

func TestFileAppenderAppendMode(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.log")

	a1, _ := NewFile(path, true, nil)
	a1.Append(testEvent(slfg.INFO, "первая"))
	a1.Close()

	a2, _ := NewFile(path, true, nil)
	a2.Append(testEvent(slfg.INFO, "вторая"))
	a2.Close()

	data, _ := os.ReadFile(path)
	content := string(data)
	if !strings.Contains(content, "первая") || !strings.Contains(content, "вторая") {
		t.Error("режим дозаписи должен сохранять обе записи")
	}
}

func TestFileAppenderTruncateMode(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.log")

	a1, _ := NewFile(path, true, nil)
	a1.Append(testEvent(slfg.INFO, "старая"))
	a1.Close()

	a2, _ := NewFile(path, false, nil)
	a2.Append(testEvent(slfg.INFO, "новая"))
	a2.Close()

	data, _ := os.ReadFile(path)
	content := string(data)
	if strings.Contains(content, "старая") {
		t.Error("режим перезаписи должен удалить старые записи")
	}
	if !strings.Contains(content, "новая") {
		t.Error("новая запись должна присутствовать")
	}
}

type denyAll struct{}

func (f *denyAll) Decide(*slfg.LoggingEvent) slfg.FilterReply { return slfg.DENY }
