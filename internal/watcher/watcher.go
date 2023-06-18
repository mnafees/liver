package watcher

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	internalWatcher *fsnotify.Watcher
}

func NewWatcher() *Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	return &Watcher{
		internalWatcher: watcher,
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

	if filepath.Base(path) == "liver.json" {
		return nil
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
