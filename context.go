package slfg

import "context"

type mdcKey struct{}

// WithMDC возвращает новый контекст с добавленным MDC-полем.
func WithMDC(ctx context.Context, key, value string) context.Context {
	mdc := MDCFrom(ctx)
	updated := make(map[string]string, len(mdc)+1)
	for k, v := range mdc {
		updated[k] = v
	}
	updated[key] = value
	return context.WithValue(ctx, mdcKey{}, updated)
}

// MDCFrom извлекает все MDC-поля из контекста.
// Возвращает пустую карту, если MDC не задан.
func MDCFrom(ctx context.Context) map[string]string {
	if mdc, ok := ctx.Value(mdcKey{}).(map[string]string); ok {
		return mdc
	}
	return map[string]string{}
}
