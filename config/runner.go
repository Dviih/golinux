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
	"context"
	"errors"
	"github.com/Dviih/golinux/util"
	"gopkg.in/yaml.v3"
	"io"
	"strings"
)

type RunnerKind int

const (
	RunnerKindCommand RunnerKind = iota
	RunnerKindQEMU
	RunnerKindKVM
)

var namedRunnerKind = map[RunnerKind]string{
	RunnerKindCommand: "command",
	RunnerKindQEMU:    "qemu",
	RunnerKindKVM:     "kvm",
}

func (kind *RunnerKind) UnmarshalYAML(node *yaml.Node) error {
	var s string

	if err := node.Decode(&s); err != nil {
		return err
	}

	switch strings.ToLower(s) {
	case namedRunnerKind[RunnerKindCommand]:
		*kind = RunnerKindCommand
	case namedRunnerKind[RunnerKindQEMU]:
		*kind = RunnerKindQEMU
	case namedRunnerKind[RunnerKindKVM]:
		*kind = RunnerKindKVM
	default:
		return errors.New("invalid RunnerKind")
	}

	return nil
}

