package pickyhelpers

import (
	"fmt"
	"path/filepath"

	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/utils"
)

func InstallDependencies(pkgManager string, path ...string) error {
	filePath := filepath.Join(path...)
	err := utils.RunCommandWithLogs(filePath, pkgManager, "install")
	return err
}

func BuildSST(pkgManager string) error {
	err := utils.RunCommandWithLogs("", "yarn", "build")
	return err
}

func DeploySST(pkgManager, environment string) error {
	environment = utils.GetShortEnvName(environment)
	arg := fmt.Sprintf("deploy:%s", environment)
	err := utils.RunCommandWithLogs("", pkgManager, arg)
	return err
}

func RemoveDeploy(pkgManager, environment string) error {
	environment = utils.GetShortEnvName(environment)
	arg := fmt.Sprintf("remove:%s", environment)
	err := utils.RunCommandWithLogs("", pkgManager, "run", arg)
	return err
}

func ParseDeployOutputs() error {
	err := utils.RunCommandWithoutLogs("", "node", constants.ParseSstOutputs)
	return err
}
