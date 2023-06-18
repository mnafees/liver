package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/mnafees/liver/internal/watcher"
)

type config struct {
	Paths []string          `json:"paths"`
	Procs map[string]string `json:"procs"`
}

func main() {
	fileBytes, err := os.ReadFile("liver.json")
	if err != nil {
		log.Fatalf("error reading liver.json: %v\n", err)
	}

	c := &config{
		Paths: make([]string, 0),
		Procs: make(map[string]string),
	}

	err = json.Unmarshal(fileBytes, c)
	if err != nil {
		log.Fatalf("error unmarshalling liver.json: %v\n", err)
	}

	if len(c.Paths) == 0 {
		log.Fatalln("no paths specified")
	}

	if len(c.Procs) == 0 {
		log.Fatalln("no procs specified")
	}

	watcher := watcher.NewWatcher()
	defer watcher.Close()

	for _, path := range c.Paths {
		err = watcher.Add(path)
		if err != nil {
			log.Fatalf("error adding watcher for %s: %v\n", path, err)
		}
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events():
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Write) {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors():
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Block main goroutine forever.
	<-make(chan struct{})
}
