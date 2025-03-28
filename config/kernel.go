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

package config

import (
	"context"
	"github.com/Dviih/golinux/util"
	"io"
	"os"
	"os/exec"
)

type Kernel struct {
	name     string    `yaml:"-"`
	compiler *Compiler `yaml:"-"`

	Path     string `yaml:"path"`
	Config   string `yaml:"config"`
	Compiler string `yaml:"compiler"`
}

func (kernel *Kernel) Name() string {
	return kernel.name
}

func (kernel *Kernel) Menu(ctx context.Context) error {
	compiler := &Compiler{
		name:        "menuconfig",
		project:     kernel.compiler.project,
		Call:        "make menuconfig",
		Environment: kernel.compiler.Environment,
		Arguments:   kernel.compiler.Arguments,
	}

	if err := compiler.compile(ctx, os.Stdin, os.Stdout, os.Stderr, kernel.Path); err != nil {
		return err
	}

	return nil
}

var configMap = map[string]string{
	"CONFIG_DEFAULT_HOSTNAME": "golinux",
}

func (kernel *Kernel) config(ctx context.Context) error {
	configMap := configMap
	configMap["CONFIG_INITRAMFS_SOURCE"] = util.WDInitramfs(kernel.compiler.project)

	for property, value := range configMap {
		var cmd *exec.Cmd

		if value == "" {
			cmd = exec.CommandContext(ctx, "./scripts/config", "--enable", property)
		} else {
			cmd = exec.CommandContext(ctx, "./scripts/config", "--set-str", property, value)
		}

		cmd.Dir = util.WDKernel(kernel.compiler.project, kernel.Name())

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

func (kernel *Kernel) Build(ctx context.Context, writer io.Writer) error {
	if err := kernel.config(ctx); err != nil {
		return err
	}

	return kernel.compiler.Compile(ctx, writer, kernel.Path)
}
