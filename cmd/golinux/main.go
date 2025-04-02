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
	"github.com/charmbracelet/bubbles/textarea"
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

	configAreaActive bool
	configArea       textarea.Model

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
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		main.size = msg
		main.help.Width = msg.Width

		for i, model := range main.models {
			main.models[i], cmd = model.Update(msg)
			cmds = append(cmds, cmd)
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
			main.exec, cmd = main.exec.Update(main.size)

			return main, tea.Batch(append(cmds, cmd)...)
		default:
		}
	}


	if main.exec != nil {
		main.exec, cmd = main.exec.Update(msg)
		return main, cmd
	}

	state, ok := main.models[main.state]
	if !ok {
		return main, nil
	}

	main.models[main.state], cmd = state.Update(msg)
	cmds = append(cmds, cmd)

	return main, tea.Batch(cmds...)
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

func NewMain(config *config.Config) Main {
	return Main{
		state:    0,
		models:   make(map[State]tea.Model),
		config:   config,
		bindings: helpMap,
		help:     help.New(),
		style:    lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
	}
}

func toString(value reflect.Value) string {
	switch value.Kind() {
	case reflect.Invalid:
		return "<nil>"
	case reflect.Bool:
		if value.Bool() {
			return "true"
		}

		return "false"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10)
	case reflect.Uintptr, reflect.Array, reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.Struct, reflect.UnsafePointer:
		return fmt.Sprintf("%s", value)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'g', -1, 64)
	case reflect.Complex64, reflect.Complex128:
		return strconv.FormatComplex(value.Complex(), 'g', -1, 128)
	case reflect.String:
		return value.String()
	default:
		panic("invalid value")
	}
}

func setValue(value reflect.Value, s string) {
	if !value.CanSet() {
		return
	}

	switch value.Kind() {
	case reflect.Invalid:
		return
	case reflect.Bool:
		if s == "true" {
			value.SetBool(true)
			return
		}

		value.SetBool(false)
		return
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(s, 10, value.Type().Bits())
		if err != nil {
			panic(err)
		}

		value.SetInt(i)
		return
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(s, 10, value.Type().Bits())
		if err != nil {
			panic(err)
		}

		value.SetUint(u)
		return
	case reflect.Uintptr, reflect.Array, reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.Struct, reflect.UnsafePointer:
		return // can't set
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(s, value.Type().Bits())
		if err != nil {
			panic(err)
		}

		value.SetFloat(f)
		return
	case reflect.Complex64, reflect.Complex128:
		c, err := strconv.ParseComplex(s, value.Type().Bits())
		if err != nil {
			panic(err)
		}

		value.SetComplex(c)
		return
	case reflect.String:
		value.SetString(s)
		return
	default:
		panic("invalid value")
	}
}

func rvAbs(value reflect.Value) reflect.Value {
	for value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	return value
}
