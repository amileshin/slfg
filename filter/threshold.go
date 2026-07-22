package filter

import "slfg"

// ThresholdFilter — пропускает события с уровнем ≥ заданного.
type ThresholdFilter struct {
	Level slfg.Level
}

func NewThreshold(level slfg.Level) *ThresholdFilter {
	return &ThresholdFilter{Level: level}
}

func (f *ThresholdFilter) Decide(event *slfg.LoggingEvent) slfg.FilterReply {
	if event.Level >= f.Level {
		return slfg.NEUTRAL
	}
	return slfg.DENY
}
