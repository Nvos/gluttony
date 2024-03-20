package main

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Buf mg.Namespace

func (Buf) Codegen() error {
	buf := sh.RunCmd("buf")
	if err := buf("generate"); err != nil {
		return fmt.Errorf("buf generate: %w", err)
	}

	return nil
}

func (Buf) Install() error {
	goInstall := sh.RunCmd("go", "install")
	npmInstall := sh.RunCmd("npm", "install", "-g")

	goModules := []string{
		"github.com/bufbuild/buf/cmd/buf@latest",
		"github.com/fullstorydev/grpcurl/cmd/grpcurl@latest",
		"google.golang.org/protobuf/cmd/protoc-gen-go@latest",
		"connectrpc.com/connect/cmd/protoc-gen-connect-go@latest",
	}
	for i := range goModules {
		if err := goInstall(goModules[i]); err != nil {
			return fmt.Errorf("buf install go module=%s: %w", goModules[i], err)
		}
	}

	npmModules := []string{
		"@connectrpc/protoc-gen-connect-es",
		"@bufbuild/protoc-gen-es",
	}
	for i := range npmModules {
		if err := npmInstall(npmModules[i]); err != nil {
			return fmt.Errorf("buf install npm package=%s: %w", npmModules[i], err)
		}
	}

	return nil
}
