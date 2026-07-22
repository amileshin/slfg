package slfg

import (
	"context"
	"testing"
)

func TestWithMDC(t *testing.T) {
	ctx := context.Background()
	ctx = WithMDC(ctx, "requestId", "abc-123")
	ctx = WithMDC(ctx, "userId", "42")

	mdc := MDCFrom(ctx)
	if mdc["requestId"] != "abc-123" {
		t.Errorf("requestId = %q, want %q", mdc["requestId"], "abc-123")
	}
	if mdc["userId"] != "42" {
		t.Errorf("userId = %q, want %q", mdc["userId"], "42")
	}
}

func TestMDCFromEmpty(t *testing.T) {
	mdc := MDCFrom(context.Background())
	if len(mdc) != 0 {
		t.Errorf("ожидалась пустая карта, получено: %v", mdc)
	}
}

func TestWithMDCImmutable(t *testing.T) {
	ctx1 := context.Background()
	ctx2 := WithMDC(ctx1, "key", "val1")
	ctx3 := WithMDC(ctx2, "key", "val2")

	if MDCFrom(ctx2)["key"] != "val1" {
		t.Error("исходный контекст не должен меняться")
	}
	if MDCFrom(ctx3)["key"] != "val2" {
		t.Error("новый контекст должен содержать обновлённое значение")
	}
}
