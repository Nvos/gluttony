package main

import (
	"context"
	"fmt"
	"github.com/alecthomas/kong"
	"gluttony/cmd/admin"
	"gluttony/cmd/run"
	"gluttony/internal/config"
	"gluttony/internal/user"
)

type Admin struct {
	Migrate    AdminMigrateCommand    `cmd:"" help:"Run database migrations."`
	Seed       AdminSeedCommand       `cmd:"" help:"Seed database with sample data."`
	CreateUser AdminCreateUserCommand `cmd:"" help:"Create user."`
}

type Globals struct {
	Config string `help:"Path to config file." default:"config.toml" short:"c"`

	cfg *config.Config `kong:"-"`
	sec *config.Secret `kong:"-"`
}

type CLI struct {
	Globals

	Admin Admin      `cmd:"" help:"Admin commands."`
	Run   RunCommand `cmd:"" help:"Run Gluttony application."`
}

type RunCommand struct {
}

type AdminCreateUserCommand struct {
	Username string `arg:"" help:"User username."`
	Password string `arg:"" help:"User password."`
	Role     string `arg:"" help:"User role." enum:"admin,user" default:"user"`
}

type AdminMigrateCommand struct {
}

type AdminSeedCommand struct {
}

func (c *AdminMigrateCommand) Run(cli *CLI) error {
	if err := admin.RunMigrations(context.Background(), cli.cfg, cli.sec); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}

func (c *AdminSeedCommand) Run(cli *CLI) error {
	if err := admin.RunSeed(context.Background(), cli.cfg, cli.sec); err != nil {
		return fmt.Errorf("run seed: %w", err)
	}

	return nil
}

func (c *AdminCreateUserCommand) Run(cli *CLI) error {
	role, err := user.NewRole(cli.Admin.CreateUser.Role)
	if err != nil {
		return fmt.Errorf("parse role: %w", err)
	}

	err = admin.AddUser(
		context.Background(),
		cli.cfg,
		cli.sec,
		cli.Admin.CreateUser.Username,
		cli.Admin.CreateUser.Password,
		role,
	)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (c *RunCommand) Run(cli *CLI) error {
	if err := run.Run(context.Background(), cli.cfg, cli.sec); err != nil {
		return fmt.Errorf("run: %w", err)
	}

	return nil
}

func main() {
	//nolint:exhaustruct // zero values are fine, no need to fill anything
	cli := CLI{}

	ctx := kong.Parse(
		&cli,
		kong.Name("gluttony"),
		kong.ConfigureHelp(kong.HelpOptions{Compact: true}),
		kong.WithAfterApply(func(ctx *kong.Context) error {
			cfg, err := config.NewConfig(cli.Config)
			if err != nil {
				return fmt.Errorf("parse config: %w", err)
			}

			sec, err := config.NewSecret()
			if err != nil {
				return fmt.Errorf("parse secret: %w", err)
			}

			cli.cfg = cfg
			cli.sec = sec

			return nil
		}),
	)

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
