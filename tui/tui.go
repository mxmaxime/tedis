package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-redis/redis/v8"
	"github.com/mxmaxime/tedis/tui/uilist"
	"github.com/mxmaxime/tedis/utils"
)

type sessionState = int

const (
	detailView sessionState = iota
	listView
)

type MainModel struct {
	state sessionState

	listModel   tea.Model
	detailModel tea.Model
}

func RedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     utils.StringOr(os.Getenv("REDIS_HOST"), "localhost:6379"),
		Password: utils.StringOr(os.Getenv("REDIS_PASSWORD"), ""),
		DB:       utils.StringOrInt(os.Getenv("REDIS_DATABASE"), 0),
	})
}

func New() *MainModel {
	cli := RedisClient()

	return &MainModel{
		state:     listView,
		listModel: uilist.New(cli),
	}
}

func (m MainModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

// handles all keypresses message and produces cmds if needed
func (m *MainModel) handleKeys(msg tea.KeyMsg) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg.Type {
	case tea.KeyCtrlC:
		cmd = tea.Quit
		cmds = append(cmds, cmd)
	}

	switch m.state {
	}

	return tea.Batch(cmds...)
}

func (m *MainModel) OnSizeChange(msg tea.WindowSizeMsg) {
	width, height := msg.Width, msg.Height
	fmt.Printf("width=%d height=%d\n", width, height)

	//statusBarHeight := lipgloss.Height(m.statusView())
	//height := m.height - statusBarHeight

	//listViewWidth := int(constant.ListProportion * float64(m.width))
	//listWidth := listViewWidth - listViewStyle.GetHorizontalFrameSize()
	//m.list.SetSize(width, height)

	//detailViewWidth := m.width - listViewWidth
	//m.viewport = viewport.New(detailViewWidth, height)
	//m.viewport.SetContent(m.detailView())
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		fmt.Println()
		//m.OnSizeChange(msg)
	// Is it a key press?
	case tea.KeyMsg:
		cmd = m.handleKeys(msg)
		cmds = append(cmds, cmd)
	}

	// update children views
	switch m.state {
	case listView:
		newList, cmd := m.listModel.Update(msg)
		listModel, ok := newList.(uilist.ListModel)
		if !ok {
			panic("could not perform assertion on uilist model")
		}
		m.listModel = listModel
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	switch m.state {
	case listView:
		return m.listModel.View()
	default:
		return "unknown state.. :("
	}
}
