package process

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

type process struct {
	internalProcessCmd *exec.Cmd
	mu                 *sync.Mutex
}

func newProcess(command string) *process {
	args := strings.Split(command, " ")

	cmd := exec.Command(args[0], args[1:]...)

	return &process{
		internalProcessCmd: cmd,
		mu:                 &sync.Mutex{},
	}
}

func (p *process) start() error {
	if p == nil || p.internalProcessCmd == nil {
		return nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	return p.internalProcessCmd.Start()
}

func (p *process) kill() error {
	if p == nil || p.internalProcessCmd == nil || p.internalProcessCmd.Process == nil {
		return nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// first try to kill the process using the internal Process
	err := p.internalProcessCmd.Process.Kill()
	if err == nil {
		return nil
	}

	// try to do a kill -9 for the process PID
	return exec.Command("kill", "-9", fmt.Sprintf("%d", p.internalProcessCmd.Process.Pid)).Run()
}
