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
	arguments := strings.Split(compiler.Call, " ")[1:]

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
