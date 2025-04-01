package main

import (
	"context"
	"errors"
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
		"build": func(ctx context.Context, config *config.Config) error {
			name := flag.Arg(1)

			if name == "" {
				name = config.DefaultPackage
			}

			pkg := config.Package(name)

			log.InfoContext(ctx, "requested package build", slog.String("package", pkg.Name()))

			if err := buildPackage(ctx, config, pkg); err != nil {
				log.ErrorContext(ctx, "failed to build package",
					slog.String("package", name),
					slog.Any("error", err),
				)

				return err
			}

			kernel := config.Kernel(config.UseKernel)

			log.InfoContext(ctx, "requested kernel build", slog.String("kernel", kernel.Name()))

			if err := kernel.Build(context.Background(), nil); err != nil {
				log.ErrorContext(ctx, "failed to build kernel",
					slog.String("kernel", kernel.Name()),
					slog.Any("error", err),
				)

				return err
			}

			return nil
		},
		"run": func(ctx context.Context, config *config.Config) error {
			if flag.Arg(1) == "" {
				return errors.New("missing runner name")
			}

			return config.Runner(flag.Arg(1)).Execute(ctx, os.Stdin, os.Stdout, os.Stderr)
		},
	}
)
