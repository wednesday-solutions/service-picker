package pickyhelpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
)

// DeleteDotGitFolder deleted .git folder from stack folder.
func DeleteDotGitFolder(dirName string) error {

	path := fmt.Sprintf("%s/%s/%s", utils.CurrentDirectory(),
		dirName,
		constants.DotGitFolder,
	)
	status, _ := utils.IsExists(path)
	if status {
		err := utils.RemoveAll(path)
		errorhandler.CheckNilErr(err)
	}
	return nil
}
