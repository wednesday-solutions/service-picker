package pickyhelpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
	"github.com/wednesday-solutions/picky/pickyhelpers/sources"
)

func CreateCI(stackDirs []string) error {
	workflowsPath := fmt.Sprintf("%s/%s", utils.CurrentDirectory(),
		constants.GithubWorkflowsDir)

	utils.CreateGithubWorkflowDir()
	var stackCIPath, stack string
	var status bool
	for _, dir := range stackDirs {
		stackCIPath = fmt.Sprintf("%s/ci-%s.yml", workflowsPath, dir)
		status, _ = utils.IsExists(stackCIPath)
		if !status {
			stack, _ = utils.FindStackAndDatabase(dir)
			err := CreateStackCI(stackCIPath, dir, stack)
			errorhandler.CheckNilErr(err)
		}
	}
	return nil
}

// CreateStackCI creates and writes CI for the given stack.
func CreateStackCI(path, stackDir, stack string) error {
	var environment, source string
	if stack == constants.NodeExpressGraphqlTemplate {
		environment = constants.Development
	} else {
		environment = constants.Dev
	}
	if stack != constants.GolangEchoTemplate {
		source = sources.CISource(stack, stackDir, environment)

		err := utils.WriteToFile(path, source)
		errorhandler.CheckNilErr(err)
	} else {
		err := utils.PrintWarningMessage(fmt.Sprintf(
			"CI of '%s' is in work in progress..!", stack,
		))
		errorhandler.CheckNilErr(err)
	}
	return nil
}
