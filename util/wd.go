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
	"os"
	"path"
)

var wd string

func init() {
	get, err := os.Getwd()
	if err != nil {
		return
	}

	wd = get
}

func WD(paths ...interface{}) string {
	if len(paths) == 0 {
		return wd
	}

	pathsString := []string{wd}

	for _, v := range paths {
		pathsString = append(pathsString, v.(string))
	}

	return path.Join(pathsString...)
}

func SetWD(s string) {
	if s[0] == '/' {
		wd = s
		return
	}

	wd = path.Join(wd, s)
}

func WDProject(project string, paths ...interface{}) string {
	return WD(wdAppend(".golinux", project, paths)...)
}

func WDInitramfs(project string, paths ...interface{}) string {
	return WDProject(project, wdAppend("initramfs", paths)...)
}

func WDKernel(project, kernel string, paths ...interface{}) string {
	return WDProject(project, wdAppend("kernel", kernel, paths)...)
}

func wdAppend(v ...interface{}) []interface{} {
	var ret []interface{}

	for _, v := range v {
		switch v := v.(type) {
		case string:
			ret = append(ret, v)
		case []string:
			for _, v := range v {
				ret = append(ret, v)
			}
		case []interface{}:
			ret = append(ret, v...)
		default:
			panic("wdAppend: invalid input")
		}
	}

	return ret
}
