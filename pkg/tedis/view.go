package tedis

import "github.com/charmbracelet/lipgloss"

var (
	listViewStyle = lipgloss.NewStyle().
			PaddingRight(1).
			MarginRight(1).
			Border(lipgloss.RoundedBorder(), false, true, false, false)

	dividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"})
)

func (m model) listView() string {
	return listViewStyle.Render(m.list.View())
}

func (m model) statusView() string {
	return "hello world"
}

func (m model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, m.listView(), m.contentViewport.View()),
		m.statusView(),
	)
}
