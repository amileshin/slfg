package slfg

import "io"

// Appender — интерфейс вывода событий логирования.
type Appender interface {
	Append(event *LoggingEvent)
	Close() error
}

// Filterable — appender с поддержкой цепочки фильтров.
type Filterable interface {
	SetFilterChain(chain *FilterChain)
}

// WriterAppender — appender, записывающий в io.Writer.
type WriterAppender struct {
	w         io.Writer
	formatter Formatter
	chain     *FilterChain
}

func NewWriterAppender(w io.Writer, formatter Formatter) *WriterAppender {
	if formatter == nil {
		formatter = DefaultFormatter
	}
	return &WriterAppender{w: w, formatter: formatter}
}

func (a *WriterAppender) SetFilterChain(chain *FilterChain) {
	a.chain = chain
}

func (a *WriterAppender) Append(event *LoggingEvent) {
	if a.chain != nil && a.chain.Decide(event) == DENY {
		return
	}
	a.w.Write(a.formatter(event))
}

func (a *WriterAppender) Close() error {
	return nil
}
