package main

import (
	"os"
	"strings"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/blang/semver/v4"
	"github.com/ecumenos/fxecumenos"
	"github.com/ecumenos/fxecumenos/fxlogger"
	postgresfx "github.com/ecumenos/fxecumenos/fxpostgres"
	"github.com/ecumenos/fxecumenos/fxpostgres/migrations"
	"github.com/ecumenos/fxecumenos/fxrf"
	"github.com/ecumenos/fxecumenos/zerodowntime"
	"github.com/ecumenos/orbis-socius-register/cmd/api/configuration"
	"github.com/ecumenos/orbis-socius-register/cmd/api/httpserver"
	cli "github.com/urfave/cli/v2"
	"golang.org/x/exp/slog"
)

var Version = semver.MustParse("0.0.1")

func main() {
	if err := run(os.Args); err != nil {
		slog.Error("exiting", "err", err)
		os.Exit(-1)
	}
}

func run(args []string) error {
	app := cli.App{
		Name:    "api",
		Usage:   "serving API",
		Version: strings.Join(Version.Build, "."),
	}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "prod",
			Value:   false,
			EnvVars: []string{"PROD"},
		},
	}

	app.Commands = []*cli.Command{
		runAppCmd,
		migrateUpCmd,
		migrateDownCmd,
	}

	return app.Run(args)
}

var runAppCmd = &cli.Command{
	Name:  "run-api-server",
	Usage: "run API HTTP server",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		app := fx.New(
			fx.Supply(fxecumenos.ServiceName("api")),
			fx.Supply(fxecumenos.Version(Version)),
			fxlogger.Module,
			configuration.Module,
			httpserver.Module,
			postgresfx.Module,
			fxrf.Module,
		)

		return zerodowntime.HandleApp(app)
	},
}

var migrateUpCmd = &cli.Command{
	Name:  "migrate-up",
	Usage: "run migrations up",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		app := fx.New(
			fx.Supply(fxecumenos.ServiceName("api")),
			fx.Supply(fxecumenos.Version(Version)),
			fxlogger.Module,
			configuration.Module,
			fx.Invoke(func(cfg *postgresfx.Config, logger *zap.Logger, shutdowner fx.Shutdowner) error {
				fn := migrations.NewMigrateUpFunc()
				if !cctx.Bool("prod") {
					logger.Info("runnning migrate up",
						zap.String("db_url", cfg.URL),
						zap.String("source_path", cfg.MigrationsPath))
				}
				return fn(cfg.MigrationsPath, cfg.URL+"?sslmode=disable", logger, shutdowner)
			}),
		)

		return zerodowntime.HandleApp(app)
	},
}

var migrateDownCmd = &cli.Command{
	Name:  "migrate-down",
	Usage: "run migrations down",
	Flags: []cli.Flag{},
	Action: func(cctx *cli.Context) error {
		app := fx.New(
			fx.Supply(fxecumenos.ServiceName("api")),
			fx.Supply(fxecumenos.Version(Version)),
			fxlogger.Module,
			configuration.Module,
			fx.Invoke(func(cfg *postgresfx.Config, logger *zap.Logger, shutdowner fx.Shutdowner) error {
				fn := migrations.NewMigrateDownFunc()
				if !cctx.Bool("prod") {
					logger.Info("runnning migrate down",
						zap.String("db_url", cfg.URL),
						zap.String("source_path", cfg.MigrationsPath))
				}
				return fn(cfg.MigrationsPath, cfg.URL+"?sslmode=disable", logger, shutdowner)
			}),
		)

		return zerodowntime.HandleApp(app)
	},
}
