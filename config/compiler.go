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

