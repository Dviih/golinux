package config

import (
	"context"
	"io"
)

type Package struct {
	name     string    `yaml:"-"`
	compiler *Compiler `yaml:"-"`
	Target   string    `yaml:"target"`
	Path     string    `yaml:"path"`
	Compiler string    `yaml:"compiler"`
}

func (pkg *Package) Name() string {
	return pkg.name
}

func (pkg *Package) Build(ctx context.Context, writer io.Writer) error {
	return pkg.compiler.Compile(ctx, writer, pkg.Name())
}
