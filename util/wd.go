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
