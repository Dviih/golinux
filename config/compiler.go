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
	"errors"
	"github.com/Dviih/golinux/util"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

type Compiler struct {
	name    string `yaml:"-"`
	project string `yaml:"-"`

	Call        string            `yaml:"call"`
	Environment map[string]string `yaml:"environment"`
	Arguments   KVS               `yaml:"arguments"`
}

func (compiler *Compiler) GetEnvironment() []string {
	env := os.Environ()

	for k, v := range compiler.Environment {
		env = append(env, k+"="+v)
	}

	return env
}

func (compiler *Compiler) GetArgs() []string {
	arguments := strings.Split(compiler.Call, " ")

	for _, argument := range compiler.Arguments {
		if argument.Value == "" {
			arguments = append(arguments, "-"+argument.Key)
			continue
		}

		arguments = append(arguments, "-"+argument.Key, argument.Value)
	}

	return arguments
}

func (compiler *Compiler) Name() string {
	return compiler.name
}

func (compiler *Compiler) Compile(ctx context.Context, writer io.Writer, target string) error {
	if writer == nil {
		writer = os.Stdout
	}

	stderr := &util.Writer{}

	switch target[0] {
	case '/':
		break
	default:
		target = path.Join(util.WD(), target)
	}

	if err := compiler.compile(ctx, nil, writer, stderr, target); err != nil {
		return stderr.Error(err)
	}

	return nil
}

func (compiler *Compiler) compile(ctx context.Context, stdin io.Reader, stdout, stderr io.Writer, target string) error {
	ctx, cancel := context.WithCancel(ctx)

	args := compiler.GetArgs()

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

	cmd.Env = compiler.GetEnvironment()
	cmd.Dir = target
	cmd.Stdin = stdin
	cmd.Stderr = stderr
	cmd.Stdout = stdout

	if err := cmd.Start(); err != nil {
		cancel()
		return err
	}

	if cmd.Err != nil {
		cancel()
		return cmd.Err
	}

	if err := cmd.Wait(); err != nil {
		cancel()
		return err
	}

	cancel()

	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.Canceled) {
			return nil
		}

		return ctx.Err()
	}
}
