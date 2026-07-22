package config

import (
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// Watcher отслеживает изменения файла конфигурации.
type Watcher struct {
	path      string
	onChange  func()
	fsWatcher *fsnotify.Watcher
	done      chan struct{}
}

// NewWatcher создаёт Watcher для указанного файла.
func NewWatcher(path string, onChange func()) (*Watcher, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		fw.Close()
		return nil, err
	}
	if err := fw.Add(filepath.Dir(abs)); err != nil {
		fw.Close()
		return nil, err
	}
	return &Watcher{
		path:      abs,
		onChange:  onChange,
		fsWatcher: fw,
		done:      make(chan struct{}),
	}, nil
}

// Start запускает отслеживание в отдельной горутине.
func (w *Watcher) Start() {
	go w.loop()
}

func (w *Watcher) loop() {
	for {
		select {
		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				return
			}
			if filepath.Clean(event.Name) != w.path {
				continue
			}
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				w.onChange()
			}
		case _, ok := <-w.fsWatcher.Errors:
			if !ok {
				return
			}
		case <-w.done:
			return
		}
	}
}

// Stop останавливает отслеживание.
func (w *Watcher) Stop() error {
	close(w.done)
	return w.fsWatcher.Close()
}
