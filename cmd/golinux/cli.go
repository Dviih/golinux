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
	"flag"
	"github.com/Dviih/golinux/config"
	"github.com/Dviih/golinux/util"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	wd string
)

func main() {
	flag.StringVar(&wd, "wd", "", "The path to the working directory")
	flag.Parse()

	util.SetWD(wd)

	c, err := config.FromPath(util.WD("golinux.yaml"))
	if err != nil {
		panic(err)
	}

	model := NewMain(c)
	model.models = map[State]tea.Model{
		StateCompilers: NewList("Compilers", c.Compilers),
		StateKernels:   NewList("Kernels", c.Kernels),
		StatePackages:  NewList("Packages", c.Packages),
		StateRunners:   NewList("Runners", c.Runners),
	}

	if _, err := tea.NewProgram(model, tea.WithAltScreen()).Run(); err != nil {
		panic(err)
	}
}
