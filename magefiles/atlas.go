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
