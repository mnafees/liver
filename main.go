package main

import (
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mnafees/liver/internal/process"
	"github.com/mnafees/liver/internal/sharedbuffer"
	"github.com/mnafees/liver/internal/tui"
	"github.com/mnafees/liver/internal/watcher"
)

type config struct {
	Paths []string            `json:"paths"`
	Procs map[string][]string `json:"procs"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func main() {
	writeChan := make(chan struct{})
	defer close(writeChan)

	bufferFactory := sharedbuffer.NewFactory(writeChan)

	logViewer := tui.NewLogViewer(bufferFactory)
	container := tui.NewContainer(logViewer)

	p := tea.NewProgram(
		container,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

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

	pm := process.NewProcessManager(bufferFactory)

	watcher := watcher.NewWatcher()
	defer watcher.Close()

	for _, path := range c.Paths {
		err = watcher.Add(path)
		if err != nil {
			log.Fatalf("error adding watcher for %s: %v\n", path, err)
		}
	}

	for path, commands := range c.Procs {
		for _, c := range commands {
			pm.Add(process.PathPrefix(path), c)
		}
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {
		if _, err := p.Run(); err != nil {
			sig <- os.Interrupt
		}

		go func() {
			for range writeChan {
				logViewer.Update(tui.UpdateLogViewerMsg{})
			}
		}()

		var (
			waitFor = 2 * time.Second

			// Keep track of the timers, as path → timer.
			mu     sync.Mutex
			timers = make(map[string]*time.Timer)
		)

		err := pm.StartAll()
		if err != nil {
			sig <- os.Interrupt
			return
		}

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
						for _, p := range procs {
							// err := p.Kill()
							// if err != nil {

							// }
							p.Kill()
						}

						for _, p := range procs {
							// err = p.Start()
							// if err != nil {
							// 	log.Fatalf("error starting processes: %v\n", err)
							// }
							p.Start()
						}

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
			case _, ok := <-watcher.Errors():
				if !ok {
					sig <- os.Interrupt
					return
				}
			}
		}
	}()

	<-sig

	err = pm.StopAll()
	if err != nil {
		log.Fatalf("error stopping processes: %v\n", err)
	}
}
