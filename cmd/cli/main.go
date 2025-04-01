package main

import (
	"context"
	"flag"
	"github.com/Dviih/golinux/config"
	"github.com/Dviih/golinux/util"
	"log/slog"
	"os"
	"path"
	"strings"
)

var (
	log = slog.Default()

	WD         string
	ConfigPath string
	version    string

)
