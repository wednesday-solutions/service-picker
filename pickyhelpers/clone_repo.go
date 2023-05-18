package pickyhelpers

import (
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
)

// StackDetails is the collection of all stack details.
type StackDetails struct {
	// Service refers to the service type of the stack.
	Service string

	// Stacks refers to Wednesday Solutions open source templates.
	Stack string

	// DirName refers to the name of directory of stack.
	DirName string

	// CurrentDir refers to the root directory.
	CurrentDir string

	// Database refers to the database of selected stack.
	Database string

	// Environment refers to the environment
	Environment string

	// StackInfo consist of all the details about stacks.
	StackInfo map[string]interface{}
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
