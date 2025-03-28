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

func (kernel *Kernel) Build(ctx context.Context, writer io.Writer) error {
	return kernel.compiler.Compile(ctx, writer, kernel.Path)
}
