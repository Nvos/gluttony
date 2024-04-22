package main

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Atlas mg.Namespace

func (Atlas) Codegen(migrationName string) error {
	diff := sh.RunCmd("atlas", "migrate", "diff")

	if err := diff("--env", "local", migrationName); err != nil {
		return fmt.Errorf("atlas migrate: %w", err)
	}

	return nil
}

func (Atlas) Lint() error {
	lint := sh.RunCmd("atlas", "migrate", "lint")

	if err := lint("--env", "local", "--latest", "1"); err != nil {
		return fmt.Errorf("atlas lint: %w", err)
	}

	return nil
}
