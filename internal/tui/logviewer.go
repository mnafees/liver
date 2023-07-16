package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mnafees/liver/internal/sharedbuffer"
	"github.com/muesli/reflow/wrap"
)

type UpdateLogViewerMsg struct{}

type LogViewer struct {
	bufferFactory  *sharedbuffer.Factory
	viewport       viewport.Model
	currentProcIdx uint
}

func NewLogViewer(factory *sharedbuffer.Factory) LogViewer {
	return LogViewer{
		bufferFactory: factory,
	}
}

// Model methods

func (lv LogViewer) Init() tea.Cmd {
	return nil
}

func (lv LogViewer) Update(msg tea.Msg) (LogViewer, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		lv.viewport.Width = msg.Width
		lv.viewport.Height = msg.Height - 20

		lv.viewport.SetContent(
			wrap.String(
				lv.bufferFactory.Get(lv.currentProcIdx).RawBuffer().String(),
				msg.Width,
			),
		)
	case UpdateLogViewerMsg:
		lv.viewport.SetContent(
			wrap.String(
				lv.bufferFactory.Get(lv.currentProcIdx).RawBuffer().String(),
				lv.viewport.Width,
			),
		)
	}

	lv.viewport, cmd = lv.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return lv, tea.Batch(cmds...)
}

func (lv LogViewer) View() string {
	return lv.viewport.View()
}

// Helpers

func (lv LogViewer) SetCurrentProcIdx(idx uint) {
	lv.currentProcIdx = idx
}
