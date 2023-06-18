package process

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type process struct {
	internalProcessCmd *exec.Cmd
	mu                 *sync.Mutex
	args               []string
}

func newProcess(command string) *process {
	args := strings.Split(command, " ")

	return &process{
		mu:   &sync.Mutex{},
		args: args,
	}
}

func (p *process) start() error {
	if p == nil {
		return nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	cmd := exec.Command(p.args[0], p.args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	p.internalProcessCmd = cmd

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
	if err != nil {
		// try to do a kill -9 for the process PID
		err = exec.Command("kill", "-9", fmt.Sprintf("%d", p.internalProcessCmd.Process.Pid)).Run()
		if err != nil {
			return fmt.Errorf("could not force kill process %d: %w", p.internalProcessCmd.Process.Pid, err)
		}
	}

	p.internalProcessCmd = nil

	return nil
}
