package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"strings"
)

type Writer struct {
	program *tea.Program
}

func (writer *Writer) Write(data []byte) (int, error) {
	writer.program.Send(ExecData{Data: data})
	return len(data), nil
}

type ExecData struct {
	Data []byte
}

type ExecProgram struct {
	Program *tea.Program
}

type ExecStyles struct {
	Title lipgloss.Style
	Info  lipgloss.Style
}

type Exec struct {
	Styles   ExecStyles
	program  *tea.Program
	viewport viewport.Model

	data  []byte
	limit int

	options     []string
	selectList  list.Model
	hasSelected bool
	selected    int

	done    bool
	Handler func(*tea.Program, string) func()
}

func (exec *Exec) Init() tea.Cmd {
	exec.viewport = viewport.New(10, 10)
	exec.viewport.YPosition = 20

	var items []list.Item

	for _, option := range exec.options {
		items = append(items, &Item{
			title: option,
		})
	}

	delegation := list.NewDefaultDelegate()
	delegation.ShowDescription = false

	exec.selectList = list.New(items, delegation, 0, 0)
	exec.selectList.Title = "Select a package to build"
	return nil
}

func (exec *Exec) Update(msg tea.Msg) (*Exec, tea.Cmd) {
	if !exec.hasSelected {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			w, h := docStyle.GetFrameSize()
			exec.selectList.SetSize(msg.Width-w, msg.Height-h-10)
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				exec.hasSelected = true
				exec.selected = exec.selectList.Index()

				exec.viewport.Width = exec.selectList.Width()
				exec.viewport.Height = exec.selectList.Height() - 10

				exec.Handler(exec.program, exec.options[exec.selected])
			}
		}

		m, cmd := exec.selectList.Update(msg)

		exec.selectList = m
		return exec, cmd
	}

	exec.done = true

	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		th := lipgloss.Height(exec.header())
		h := th + lipgloss.Height(exec.footer())

		exec.viewport.YPosition = th
		exec.viewport.Width = msg.Width
		exec.viewport.Height = msg.Height - h

		cmds = append(cmds, viewport.Sync(exec.viewport))
	case ExecData:
		exec.data = append(exec.data, msg.Data...)
	case ExecProgram:
		exec.program = msg.Program
	}

	exec.viewport.SetContent(string(exec.data))

	m, cmd := exec.viewport.Update(msg)
	cmds = append(cmds, cmd)

	exec.viewport = m

	return exec, tea.Batch(cmds...)
}

func (exec *Exec) View() string {
	if !exec.hasSelected {
		return exec.selectList.View()
	}

	return exec.header() + "\n" + exec.viewport.View() + "\n" + exec.footer()
}

func (exec *Exec) header() string {
	title := exec.Styles.Title.Render("Building kernel")
	line := strings.Repeat("─", max(0, exec.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (exec *Exec) footer() string {
	info := exec.Styles.Info.Render(fmt.Sprintf("%3.f%%", exec.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, exec.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (exec *Exec) Writer() io.Writer {
	return &Writer{program: exec.program}
}

func NewExec(options []string) *Exec {
	exec := &Exec{
		options: options,
	}

	styles := list.DefaultStyles()

	br := lipgloss.RoundedBorder()
	br.Right = "├"

	bl := lipgloss.RoundedBorder()
	bl.Left = "┤"

	exec.Styles.Title = styles.Title.BorderStyle(br).Padding(0, 1)
	exec.Styles.Info = lipgloss.NewStyle().BorderStyle(bl)

	return exec
}
