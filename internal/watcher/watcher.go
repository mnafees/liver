package watcher

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/mnafees/liver/internal/process"
)

type Watcher struct {
	internalWatcher *fsnotify.Watcher
	pm              *process.ProcessManager
}

func NewWatcher(pm *process.ProcessManager) *Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	return &Watcher{
		internalWatcher: watcher,
		pm:              pm,
	}
}

func (w *Watcher) Close() error {
	return w.internalWatcher.Close()
}

func (w *Watcher) Add(path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("error getting absolute path for %s: %w", path, err)
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening %s: %w", path, err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error stating %s: %w", path, err)
	}

	if fileInfo.IsDir() {
		// recursively add all subdirectories to the watcher
		fis, err := os.ReadDir(path)
		if err != nil {
			return fmt.Errorf("error reading directory %s: %w", path, err)
		}

		for _, fi := range fis {
			err := w.Add(fmt.Sprintf("%s/%s", path, fi.Name()))
			if err != nil {
				return err
			}
		}

		return nil
	}

	log.Printf("watching %s\n", path)

	return w.internalWatcher.Add(path)
}

func (w *Watcher) Events() chan fsnotify.Event {
	return w.internalWatcher.Events
}

func (w *Watcher) Errors() chan error {
	return w.internalWatcher.Errors
}
