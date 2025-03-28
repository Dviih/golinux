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

