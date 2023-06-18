package process

import (
	"log"
	"strings"
	"sync"
)

type ProcessManager struct {
	procs map[string][]*process
	mu    *sync.Mutex
}

func NewProcessManager() *ProcessManager {
	return &ProcessManager{
		procs: make(map[string][]*process),
		mu:    &sync.Mutex{},
	}
}

func (pm *ProcessManager) Add(path, command string) {
	p := newProcess(command)

	if _, ok := pm.procs[path]; !ok {
		pm.procs[path] = make([]*process, 0)
	}

	pm.procs[path] = append(pm.procs[path], p)
}

func (pm *ProcessManager) StartAll() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for _, procs := range pm.procs {
		for _, p := range procs {
			err := p.Start()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (pm *ProcessManager) StopAll() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for _, procs := range pm.procs {
		for _, p := range procs {
			err := p.Kill()
			if err != nil {
				log.Println(err)
			}
		}
	}

	return nil
}

func (pm *ProcessManager) GetProcs(path string) []*process {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for procPath, procs := range pm.procs {
		if strings.HasPrefix(path, procPath) {
			return procs
		}
	}

	return []*process{}
}
