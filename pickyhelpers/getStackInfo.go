package pickyhelpers

import (
	"fmt"
	"strings"

	"github.com/stoewer/go-strcase"
	"github.com/wednesday-solutions/picky/utils"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func GetStackInfo(stack, database string) map[string]interface{} {

	var webDir, mobileDir, backendDir, service string
	currentDir := fileutils.CurrentDirectory()
	splitDirs := strings.Split(currentDir, "/")
	projectName := splitDirs[len(splitDirs)-1]
	projectName = strcase.SnakeCase(projectName)

	_, _, directories := utils.ExistingStacksDatabasesAndDirectories()
	for _, dirName := range directories {
		service = utils.FindService(dirName)
		switch service {
		case constants.Web:
			webDir = dirName
		case constants.Mobile:
			mobileDir = dirName
		case constants.Backend:
			backendDir = dirName
		}
	}
	stackDestination := map[string]string{
		constants.WebStatus:     currentDir + "/" + webDir,
		constants.MobileStatus:  currentDir + "/" + mobileDir,
		constants.BackendStatus: currentDir + "/" + backendDir,
	}
	stackInfo := make(map[string]interface{})

	for status, destination := range stackDestination {
		if destination != fmt.Sprintf("%s/", currentDir) {
			stackInfo[status], _ = fileutils.IsExists(destination)
		} else {
			stackInfo[status] = false
		}
	}
	stackInfo[constants.Stack] = stack
	stackInfo[constants.Database] = database
	stackInfo[constants.ProjectName] = projectName
	stackInfo[constants.WebDirName] = webDir
	stackInfo[constants.MobileDirName] = mobileDir
	stackInfo[constants.BackendDirName] = backendDir

	return stackInfo
}
