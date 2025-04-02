/*
 *     Execute binaries on bare Linux.
 *     Copyright (C) 2025  Dviih
 *
 *     This program is free software: you can redistribute it and/or modify
 *     it under the terms of the GNU Affero General Public License as published
 *     by the Free Software Foundation, either version 3 of the License, or
 *     (at your option) any later version.
 *
 *     This program is distributed in the hope that it will be useful,
 *     but WITHOUT ANY WARRANTY; without even the implied warranty of
 *     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *     GNU Affero General Public License for more details.
 *
 *     You should have received a copy of the GNU Affero General Public License
 *     along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */

package main

import (
	"fmt"
	"github.com/Dviih/golinux/config"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"reflect"
	"strconv"
)

type State int

const (
	StateCompilers State = iota
	StateKernels
	StatePackages
	StateRunners
	StateLast
)

type HelpMap struct {
	Tab     key.Binding
	Help    key.Binding
	Config  key.Binding
	Sync    key.Binding
	Quit    key.Binding
	Zoom    key.Binding
	Rename  key.Binding
	Build   key.Binding
	Execute key.Binding
}

func (help *HelpMap) ShortHelp() []key.Binding {
	return []key.Binding{help.Help, help.Quit}
}

func (help *HelpMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{help.Tab, help.Help, help.Config, help.Sync, help.Quit, help.Zoom, help.Rename},
	}
}

var helpMap = &HelpMap{
	Tab:     key.NewBinding(key.WithKeys(tea.KeyTab.String()), key.WithHelp("tab", "switch state")),
	Help:    key.NewBinding(key.WithKeys(tea.KeyCtrlH.String()), key.WithHelp("control + h(elp)", "show help")),
	Config:  key.NewBinding(key.WithKeys(tea.KeyCtrlC.String()), key.WithHelp("control + c(onfig)", "show configuration")),
	Sync:    key.NewBinding(key.WithKeys(tea.KeyCtrlS.String()), key.WithHelp("control + s(ync)", "sync current state to config")),
	Quit:    key.NewBinding(key.WithKeys(tea.KeyCtrlQ.String(), tea.KeyEsc.String()), key.WithHelp("control + q(uit)", "quit golinux")),
	Zoom:    key.NewBinding(key.WithKeys(tea.KeyCtrlZ.String()), key.WithHelp("control + z(oom))", "focus into a tab")),
	Rename:  key.NewBinding(key.WithKeys(tea.KeyCtrlR.String()), key.WithHelp("control + r(ename)", "rename selected item")),
	Build:   key.NewBinding(key.WithKeys(tea.KeyCtrlB.String()), key.WithHelp("control + b(uild)", "build a package")),
	Execute: key.NewBinding(key.WithKeys(tea.KeyCtrlE.String()), key.WithHelp("control + e(execute)", "execute a package")),
}

type Main struct {
	state  State
	models map[State]tea.Model
	config *config.Config

	bindings *HelpMap
	help     help.Model
	exec     *Exec
	style    lipgloss.Style

	size  tea.WindowSizeMsg
	focus bool
}

func (main Main) Init() tea.Cmd {
	var cmds []tea.Cmd

	for _, model := range main.models {
		cmds = append(cmds, model.Init())
	}

	return tea.Batch(cmds...)
}

func (main Main) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		main.size = msg
		main.help.Width = msg.Width

		for _, model := range main.models {
			model.Update(msg)
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, main.bindings.Tab):
			main.state++

			if main.state >= StateLast {
				main.state = 0
			}

			return main, nil
		case key.Matches(msg, main.bindings.Help):
			main.help.ShowAll = true
			main.focus = true
		case key.Matches(msg, main.bindings.Config):
		case key.Matches(msg, main.bindings.Sync):
		case key.Matches(msg, main.bindings.Quit):
			if main.help.ShowAll {
				main.help.ShowAll = false
				main.focus = false
				return main, nil
			}

			return main, tea.Quit
		case key.Matches(msg, main.bindings.Zoom):
			main.focus = !main.focus
			return main, nil
		case key.Matches(msg, main.bindings.Build):
			if main.exec != nil {
				if main.exec.done {
					main.exec = nil
				}

				return main, nil
			}

			var packages []string

			for name := range main.config.Packages {
				packages = append(packages, name)
			}

			main.exec = NewExec(packages)
			main.exec.Handler = func(program *tea.Program, s string) func() {
				return func() {

				}
			}

			main.exec.Init()

			var cmd tea.Cmd
			main.exec, cmd = main.exec.Update(main.size)

			return main, cmd
		default:
		}
	}

	if main.exec != nil {
		m, cmd := main.exec.Update(msg)

		main.exec = m
		return main, cmd
	}

	state, ok := main.models[main.state]
	if !ok {
		return main, nil
	}

	m, cmds := state.Update(msg)
	main.models[main.state] = m

	return main, cmds
}

func (main Main) View() string {
	style := lipgloss.NewStyle().
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		AlignVertical(lipgloss.Left).
		AlignHorizontal(lipgloss.Center)

	if main.help.ShowAll {
		return style.AlignVertical(lipgloss.Center).Height(main.size.Height - 2).Width(main.size.Width - 2).Render(main.help.View(main.bindings))
	}

	if main.exec != nil {
		if !main.exec.hasSelected {
			view := style.UnsetBorderStyle().AlignVertical(lipgloss.Top).Render(main.exec.View()) + "\n" + "Golinux Alpha | Loaded Project: " + main.config.Project + "\n" + main.help.View(main.bindings)
			return style.Height(main.size.Height - 2).Width(main.size.Width - 2).Render(view)
		}

		return style.Height(main.size.Height - 2).Width(main.size.Width - 2).Render(main.exec.View() + "\n" + "Golinux Alpha | Loaded Project: " + main.config.Project + "\n" + main.help.View(main.bindings))
	}

	if main.focus {
		model, ok := main.models[main.state]
		if !ok {
			return "model not found"
		}

		view := style.UnsetBorderStyle().AlignVertical(lipgloss.Top).Render(model.View()) + "\n" + "Golinux Alpha | Loaded Project: " + main.config.Project + "\n" + main.help.View(main.bindings)
		return style.Height(main.size.Height - 2).Width(main.size.Width - 2).Render(view)
	}

	var views []string

	for state := State(0); state <= StateLast; state++ {
		model, ok := main.models[state]
		if !ok {
			continue
		}

		style := style

		if state == main.state {
			style = style.BorderForeground(lipgloss.Color("#3C3C3C"))
		}

		views = append(views, style.Height(main.size.Height-8).Width(main.size.Width/(len(main.models)+1)).Render(model.View()))
	}

	view := lipgloss.JoinHorizontal(lipgloss.Top, views...) + "\n" + "Golinux Alpha | Loaded Project: " + main.config.Project + "\n" + main.help.View(main.bindings)
	return style.Height(main.size.Height - 2).Width(main.size.Width - 2).Render(view)
}

