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
	kernel.Path = util.WDKernel(config.Project)

	return kernel
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
