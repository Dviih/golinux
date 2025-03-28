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
