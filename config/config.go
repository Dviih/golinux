package config

import (
	"gopkg.in/yaml.v3"
	"io/fs"
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

type Config struct {
	file fs.File `yaml:"-"`

}

func FromFile(file fs.File) (*Config, error) {
	config := &Config{file: file}

	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}


	go func() {
		return
		if err := config.update(); err != nil {
			panic(err)
		}
	}()

	return config, nil
}
