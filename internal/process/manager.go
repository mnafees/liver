package process

import (
	"fmt"
	"log"
	"strings"
)

type ProcessManager struct {
	root  *trie
	procs []*process
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
	for _, p := range pm.procs {
		err := p.start()
		if err != nil {
			return err
		}
	}

	return nil
}

func (pm *ProcessManager) Stop() error {
	for _, p := range pm.procs {
		err := p.kill()
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}
