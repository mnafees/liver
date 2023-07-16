package process

import (
	"log"
	"strings"
	"sync"

	"github.com/mnafees/liver/internal/sharedbuffer"
)

type PathPrefix string

type ProcessManager struct {
	bufferFactory *sharedbuffer.Factory
	procs         map[PathPrefix][]*process
	mu            *sync.Mutex
}

func NewProcessManager(factory *sharedbuffer.Factory) *ProcessManager {
	return &ProcessManager{
		bufferFactory: factory,
		procs:         make(map[PathPrefix][]*process),
		mu:            &sync.Mutex{},
	}
}

func (pm *ProcessManager) Add(path PathPrefix, command string) {
	p := newProcess(uint(len(pm.procs)), pm.bufferFactory, command)

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

func (pm *ProcessManager) GetProcNames() []string {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	var procNames []string

	for _, procs := range pm.procs {
		for _, p := range procs {
			procNames = append(procNames, strings.Join(p.args, " "))
		}
	}

	return procNames
}

func (pm *ProcessManager) GetProcs(path string) []*process {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	var processes []*process

	for prefix, procs := range pm.procs {
		for _, proc := range procs {
			if strings.HasPrefix(path, string(prefix)) {
				processes = append(processes, proc)
			}
		}
	}

	return processes
}
