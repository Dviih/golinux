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

	commands = map[string]func(context.Context, *config.Config) error{
		"version": func(ctx context.Context, _ *config.Config) error {
			log.InfoContext(ctx, "version requested", slog.String("version", version))
			return nil
		},
		"project": func(ctx context.Context, config *config.Config) error {
			log.InfoContext(ctx, "project requested", slog.String("project", config.Project))
			return nil
		},
		"defaults": func(ctx context.Context, config *config.Config) error {
			log.InfoContext(ctx, "defaults result",
				slog.String("project", config.Project),
				slog.String("kernel", config.UseKernel),
				slog.String("package", config.DefaultPackage),
			)

			return nil
		},
		"compilers": func(ctx context.Context, config *config.Config) error {
			log.InfoContext(ctx, "compilers result",
				slog.String("project", config.Project),
				slog.String("kernel", config.UseKernel),
				slog.String("package", config.DefaultPackage),
				slog.Any("compilers", Keys(config.Compilers)),
			)

			return nil
		},
		"kernels": func(ctx context.Context, config *config.Config) error {
			log.InfoContext(ctx, "kernels result",
				slog.String("project", config.Project),
				slog.String("kernel", config.UseKernel),
				slog.String("package", config.DefaultPackage),
				slog.Any("kernels", Keys(config.Kernels)),
			)

			return nil
		},
		"packages": func(ctx context.Context, config *config.Config) error {
			log.InfoContext(ctx, "packages result",
				slog.String("project", config.Project),
				slog.String("kernel", config.UseKernel),
				slog.String("package", config.DefaultPackage),
				slog.Any("packages", Keys(config.Packages)),
			)

			return nil
		},
		"runners": func(ctx context.Context, config *config.Config) error {
			log.InfoContext(ctx, "runners result",
				slog.String("project", config.Project),
				slog.String("kernel", config.UseKernel),
				slog.String("package", config.DefaultPackage),
				slog.Any("runners", Keys(config.Runners)),
			)

			return nil
		},
	}
)
