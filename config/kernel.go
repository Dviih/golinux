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
	return kernel.compiler.Compile(ctx, writer, kernel.Path)
}
