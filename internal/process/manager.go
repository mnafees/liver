package process

import (
	"fmt"
	"log"
	"strings"
	"sync"
)

type ProcessManager struct {
	root  *trie
	procs []*process
	mu    *sync.Mutex
}

func NewProcessManager() *ProcessManager {
	rootNode := &trieNode{
		children: make(map[string]*trieNode),
	}

	rootNode.children["/"] = newNode()

	return &ProcessManager{
		root: &trie{
			node: rootNode,
		},
		procs: make([]*process, 0),
		mu:    &sync.Mutex{},
	}
}

func (pm *ProcessManager) Add(path, command string) error {
	p := newProcess(command)
	chars := strings.Split(path, "")

	if chars[0] != "/" {
		return fmt.Errorf("expected UNIX filesystem absolute path starting with '/', got %s", path)
	}

	chars = chars[1:]

	curr := pm.root.node

	for _, c := range chars {
		if _, ok := curr.children[c]; !ok {
			curr.children[c] = newNode()
		}

		curr = curr.children[c]
	}

	curr.proc = p

	pm.procs = append(pm.procs, p)

	return nil
}

func (pm *ProcessManager) Start() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	log.Println("starting processes")

	for _, p := range pm.procs {
		err := p.start()
		if err != nil {
			return err
		}
	}

	return nil
}

func (pm *ProcessManager) Stop() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	log.Println("stopping processes")

	for _, p := range pm.procs {
		err := p.kill()
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

func (pm *ProcessManager) Valid(path string) bool {
	chars := strings.Split(path, "")

	if chars[0] != "/" {
		return false
	}

	chars = chars[1:]

	curr := pm.root.node

	for _, c := range chars {
		if _, ok := curr.children[c]; !ok {
			return false
		}

		curr = curr.children[c]
	}

	return curr.isLeaf()
}
