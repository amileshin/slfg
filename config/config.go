package config

import (
	"fmt"
	"os"
	"strings"

	"slfg"
	"slfg/appender"
	"slfg/filter"
)

// Load читает и разбирает файл конфигурации (logback.xml или logback.yaml).
func Load(path string) (*Configuration, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("slfg/config: %w", err)
	}
	data = ResolveVars(data)

	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		return parseYAML(data)
	}
	return parseXML(data)
}

// Apply применяет конфигурацию: создаёт appender'ы, фильтры и подменяет Provider.
func Apply(cfg *Configuration) error {
	appenders, err := buildAppenders(cfg.Appenders)
	if err != nil {
		return err
	}

	var rootAppenders []slfg.Appender
	for _, ref := range cfg.Root.AppenderRefs {
		a, ok := appenders[ref]
		if !ok {
			return fmt.Errorf("slfg/config: appender %q не найден", ref)
		}
		rootAppenders = append(rootAppenders, a)
	}

	level := slfg.INFO
	if cfg.Root.Level != "" {
		level, err = slfg.ParseLevel(cfg.Root.Level)
		if err != nil {
			return err
		}
	}

	slfg.SetProvider(&configuredProvider{level: level, appenders: rootAppenders})
	return nil
}

// LoadAndWatch загружает конфигурацию и запускает hot reload.
func LoadAndWatch(path string) (*Watcher, error) {
	if err := load(path); err != nil {
		return nil, err
	}
	w, err := NewWatcher(path, func() {
		load(path)
	})
	if err != nil {
		return nil, err
	}
	w.Start()
	return w, nil
}

func load(path string) error {
	cfg, err := Load(path)
	if err != nil {
		return err
	}
	return Apply(cfg)
}

// --- Построение appender'ов и фильтров ---

func buildAppenders(configs []AppenderConfig) (map[string]slfg.Appender, error) {
	result := make(map[string]slfg.Appender, len(configs))
	for _, ac := range configs {
		a, err := buildAppender(ac)
		if err != nil {
			return nil, fmt.Errorf("slfg/config: appender %q: %w", ac.Name, err)
		}
		if len(ac.Filters) > 0 {
			chain, err := buildFilters(ac.Filters)
			if err != nil {
				return nil, fmt.Errorf("slfg/config: appender %q: %w", ac.Name, err)
			}
			if f, ok := a.(slfg.Filterable); ok {
				f.SetFilterChain(chain)
			}
		}
		result[ac.Name] = a
	}
	return result, nil
}

func buildAppender(ac AppenderConfig) (slfg.Appender, error) {
	switch {
	case matchesClass(ac.Class, "ConsoleAppender"):
		return appender.NewConsole(ac.Target, nil), nil
	case matchesClass(ac.Class, "FileAppender"):
		if ac.File == "" {
			return nil, fmt.Errorf("не указан <file>")
		}
		return appender.NewFile(ac.File, ac.Append, nil)
	default:
		return nil, fmt.Errorf("неизвестный класс %q", ac.Class)
	}
}

func buildFilters(configs []FilterConfig) (*slfg.FilterChain, error) {
	chain := slfg.NewFilterChain()
	for _, fc := range configs {
		level, err := slfg.ParseLevel(fc.Level)
		if err != nil {
			return nil, err
		}
		switch {
		case matchesClass(fc.Class, "ThresholdFilter"):
			chain.Add(filter.NewThreshold(level))
		case matchesClass(fc.Class, "LevelFilter"):
			chain.Add(filter.NewLevel(level))
		default:
			return nil, fmt.Errorf("неизвестный класс фильтра %q", fc.Class)
		}
	}
	return chain, nil
}

// matchesClass проверяет совпадение по полному имени класса или по короткому.
func matchesClass(class, short string) bool {
	return class == short || strings.HasSuffix(class, "."+short)
}

// --- configuredProvider ---

type configuredProvider struct {
	level     slfg.Level
	appenders []slfg.Appender
}

func (p *configuredProvider) Logger(name string) slfg.Logger {
	return slfg.NewLogger(name, p.level, p.appenders...)
}
