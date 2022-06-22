package uilist

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mxmaxime/tedis/tui/constants"
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

func (m *ListModel) activeKey() string {
	// old way of finding the active key,
	// doesn't work with the built-in list filter
	//items := m.list.Items()
	//activeItem := items[m.list.Index()]

	item, ok := m.list.SelectedItem().(ListItem)
	if !ok {
		// todo
	}

	return item.Key
}

func selectItemCmd(activeKey string) tea.Cmd {
	return func() tea.Msg {
		fmt.Println("in select item cmd, selected key: ", activeKey)
		return SelectMsg{ActiveRedisKey: activeKey}
	}
}

func (m *ListModel) onSizeChange(msg tea.WindowSizeMsg) {
	width, height := msg.Width, msg.Height
	//fmt.Printf("width=%d height=%d\n", width, height)

	m.list.SetSize(width, height)
}

// when a key is pressed
func (m ListModel) onKey(msg tea.KeyMsg) (ListModel, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	if m.list.FilterState() == list.Filtering {
		//fmt.Println(filterVal)
		//filterVal := m.list.FilterValue()
	}

	switch {
	case key.Matches(msg, constants.Keymap.Enter):
		cmd = selectItemCmd(m.activeKey())
		//fmt.Printf("onKey pressed in list view, cmd = %v\n", cmd)
	}

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.onSizeChange(msg)
	case tea.KeyMsg:
		m, cmd = m.onKey(msg)
	}

	cmds = append(cmds, cmd)

	// built-in list update (navigation..)
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m ListModel) View() string {
	return listViewStyle.Render(m.list.View())
}
