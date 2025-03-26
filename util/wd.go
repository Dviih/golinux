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

	return path.Join(wdAppend(wd, paths)...)
}

func SetWD(s string) {
	if s[0] == '/' {
		wd = s
		return
	}

	wd = path.Join(wd, s)
}

func WDProject(project string, paths ...interface{}) string {
	return WD(".golinux", project, paths)
}

func WDInitramfs(project string, paths ...interface{}) string {
	return WDProject(project, "initramfs", paths)
}

func WDKernel(project string, paths ...interface{}) string {
	return WDProject(project, "kernel", paths)
}

