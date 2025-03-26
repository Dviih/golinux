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

