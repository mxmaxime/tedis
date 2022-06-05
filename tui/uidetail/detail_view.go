package uidetail

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/tidwall/pretty"

	"github.com/mxmaxime/tedis/myredis"
	"github.com/mxmaxime/tedis/tui/constants"
	"github.com/mxmaxime/tedis/utils"
)

const minWidth = 30

type Item struct {
	Key string
	Val map[string]string

	// json pretty representation
	PrettyJson []byte

	// pretty + colorized, for terminal output
	PrettyColorizedJson []byte

	// json
	Json string
}

type sessionState = int

const (
	jsonView sessionState = iota
	tableView
	confirmationView
)

type confirmationModel struct {
	diff []byte

	// file where the modification is written to
	filepath string
}

type DetailModel struct {
	state    sessionState
	RedisCli *redis.Client

	RedisRepo *myredis.RedisRepo

	item Item
	err  string

	totalMargin int
	totalWidth  int
	table       table.Model
	viewport    viewport.Model
	ready       bool

	confirmationModel confirmationModel
}

func New(cli *redis.Client, selectedKey string) *DetailModel {
	ctx := context.TODO()

	repo := myredis.RedisRepo{Cli: cli}

	m := DetailModel{
		RedisCli:  cli,
		RedisRepo: &repo,
	}

	val, err := cli.HGetAll(ctx, selectedKey).Result()
	if err != nil {
		m.err = err.Error()
	}

	valStr, err := json.Marshal(val)
	if err != nil {
		m.err = err.Error()
	}

	niceJson := pretty.Pretty(valStr)

	colorized := pretty.Color(niceJson, pretty.TerminalStyle)

	item := Item{
		Key:                 selectedKey,
		Val:                 val,
		PrettyJson:          niceJson,
		PrettyColorizedJson: colorized,
	}

	m.item = item

	// build table
	var rows []table.Row

	for key, val := range m.item.Val {
		row := table.NewRow(table.RowData{
			columnKey:   key,
			columnValue: val,
		})
		rows = append(rows, row)
	}

	table := table.New([]table.Column{
		table.NewFlexColumn(columnKey, "Key", 20),
		table.NewFlexColumn(columnValue, "Value", 20),
	}).WithRows(rows)

	m.table = table

	// build view port
	vp := viewport.New(100, 10)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	m.viewport = vp
	m.setViewportContent()

	return &m
}

func (m *DetailModel) setViewportContent() {
	str, err := json.Marshal(m.item.Val)
	if err != nil {
		m.viewport.SetContent("got an error while decoding the data.")
	}

	niceJson := pretty.Pretty(str)
	colorized := pretty.Color(niceJson, pretty.TerminalStyle)

	m.viewport.SetContent(string(colorized))
}

// Init run any intial IO on program start
func (m DetailModel) Init() tea.Cmd {
	return nil
}

/**
 * Messages
 */

// BackMsg change state back to project view
type BackMsg bool

type errMsg struct{ error } // TODO: have this implement Error()
type updateEntryListMsg struct{ input []byte }
type updatedMsg struct{}

// editor closed
type editorFinishedMsg struct {
	err      error
	filepath string
}

type confirmationMsg struct {
	state confirmationState
}

type confirmationState = int

const (
	write int = iota
	cancel
	continueEdit
)

/**
 * Commands
 */

func msgToCmd(msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

func BackCmd() tea.Cmd {
	return func() tea.Msg {
		return BackMsg(true)
	}
}

func errCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return errMsg{err}
	}
}

func (m DetailModel) openEditorCmd() tea.Cmd {
	file, err := os.CreateTemp(os.TempDir(), "")
	if err != nil {
		return errCmd(errors.Wrap(err, "cannot create temp file"))
	}
	defer file.Close()

	if _, err := file.Write(m.item.PrettyJson); err != nil {
		return errCmd(errors.Wrap(err, "cannot write temp file"))
	}

	editorEnv, ok := os.LookupEnv("EDITOR")
	if !ok {
		editorEnv = "vi"
	}

	filename := file.Name()
	c := exec.Command(editorEnv, filename)

	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err, filename}
	})
}

func (m DetailModel) onEditorFinished(msg editorFinishedMsg) (DetailModel, tea.Cmd) {
	_, err := os.ReadFile(msg.filepath)
	if err != nil {
		// design choice: if I can't read the file to show the diff that will be saved, skip the saving.
		return m, errCmd(errors.Wrap(err, "cannot read file to diff it. Save to redis is canceled."))
	}

	//newData := string(dataBytes)
	//fmt.Println(newData)
	m.state = confirmationView
	m.confirmationModel = confirmationModel{diff: []byte("some diff to implement here"), filepath: msg.filepath}
	return m, nil
}

func (m DetailModel) onWrite(msg confirmationMsg) (DetailModel, tea.Cmd) {

	return m, nil
}

func (m DetailModel) onCancel(msg confirmationMsg) (DetailModel, tea.Cmd) {
	// not needed anymore
	defer os.Remove(m.confirmationModel.filepath)

	m.state = jsonView

	return m, nil
}

func (m DetailModel) onContinueEdit(msg confirmationMsg) (DetailModel, tea.Cmd) {
	filename := m.confirmationModel.filepath

	editorEnv, ok := os.LookupEnv("EDITOR")
	if !ok {
		editorEnv = "vi"
	}

	c := exec.Command(editorEnv, filename)

	return m, tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err, filename}
	})
}

//func (m DetailModel) onEditorFinished(msg askEditConfirmationMsg) (DetailModel, tea.Cmd) {
//}

func (m DetailModel) onKeys(msg tea.KeyMsg) (DetailModel, tea.Cmd) {
	var cmd tea.Cmd

	switch m.state {
	case confirmationView:
		switch {
		// cancel: delete changes + don't write
		case key.Matches(msg, constants.Keymap.Cancel):
			// write things to redis
			return m, msgToCmd(confirmationMsg{state: cancel})
		case key.Matches(msg, constants.Keymap.Save):
			// continue edditing
			return m, msgToCmd(confirmationMsg{state: write})
		case key.Matches(msg, constants.Keymap.Create):
			return m, msgToCmd(confirmationMsg{state: continueEdit})
		}
	default:
		switch {
		case key.Matches(msg, constants.Keymap.Back):
			cmd = BackCmd()
		case key.Matches(msg, constants.Keymap.Create):
			return m, m.openEditorCmd()
		}
	}
	/**
	used for table
	switch msg.String() {
	case "left":
		if m.totalWidth-m.totalMargin > minWidth {
			m.totalMargin++
			m.recalculateTable()
		}

	case "right":
		if m.totalMargin > 0 {
			m.totalMargin--
			m.recalculateTable()
		}
	}
	*/

	return m, cmd
}

func (m DetailModel) onWindowSizeChange(msg tea.WindowSizeMsg) (DetailModel, tea.Cmd) {
	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())

	verticalMarginHeight := headerHeight + footerHeight

	// Since this program is using the full size of the viewport we
	// need to wait until we've received the window dimensions before
	// we can initialize the viewport. The initial dimensions come in
	// quickly, though asynchronously, which is why we wait for them
	// here.
	m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
	//m.viewport.Width = msg.Width
	//m.viewport.Height = msg.Height - verticalMarginHeight

	m.viewport.YPosition = headerHeight
	m.setViewportContent()

	// This is only necessary for high performance rendering, which in
	// most cases you won't need.
	//
	// Render the viewport one line below the header.
	m.viewport.YPosition = headerHeight + 1
	/*
		m.setViewportContent()
	*/
	return m, nil
}

func (m DetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	//fmt.Println("detail view got updated")

	//m.table, cmd = m.table.Update(msg)
	//cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m, cmd = m.onKeys(msg)
		cmds = append(cmds, cmd)
	case tea.WindowSizeMsg:
		//fmt.Println("window size changed")
		m, cmd = m.onWindowSizeChange(msg)
		cmds = append(cmds, cmd)

		// was used to table
		//m.totalWidth = msg.Width
		//m.recalculateTable()
	case errMsg:
		// should I stop the switch case?
		m.err = msg.Error()
	case editorFinishedMsg:
		m, cmd = m.onEditorFinished(msg)
	case confirmationMsg:
		fmt.Println("got confirmation msg")
		switch msg.state {
		case write:
		case continueEdit:
			return m.onContinueEdit(msg)
		case cancel:
			return m.onCancel(msg)
		default:
			return m, errCmd(errMsg{errors.New("unkown state.")})
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

const (
	columnKey   = "key"
	columnValue = "value"
)

func (m *DetailModel) recalculateTable() {
	fmt.Println("recalculateTable called")
	m.table = m.table.WithTargetWidth(m.totalWidth - m.totalMargin)
}

/**
 * VIEWS
 */

func (m DetailModel) headerView() string {
	title := titleStyle.Render(fmt.Sprintf("Viewing key '%s'", m.item.Key))

	line := strings.Repeat("─", utils.MaxInt(0, m.viewport.Width-lipgloss.Width(title)))

	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m DetailModel) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))

	line := strings.Repeat("─", utils.MaxInt(0, m.viewport.Width-lipgloss.Width(info)))

	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m DetailModel) View() string {
	switch m.state {
	case confirmationView:
		return "Press q to cancel, w to write your changes in Redis, and i to go back in the edditing mode"
	default:
		return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
	}

	/*
		wip table stuff
		strs := []string{
			"A flexible table that fills available space (Name is fixed-width)",
			fmt.Sprintf("Target size: %d (left/right to adjust)", m.totalWidth-m.totalMargin),
			"Press q or ctrl+c to quit",
			m.table.View(),
		}

		return lipgloss.JoinVertical(lipgloss.Left, strs...) + "\n"
	*/
}
