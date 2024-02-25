package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"ariga.io/atlas-go-sdk/atlasexec"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// Define the execution context, supplying a migration directory
	// and potentially an `atlas.hcl` configuration file using `atlasexec.WithHCL`.
	dir := os.DirFS(filepath.Join(pwd, "internal", "atlas", "migrations"))
	workdir, err := atlasexec.NewWorkingDir(
		atlasexec.WithMigrations(
			dir,
		),
	)
	if err != nil {
		log.Fatalf("failed to load working directory: %v", err)
	}
	// atlasexec works on a temporary directory, so we need to close it
	defer workdir.Close()

	// Initialize the client.
	client, err := atlasexec.NewClient(workdir.Path(), "atlas")
	if err != nil {
		log.Fatalf("failed to initialize client: %v", err)
	}
	// Run `atlas migrate apply` on a SQLite database under /tmp.
	res, err := client.MigrateApply(context.Background(), &atlasexec.MigrateApplyParams{
		URL: "sqlite:///tmp/demo.db?_fk=1&cache=shared",
	})
	if err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}
	fmt.Printf("Applied %d migrations\n", len(res.Applied))
}
