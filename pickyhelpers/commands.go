package pickyhelpers

import (
	"fmt"
	"path/filepath"

	"github.com/wednesday-solutions/picky/internal/utils"
)

func InstallDependencies(pkgManager string, path ...string) error {
	filePath := filepath.Join(path...)
	err := utils.RunCommandWithLogs(filePath, pkgManager, "install")
	return err
}

func BuildSST() error {
	err := utils.RunCommandWithLogs("", "yarn", "build")
	return err
}

func DeploySST(environment string) error {
	environment = utils.GetEnvironment(environment)
	arg := fmt.Sprintf("deploy:%s", environment)
	err := utils.RunCommandWithLogs("", "yarn", arg)
	return err
}
