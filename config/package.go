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
