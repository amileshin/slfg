package filter

import "slfg"

// LevelFilter — пропускает события строго заданного уровня.
// Совпадение → ACCEPT, несовпадение → DENY.
type LevelFilter struct {
	Level slfg.Level
}

func NewLevel(level slfg.Level) *LevelFilter {
	return &LevelFilter{Level: level}
}

func (f *LevelFilter) Decide(event *slfg.LoggingEvent) slfg.FilterReply {
	if event.Level == f.Level {
		return slfg.ACCEPT
	}
	return slfg.DENY
}
