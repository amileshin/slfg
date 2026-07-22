package appender

import (
	"io"
	"os"

	"slfg"
)

// ConsoleAppender — вывод в stdout или stderr.
type ConsoleAppender struct {
	writer    io.Writer
	formatter slfg.Formatter
	chain     *slfg.FilterChain
}

// NewConsole создаёт ConsoleAppender.
// target: "System.out" / "stdout" → stdout, иначе stderr.
func NewConsole(target string, formatter slfg.Formatter) *ConsoleAppender {
	var w io.Writer
	switch target {
	case "System.out", "stdout":
		w = os.Stdout
	default:
		w = os.Stderr
	}
	if formatter == nil {
		formatter = slfg.DefaultFormatter
	}
	return &ConsoleAppender{writer: w, formatter: formatter}
}

func (a *ConsoleAppender) SetFilterChain(chain *slfg.FilterChain) {
	a.chain = chain
}

func (a *ConsoleAppender) Append(event *slfg.LoggingEvent) {
	if a.chain != nil && a.chain.Decide(event) == slfg.DENY {
		return
	}
	a.writer.Write(a.formatter(event))
}

func (a *ConsoleAppender) Close() error {
	return nil
}
