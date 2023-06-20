package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/mnafees/liver/internal/process"
	"github.com/mnafees/liver/internal/watcher"
)

type config struct {
	Paths []string            `json:"paths"`
	Procs map[string][]string `json:"procs"`
}

func main() {
	fileBytes, err := os.ReadFile("liver.json")
	if err != nil {
		log.Fatalf("error reading liver.json: %v\n", err)
	}

	c := &config{
		Paths: make([]string, 0),
		Procs: make(map[string][]string),
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

	idx := uint(0)

	fmt.Println()

	for path, commands := range c.Procs {
		for _, c := range commands {
			fmt.Printf("setting process %d for: %s\n", idx, c)

			pm.Add(idx, path, c)
			idx++
		}
	}

	fmt.Println()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {
		var (
			waitFor = 2 * time.Second

			// Keep track of the timers, as path → timer.
			mu     sync.Mutex
			timers = make(map[string]*time.Timer)
		)

		log.Println("starting all processes")

		err := pm.StartAll()
		if err != nil {
			log.Fatalf("error starting processes: %v\n", err)
			sig <- os.Interrupt
			return
		}

		log.Printf("started all processes\n\n")

		for {
			select {
			case event, ok := <-watcher.Events():
				if !ok {
					sig <- os.Interrupt
					return
				}

				procs := pm.GetProcs(event.Name)
				if len(procs) == 0 {
					continue
				}

				mu.Lock()
				t, ok := timers[event.Name]
				mu.Unlock()

				if !ok {
					t = time.AfterFunc(math.MaxInt64, func() {
						fmt.Println()
						log.Printf("stopping processes")

						for _, p := range procs {
							err := p.Kill()
							if err != nil {
								log.Fatalf("error stopping processes: %v\n", err)
							}
						}

						log.Println("restarting processes")

						for _, p := range procs {
							err = p.Start()
							if err != nil {
								log.Fatalf("error starting processes: %v\n", err)
							}
						}

						log.Printf("restarted processes")
						fmt.Println()

						mu.Lock()
						delete(timers, event.Name)
						mu.Unlock()
					})
					t.Stop()

					mu.Lock()
					timers[event.Name] = t
					mu.Unlock()
				}

				t.Reset(waitFor)
			case err, ok := <-watcher.Errors():
				if !ok {
					sig <- os.Interrupt
					return
				}

				log.Println("watcher error: ", err)
			}
		}
	}()

	<-sig

	fmt.Println()
	log.Printf("stopping all processes")

	err = pm.StopAll()
	if err != nil {
		log.Fatalf("error stopping processes: %v\n", err)
	}
}
