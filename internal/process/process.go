package process

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

type errWriter struct {
	idx uint
}

func (ew *errWriter) Write(p []byte) (int, error) {
	return fmt.Fprintf(os.Stderr, "[%d]: %s", ew.idx, string(p))
}

type outWriter struct {
	idx uint
}

func (ow *outWriter) Write(p []byte) (int, error) {
	return fmt.Fprintf(os.Stdout, "[%d]: %s", ow.idx, string(p))
}

type process struct {
	internalProcessCmd *exec.Cmd
	mu                 *sync.Mutex
	args               []string
	idx                uint
}

func newProcess(idx uint, command string) *process {
	args := strings.Split(command, " ")

	return &process{
		mu:   &sync.Mutex{},
		args: args,
		idx:  idx,
	}
}

func (p *process) Start() error {
	if p == nil {
		return nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	cmd := exec.Command(p.args[0], p.args[1:]...)
	cmd.Stderr = &errWriter{idx: p.idx}
	cmd.Stdout = &outWriter{idx: p.idx}
	cmd.Stdin = os.Stdin

	p.internalProcessCmd = cmd

	return p.internalProcessCmd.Start()
}

func (p *process) Kill() error {
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
