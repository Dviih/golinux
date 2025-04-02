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

