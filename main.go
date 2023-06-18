package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/mnafees/liver/internal/process"
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

	pm := process.NewProcessManager()

	watcher := watcher.NewWatcher()
	defer watcher.Close()

	for _, path := range c.Paths {
		err = watcher.Add(path)
		if err != nil {
			log.Fatalf("error adding watcher for %s: %v\n", path, err)
		}
	}

	for path, command := range c.Procs {
		err = pm.Add(path, command)
		if err != nil {
			log.Fatalf("error adding process for path %s: %v\n", path, err)
		}
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig)

	go func() {
		err := pm.Start()
		if err != nil {
			log.Fatalf("error starting processes: %v\n", err)
			sig <- os.Interrupt
			return
		}

		for {
			select {
			case event, ok := <-watcher.Events():
				if !ok || !pm.Valid(event.Name) {
					return
				}

				err := pm.Stop()
				if err != nil {
					log.Fatalf("error stopping processes: %v\n", err)
				}

				err = pm.Start()
				if err != nil {
					log.Fatalf("error starting processes: %v\n", err)
				}
			case err, ok := <-watcher.Errors():
				if !ok {
					return
				}

				log.Println("watcher error: ", err)
			}
		}
	}()

	<-sig

	log.Println("exiting gracefully")

	err = pm.Stop()
	if err != nil {
		log.Fatalf("error stopping processes: %v\n", err)
	}
}
