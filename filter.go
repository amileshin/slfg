package slfg

// FilterReply — результат решения фильтра.
type FilterReply int

const (
	// DENY — событие отклоняется, дальнейшие фильтры не проверяются.
	DENY FilterReply = iota
	// NEUTRAL — фильтр не принял решения, проверка продолжается.
	NEUTRAL
	// ACCEPT — событие принимается, дальнейшие фильтры не проверяются.
	ACCEPT
)

// Filter — интерфейс фильтра событий логирования.
type Filter interface {
	Decide(event *LoggingEvent) FilterReply
}

// FilterChain — цепочка фильтров: ACCEPT → DENY → NEUTRAL.
// Если все фильтры вернули NEUTRAL, событие принимается.
type FilterChain struct {
	filters []Filter
}

func NewFilterChain(filters ...Filter) *FilterChain {
	return &FilterChain{filters: filters}
}

func (c *FilterChain) Decide(event *LoggingEvent) FilterReply {
	for _, f := range c.filters {
		reply := f.Decide(event)
		if reply == ACCEPT || reply == DENY {
			return reply
		}
	}
	return NEUTRAL
}

// Add добавляет фильтр в конец цепочки.
func (c *FilterChain) Add(f Filter) {
	c.filters = append(c.filters, f)
}
