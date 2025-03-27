package config

import (
	"bytes"
	"context"
	"io"
	"os"
	"path"
)

type Kernel struct {
	name     string    `yaml:"-"`
	compiler *Compiler `yaml:"-"`

	Path     string `yaml:"path"`
	Config   string `yaml:"config"`
	Compiler string `yaml:"compiler"`
}

