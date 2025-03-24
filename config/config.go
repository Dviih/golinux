package config

import (
	"gopkg.in/yaml.v3"
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

