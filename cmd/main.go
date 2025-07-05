package main

import (
	"flag"
	"fmt"
	"gluttony/cmd/admin"
	"gluttony/cmd/run"
	"gluttony/internal/config"
	"gluttony/internal/user"
	"golang.org/x/net/context"
	"os"
)

func main() {
	// Parse global flags
	var envFile string
	flag.StringVar(&envFile, "env-file", "", "Path to environment file")
	flag.Parse()

	if envFile != "" {
		if err := config.LoadEnvFile(envFile); err != nil {
			fmt.Fprintf(os.Stderr, "Load env file=%q: %v\n", envFile, err)
			os.Exit(1)
		}
	}

	args := flag.Args()

	if len(args) < 1 {
		printUsage()
		os.Exit(1)
	}

	ctx := context.Background()
	command := args[0]

	switch command {
	case "migrate":
		if err := admin.RunMigrations(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Migration failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migrations completed successfully")

	case "seed":
		if err := admin.RunSeed(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Seeding failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Database seeded successfully")

	case "user":
		if len(args) < 3 {
			fmt.Fprintf(os.Stderr, "Usage: %s user <subcommand>\n", args[0])
			printUserUsage()
			os.Exit(1)
		}
		handleUserCommand(ctx, args[2:])
	case "run":
		if err := run.Run(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Run failed: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func handleUserCommand(ctx context.Context, args []string) {
	if len(args) == 0 {
		printUserUsage()
		os.Exit(1)
	}

	subcommand := args[0]

	switch subcommand {
	case "add":
		if len(args) < 4 {
			fmt.Fprintf(os.Stderr, "Usage: %s user add <role> <username> <password>\n", args[0])
			fmt.Fprintf(os.Stderr, "Roles: admin, user\n")
			os.Exit(1)
		}
		rawRole := args[1]
		username := args[2]
		password := args[3]
		role, err := user.NewRole(rawRole)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Parse role: %v\n", err)
			os.Exit(1)
		}

		if err := admin.AddUser(ctx, username, password, role); err != nil {
			fmt.Fprintf(os.Stderr, "Adding user failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("User created successfully with role: %s\n", role)

	default:
		fmt.Fprintf(os.Stderr, "Unknown user subcommand: %s\n", subcommand)
		printUserUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Print(`Usage: cli [--env-file=<path>] <command>

Global flags:
  --env-file=<path>        Load environment variables from file (optional)

Commands:
  migrate                  Run database migrations
  seed                     Seed database with sample data
  user                     User management commands
  run                      Start Gluttony service
`)
}

func printUserUsage() {
	fmt.Print(`User commands:
  add <role> <username> <password>  Create user with specified role
    Roles: admin, user
`)
}