package uilist

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	listViewStyle = lipgloss.NewStyle().
		PaddingRight(1).
		MarginRight(1).
		Border(lipgloss.RoundedBorder(), false, true, false, false)
)

// SelectMsg the message to change the view to the selected entry
type SelectMsg struct {
	ActiveRedisKey string
}

func (m *ListModel) onSizeChange(msg tea.WindowSizeMsg) {
	width, height := msg.Width, msg.Height
	fmt.Printf("width=%d height=%d\n", width, height)

	m.list.SetSize(width, height)
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.onSizeChange(msg)
	}

	// built-in list update (navigation..)
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m ListModel) View() string {
	return listViewStyle.Render(m.list.View())
}
