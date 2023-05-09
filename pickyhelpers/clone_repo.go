package pickyhelpers

import (
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
)

type StackDetails struct {
	Service     string
	Stack       string
	DirName     string
	CurrentDir  string
	Database    string
	Environment string
	StackInfo   map[string]interface{}
}

func (s StackDetails) CloneRepo() error {

	// Download the selected stack.
	err := utils.RunCommandWithoutLogs("", "git", "clone", constants.Repos()[s.Stack], s.DirName)
	errorhandler.CheckNilErr(err)

	// Delete cd.yml file from the cloned repo.
	cdFilePatch := s.CurrentDir + "/" + s.DirName + constants.CDFilePathURL
	status, _ := utils.IsExists(cdFilePatch)
	if status {
		err = utils.RemoveFile(cdFilePatch)
		errorhandler.CheckNilErr(err)
	}
	return nil
}
