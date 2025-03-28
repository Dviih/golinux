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
	"errors"
	"github.com/Dviih/golinux/util"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"os"
)

type KV struct {
	Key   string
	Value string
}

type KVS []*KV

func (kvs *KVS) UnmarshalYAML(value *yaml.Node) error {
	current := &KV{}

	for _, content := range value.Content {
		for _, content := range content.Content {
			if current.Key != "" {
				current.Value = content.Value

				*kvs = append(*kvs, current)
				current = &KV{}

				continue
			}

			current.Key = content.Value
		}
	}

	return nil
}

func (kv *KV) MarshalYAML() (interface{}, error) {
	return map[string]string{kv.Key: kv.Value}, nil
}

type Config struct {
	file fs.File `yaml:"-"`

	Project   string               `yaml:"project"`
	Compilers map[string]*Compiler `yaml:"compilers"`
	Kernels   map[string]*Kernel   `yaml:"kernel"`
	Packages  map[string]*Package  `yaml:"packages"`
	Runners   map[string]*Runner   `yaml:"runners"`

	DefaultPackage string `yaml:"default_package"`
}

func (config *Config) Sync() error {
	seeker, ok := config.file.(io.Seeker)
	if !ok {
		return errors.New("unsupported: missing io.Seeker")
	}

	if _, err := seeker.Seek(0, io.SeekStart); err != nil {
		return err
	}

	writer, ok := config.file.(io.Writer)
	if !ok {
		return errors.New("unsupported: missing io.Writer")
	}

	return yaml.NewEncoder(writer).Encode(config)
}

func (config *Config) Close() error {
	err := config.file.Close()
	*config = Config{}

	return err
}

func (config *Config) Compiler(name string) *Compiler {
	compiler, ok := config.Compilers[name]
	if !ok {
		return &Compiler{name: name}
	}

	compiler.name = name
	compiler.project = config.Project

	return compiler
}

func (config *Config) Kernel(name string) *Kernel {
	kernel, ok := config.Kernels[name]
	if !ok {
		return &Kernel{name: name}
	}

	kernel.name = name
	kernel.compiler = config.Compiler(kernel.Compiler)
	kernel.Path = util.WDKernel(config.Project, kernel.Name())

	return kernel
}

func (config *Config) Package(name string) *Package {
	pkg, ok := config.Packages[name]
	if !ok {
		return &Package{}
	}

	pkg.name = name
	pkg.compiler = config.Compiler(pkg.Compiler)

	if pkg.Target != "" && pkg.Path == "" {
		pkg.Path = util.WD(pkg.Target)
		pkg.Target = ""
	}

	return pkg
}

func (config *Config) Runner(name string) *Runner {
	runner, ok := config.Runners[name]
	if !ok {
		return &Runner{}
	}

	runner.name = name
	runner.project = config.Project

	return runner
}

func FromPath(path string) (*Config, error) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return FromFile(file)
}

func FromFile(file fs.File) (*Config, error) {
	config := &Config{file: file}

	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	if err := config.Sync(); err != nil {
		return nil, err
	}

	go func() {
		if err := config.update(); err != nil {
			println("failed to load update config:", err.Error())
		}
	}()

	return config, nil
}
