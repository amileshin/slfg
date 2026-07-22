package config

import (
	"os"
	"regexp"
	"strings"
)

var varPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// ResolveVars подставляет переменные окружения в формате ${VAR} и ${VAR:default}.
// Если переменная не найдена и default не задан — выражение остаётся как есть.
func ResolveVars(data []byte) []byte {
	return varPattern.ReplaceAllFunc(data, func(match []byte) []byte {
		inner := string(match[2 : len(match)-1])
		name, defaultVal, hasDefault := strings.Cut(inner, ":")
		if val := os.Getenv(name); val != "" {
			return []byte(val)
		}
		if hasDefault {
			return []byte(defaultVal)
		}
		return match
	})
}
