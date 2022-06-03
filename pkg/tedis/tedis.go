package tedis

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-redis/redis/v8"
)

type State = int

const (
	defaultState State = iota
	searchState
)

type model struct {
	state State

	contentViewport viewport.Model
	redis_repo      *redis_repo

	list list.Model
}

func RedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     StringOr(os.Getenv("REDIS_HOST"), "localhost:6379"),
		Password: StringOr(os.Getenv("REDIS_PASSWORD"), ""),
		DB:       StringOrInt(os.Getenv("REDIS_DATABASE"), 0),
	})
}

func InitialModel() *model {
	ctx := context.TODO()

	cli := RedisClient()
	repo := redis_repo{cli: cli}

	var items []list.Item

	keys, err := repo.GetKeys(ctx, 0, "", -1)
	if err != nil {
		panic(err)
		// well, do something I guess
	}

	fmt.Println(keys)

	for _, key := range keys {
		kt := cli.Type(ctx, key).Val()
		fmt.Printf("key: %s kt %s\n", key, kt)
		items = append(items, item{key: key})
	}

	// init the list
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Keys"
	l.SetShowFilter(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)

	return &model{
		list:       l,
		redis_repo: &repo,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

// handles all keypresses message and produces cmds if needed
func (m *model) handleKeys(msg tea.KeyMsg) tea.Cmd {
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		width, height := msg.Width, msg.Height

		//statusBarHeight := lipgloss.Height(m.statusView())
		//height := m.height - statusBarHeight

		//listViewWidth := int(constant.ListProportion * float64(m.width))
		//listWidth := listViewWidth - listViewStyle.GetHorizontalFrameSize()
		m.list.SetSize(width, height)

		//detailViewWidth := m.width - listViewWidth
		//m.viewport = viewport.New(detailViewWidth, height)
		//m.viewport.SetContent(m.detailView())
	// Is it a key press?
	case tea.KeyMsg:
		cmd = m.handleKeys(msg)
		cmds = append(cmds, cmd)
	}

	// built-in list update (navigation..)
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
