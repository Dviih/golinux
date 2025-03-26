package config

import (
	"errors"
	"github.com/Dviih/golinux/util"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
	"io"
)

func (config *Config) update() error {
	reader, ok := config.file.(io.Reader)
	if !ok {
		return errors.New("unsupported: io.Reader")
	}

	seeker, ok := config.file.(io.Seeker)
	if !ok {
		return errors.New("unsupported: io.Seeker")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	stat, err := config.file.Stat()
	if err != nil {
		return err
	}

	if err = watcher.Add(util.WD(stat.Name())); err != nil {
		return err
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op != fsnotify.Write {
				continue
			}

			if _, err = seeker.Seek(0, io.SeekStart); err != nil {
				return err
			}

			if err = yaml.NewDecoder(reader).Decode(&config); err != nil {
				return err
			}
		case err = <-watcher.Errors:
			return err
		}
	}
}
