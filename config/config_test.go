package config

import (
	"os"
	"path/filepath"
	"testing"

	"slfg"
)

func TestResolveVars(t *testing.T) {
	t.Setenv("SLFG_TEST_DIR", "/tmp/logs")

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"подстановка", "${SLFG_TEST_DIR}/app.log", "/tmp/logs/app.log"},
		{"default", "${SLFG_MISSING:/var/log}/app.log", "/var/log/app.log"},
		{"нет переменной", "${SLFG_MISSING}", "${SLFG_MISSING}"},
		{"без переменных", "plain text", "plain text"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(ResolveVars([]byte(tt.input)))
			if got != tt.want {
				t.Errorf("ResolveVars(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseXML(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<configuration>
    <appender name="STDOUT" class="ch.qos.logback.core.ConsoleAppender">
        <target>System.out</target>
        <filter class="ch.qos.logback.classic.filter.ThresholdFilter">
            <level>INFO</level>
        </filter>
    </appender>
    <appender name="FILE" class="ch.qos.logback.core.FileAppender">
        <file>app.log</file>
        <append>true</append>
    </appender>
    <root level="DEBUG">
        <appender-ref ref="STDOUT"/>
        <appender-ref ref="FILE"/>
    </root>
</configuration>`

	cfg, err := parseXML([]byte(xml))
	if err != nil {
		t.Fatalf("parseXML: %v", err)
	}

	if len(cfg.Appenders) != 2 {
		t.Fatalf("ожидалось 2 appender'а, получено %d", len(cfg.Appenders))
	}

	stdout := cfg.Appenders[0]
	if stdout.Name != "STDOUT" {
		t.Errorf("имя = %q, want STDOUT", stdout.Name)
	}
	if stdout.Target != "System.out" {
		t.Errorf("target = %q, want System.out", stdout.Target)
	}
	if len(stdout.Filters) != 1 || stdout.Filters[0].Level != "INFO" {
		t.Error("ожидался ThresholdFilter с уровнем INFO")
	}

	file := cfg.Appenders[1]
	if file.File != "app.log" {
		t.Errorf("file = %q, want app.log", file.File)
	}
	if !file.Append {
		t.Error("append должен быть true")
	}

	if cfg.Root.Level != "DEBUG" {
		t.Errorf("root level = %q, want DEBUG", cfg.Root.Level)
	}
	if len(cfg.Root.AppenderRefs) != 2 {
		t.Fatalf("ожидалось 2 appender-ref, получено %d", len(cfg.Root.AppenderRefs))
	}
}

func TestParseYAML(t *testing.T) {
	yamlData := `
configuration:
  appender:
    - name: STDOUT
      class: ch.qos.logback.core.ConsoleAppender
      target: System.out
      filter:
        - class: ch.qos.logback.classic.filter.ThresholdFilter
          level: INFO
    - name: FILE
      class: ch.qos.logback.core.FileAppender
      file: app.log
  root:
    level: DEBUG
    appender-ref:
      - ref: STDOUT
      - ref: FILE
`

	cfg, err := parseYAML([]byte(yamlData))
	if err != nil {
		t.Fatalf("parseYAML: %v", err)
	}

	if len(cfg.Appenders) != 2 {
		t.Fatalf("ожидалось 2 appender'а, получено %d", len(cfg.Appenders))
	}
	if cfg.Appenders[0].Name != "STDOUT" {
		t.Errorf("имя = %q, want STDOUT", cfg.Appenders[0].Name)
	}
	if cfg.Root.Level != "DEBUG" {
		t.Errorf("root level = %q, want DEBUG", cfg.Root.Level)
	}
	if len(cfg.Root.AppenderRefs) != 2 {
		t.Fatalf("ожидалось 2 appender-ref, получено %d", len(cfg.Root.AppenderRefs))
	}
}

func TestLoadXML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "logback.xml")

	t.Setenv("SLFG_LOG_FILE", filepath.Join(dir, "app.log"))

	content := `<configuration>
    <appender name="FILE" class="ch.qos.logback.core.FileAppender">
        <file>${SLFG_LOG_FILE}</file>
    </appender>
    <root level="INFO">
        <appender-ref ref="FILE"/>
    </root>
</configuration>`
	os.WriteFile(path, []byte(content), 0o644)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	want := filepath.Join(dir, "app.log")
	if cfg.Appenders[0].File != want {
		t.Errorf("file = %q, want %q", cfg.Appenders[0].File, want)
	}
}

func TestLoadYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "logback.yaml")

	content := `
configuration:
  appender:
    - name: CONSOLE
      class: ConsoleAppender
      target: stdout
  root:
    level: WARN
    appender-ref:
      - ref: CONSOLE
`
	os.WriteFile(path, []byte(content), 0o644)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Root.Level != "WARN" {
		t.Errorf("root level = %q, want WARN", cfg.Root.Level)
	}
}

func TestApply(t *testing.T) {
	dir := t.TempDir()
	logFile := filepath.Join(dir, "test.log")

	cfg := &Configuration{
		Appenders: []AppenderConfig{
			{Name: "FILE", Class: "FileAppender", File: logFile, Append: true},
		},
		Root: RootConfig{
			Level:        "INFO",
			AppenderRefs: []string{"FILE"},
		},
	}

	if err := Apply(cfg); err != nil {
		t.Fatalf("Apply: %v", err)
	}

	// Проверяем, что логгер пишет в файл
	logger := slfg.GetLogger("test")
	logger.Info("запись через Apply")

	data, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if len(data) == 0 {
		t.Error("файл пуст после записи")
	}
}

func TestApplyUnknownAppenderRef(t *testing.T) {
	cfg := &Configuration{
		Root: RootConfig{
			Level:        "INFO",
			AppenderRefs: []string{"NONEXISTENT"},
		},
	}
	if err := Apply(cfg); err == nil {
		t.Error("ожидалась ошибка для несуществующего appender-ref")
	}
}

func TestApplyUnknownAppenderClass(t *testing.T) {
	cfg := &Configuration{
		Appenders: []AppenderConfig{
			{Name: "X", Class: "UnknownAppender"},
		},
		Root: RootConfig{
			Level:        "INFO",
			AppenderRefs: []string{"X"},
		},
	}
	if err := Apply(cfg); err == nil {
		t.Error("ожидалась ошибка для неизвестного класса")
	}
}

func TestMatchesClass(t *testing.T) {
	tests := []struct {
		class string
		short string
		want  bool
	}{
		{"ch.qos.logback.core.ConsoleAppender", "ConsoleAppender", true},
		{"ConsoleAppender", "ConsoleAppender", true},
		{"ch.qos.logback.core.FileAppender", "ConsoleAppender", false},
	}
	for _, tt := range tests {
		if got := matchesClass(tt.class, tt.short); got != tt.want {
			t.Errorf("matchesClass(%q, %q) = %v, want %v", tt.class, tt.short, got, tt.want)
		}
	}
}
