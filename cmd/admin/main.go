package main

import (
	"fmt"
	"gluttony/cmd/admin/commands"
	"golang.org/x/net/context"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	ctx := context.Background()
	command := os.Args[1]

	switch command {
	case "migrate":
		if err := commands.RunMigrations(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Migration failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Migrations completed successfully")

	case "seed":
		if err := commands.RunSeed(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Seeding failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Database seeded successfully")

	case "add-admin":
		if len(os.Args) < 4 {
			fmt.Fprintf(os.Stderr, "Usage: %s add-admin <username> <password>\n", os.Args[0])
			os.Exit(1)
		}
		username := os.Args[2]
		password := os.Args[3]

		if err := commands.AddAdmin(ctx, username, password); err != nil {
			fmt.Fprintf(os.Stderr, "Adding admin failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Admin user created successfully")

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: cli <command>")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  migrate              Run database migrations")
	fmt.Println("  seed                 Seed database with sample data")
	fmt.Println("  add-admin <username> <password>  Create admin user")
}
