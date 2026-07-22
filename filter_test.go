package slfg

import (
	"testing"
	"time"
)

func TestFilterChainAccept(t *testing.T) {
	chain := NewFilterChain(&acceptFilter{})
	event := &LoggingEvent{Level: INFO, Time: time.Now()}

	if chain.Decide(event) != ACCEPT {
		t.Error("ожидался ACCEPT")
	}
}

func TestFilterChainDeny(t *testing.T) {
	chain := NewFilterChain(&denyFilter{})
	event := &LoggingEvent{Level: INFO, Time: time.Now()}

	if chain.Decide(event) != DENY {
		t.Error("ожидался DENY")
	}
}

func TestFilterChainNeutralAccepts(t *testing.T) {
	chain := NewFilterChain(&neutralFilter{}, &neutralFilter{})
	event := &LoggingEvent{Level: INFO, Time: time.Now()}

	if chain.Decide(event) != NEUTRAL {
		t.Error("все NEUTRAL → событие принимается (NEUTRAL)")
	}
}

func TestFilterChainShortCircuit(t *testing.T) {
	chain := NewFilterChain(&neutralFilter{}, &denyFilter{}, &acceptFilter{})
	event := &LoggingEvent{Level: INFO, Time: time.Now()}

	if chain.Decide(event) != DENY {
		t.Error("DENY должен прервать цепочку до ACCEPT")
	}
}

func TestFilterChainEmpty(t *testing.T) {
	chain := NewFilterChain()
	event := &LoggingEvent{Level: INFO, Time: time.Now()}

	if chain.Decide(event) != NEUTRAL {
		t.Error("пустая цепочка → NEUTRAL")
	}
}

type acceptFilter struct{}

func (f *acceptFilter) Decide(*LoggingEvent) FilterReply { return ACCEPT }

type denyFilter struct{}

func (f *denyFilter) Decide(*LoggingEvent) FilterReply { return DENY }

type neutralFilter struct{}

func (f *neutralFilter) Decide(*LoggingEvent) FilterReply { return NEUTRAL }
