package main

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Sqlc mg.Namespace

func (Sqlc) Install() error {
	goInstall := sh.RunCmd("go", "install")
	if err := goInstall("github.com/sqlc-dev/sqlc/cmd/sqlc@latest"); err != nil {
		return fmt.Errorf("go install sqlc: %w", err)
	}

	return nil
}

func (Sqlc) Codegen() error {
	buf := sh.RunCmd("sqlc")
	if err := buf("generate"); err != nil {
		return fmt.Errorf("sqlc generate: %w", err)
	}

	return nil
}
