package pickyhelpers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func InstallDependencies(pkgManager string, path ...string) error {
	filePath := filepath.Join(path...)
	err := RunCommand(filePath, pkgManager, "install")
	return err
}

func BuildSST() error {
	err := RunCommand("", "yarn", "build")
	return err
}

func DeploySST(environment string) error {
	arg := fmt.Sprintf("deploy:%s", environment)
	err := RunCommand("", "yarn", arg)
	return err
}

func RunCommand(path string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if path != "" {
		cmd.Dir = path
	}
	err := cmd.Run()
	fmt.Printf("\n")
	return err
}
