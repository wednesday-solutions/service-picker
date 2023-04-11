package pickyhelpers

import (
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
)

func CloneRepo(stack, dirName, path string) error {

	// Download the selected stack.
	err := utils.RunCommandWithoutLogs("", "git", "clone", constants.Repos()[stack], dirName)
	errorhandler.CheckNilErr(err)

	// Delete cd.yml file from the cloned repo.
	cdFilePatch := path + "/" + dirName + constants.CDFilePathURL
	status, _ := utils.IsExists(cdFilePatch)
	if status {
		err = utils.RemoveFile(cdFilePatch)
		errorhandler.CheckNilErr(err)
	}
	return nil
}
