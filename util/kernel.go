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

package util

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

const KernelMirror = "https://cdn.kernel.org/pub/linux/kernel/v%s.x/linux-%s.tar.gz"

func GetKernel(project, target, version string) error {
	switch target[0] {
	case 'v':
		target = path.Join(WD(), ".golinux", project, "kernel", target[1:])
	case '/':
		break
	default:
		target = path.Join(WD(), ".golinux", project, "kernel", target)
	}

	ctx, cancel := context.WithCancelCause(context.Background())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(KernelMirror, version[:strings.IndexByte(version, '.')], version), nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			cancel(err)
		}
	}(res.Body)

	greader, err := gzip.NewReader(res.Body)
	if err != nil {
		return err
	}

	defer func(reader *gzip.Reader) {
		if err = reader.Close(); err != nil {
			cancel(err)
		}
	}(greader)

	treader := tar.NewReader(greader)

}
