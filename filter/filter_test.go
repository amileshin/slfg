package filter

import (
	"testing"
	"time"

	"slfg"
)

func event(level slfg.Level) *slfg.LoggingEvent {
	return &slfg.LoggingEvent{Level: level, Time: time.Now()}
}

func TestThresholdFilter(t *testing.T) {
	f := NewThreshold(slfg.WARN)

	tests := []struct {
		level slfg.Level
		want  slfg.FilterReply
	}{
		{slfg.TRACE, slfg.DENY},
		{slfg.DEBUG, slfg.DENY},
		{slfg.INFO, slfg.DENY},
		{slfg.WARN, slfg.NEUTRAL},
		{slfg.ERROR, slfg.NEUTRAL},
		{slfg.FATAL, slfg.NEUTRAL},
	}
	for _, tt := range tests {
		if got := f.Decide(event(tt.level)); got != tt.want {
			t.Errorf("ThresholdFilter(WARN).Decide(%s) = %v, want %v", tt.level, got, tt.want)
		}
	}
}

func TestLevelFilter(t *testing.T) {
	f := NewLevel(slfg.ERROR)

	tests := []struct {
		level slfg.Level
		want  slfg.FilterReply
	}{
		{slfg.DEBUG, slfg.DENY},
		{slfg.INFO, slfg.DENY},
		{slfg.ERROR, slfg.ACCEPT},
		{slfg.FATAL, slfg.DENY},
	}
	for _, tt := range tests {
		if got := f.Decide(event(tt.level)); got != tt.want {
			t.Errorf("LevelFilter(ERROR).Decide(%s) = %v, want %v", tt.level, got, tt.want)
		}
	}
}
