package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	brand = ` ___       __ ___      ___ _______  _______
|"  |     |" |"  \    /"  /"     "|/"      \
||  |     ||  \   \  //  (: ______|:        |
|:  |     |:  |\\  \/. ./ \/    | |_____/   )
 \  |___  |.  | \.    //  // ___)_ //      /
( \_|:  \ /\  |\ \\   /  (:      "|:  __   \
 \_______(__\_|_) \__/    \_______|__|  \___)

Live reloading made simple, interactive, and enjoyable.`
)

var (
	logo = lipgloss.NewStyle().
		MarginLeft(1)

	highlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

	// Tabs.

	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	tab = lipgloss.NewStyle().
		Border(tabBorder, true).
		BorderForeground(highlight).
		Padding(0, 1)

	activeTab = tab.Copy().Border(activeTabBorder, true)

	tabGap = tab.Copy().
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)
)

var width = 80

type Container struct {
	lv LogViewer
}

func NewContainer(logViewer LogViewer) Container {
	return Container{
		lv: logViewer,
	}
}

func (c Container) Init() tea.Cmd {
	return nil
}

func (c Container) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return c, tea.Quit
		}

	case tea.WindowSizeMsg:
		width = msg.Width

	case tea.MouseMsg:

	}

	c.lv, cmd = c.lv.Update(msg)
	cmds = append(cmds, cmd)

	return c, tea.Batch(cmds...)
}

func (c Container) View() string {
	doc := strings.Builder{}

	// Logo
	{
		doc.WriteString(logo.Render(brand) + "\n")
	}

	// Tabs
	{
		row := lipgloss.JoinHorizontal(
			lipgloss.Top,
			activeTab.Render("Lip Gloss"),
			tab.Render("Blush"),
			tab.Render("Eye Shadow"),
			tab.Render("Mascara"),
			tab.Render("Foundation"),
		)
		gap := tabGap.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(row)-2)))
		row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
		doc.WriteString(row + "\n\n")
	}

	// Log Viewer
	{
		doc.WriteString(c.lv.View())
	}

	return doc.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
