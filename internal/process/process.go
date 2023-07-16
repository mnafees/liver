package process

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/mnafees/liver/internal/sharedbuffer"
)

type process struct {
	bufferFactory      *sharedbuffer.Factory
	internalProcessCmd *exec.Cmd
	mu                 *sync.Mutex
	args               []string
	idx                uint
}

func newProcess(idx uint, factory *sharedbuffer.Factory, command string) *process {
	args := strings.Split(command, " ")

	return &process{
		bufferFactory: factory,
		mu:            &sync.Mutex{},
		args:          args,
		idx:           idx,
	}
}

func (p *process) Start() error {
	if p == nil {
		return nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	cmd := exec.Command(p.args[0], p.args[1:]...)
	cmd.Stderr = p.bufferFactory.Get(p.idx)
	cmd.Stdout = p.bufferFactory.Get(p.idx)
	// FIXME: implement input

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("could not start process: %w", err)
	}

	p.internalProcessCmd = cmd

	return err
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
