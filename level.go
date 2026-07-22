package slfg

import "fmt"

// Level — уровень логирования.
type Level int

const (
	TRACE Level = iota
	DEBUG
	INFO
	WARN
	ERROR
	PANIC
	FATAL
)

var levelNames = [...]string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "PANIC", "FATAL"}

func (l Level) String() string {
	if l >= TRACE && l <= FATAL {
		return levelNames[l]
	}
	return fmt.Sprintf("LEVEL(%d)", int(l))
}

// ParseLevel разбирает строковое представление уровня (регистронезависимо).
func ParseLevel(s string) (Level, error) {
	for i, name := range levelNames {
		if equalFold(s, name) {
			return Level(i), nil
		}
	}
	return 0, fmt.Errorf("slfg: неизвестный уровень %q", s)
}

func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		ca, cb := a[i], b[i]
		if ca >= 'a' && ca <= 'z' {
			ca -= 'a' - 'A'
		}
		if cb >= 'a' && cb <= 'z' {
			cb -= 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return true
}
