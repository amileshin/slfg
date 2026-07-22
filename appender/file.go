package appender

import (
	"os"

	"slfg"
)

// FileAppender — запись в файл.
type FileAppender struct {
	file      *os.File
	formatter slfg.Formatter
	chain     *slfg.FilterChain
}

// NewFile создаёт FileAppender.
// appendMode: true — дозапись, false — перезапись.
func NewFile(path string, appendMode bool, formatter slfg.Formatter) (*FileAppender, error) {
	flags := os.O_CREATE | os.O_WRONLY
	if appendMode {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}
	f, err := os.OpenFile(path, flags, 0o644)
	if err != nil {
		return nil, err
	}
	if formatter == nil {
		formatter = slfg.DefaultFormatter
	}
	return &FileAppender{file: f, formatter: formatter}, nil
}

func (a *FileAppender) SetFilterChain(chain *slfg.FilterChain) {
	a.chain = chain
}

func (a *FileAppender) Append(event *slfg.LoggingEvent) {
	if a.chain != nil && a.chain.Decide(event) == slfg.DENY {
		return
	}
	a.file.Write(a.formatter(event))
}

func (a *FileAppender) Close() error {
	return a.file.Close()
}
