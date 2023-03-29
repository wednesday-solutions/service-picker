package pickyhelpers

import (
	"fmt"
	"strings"

	"github.com/stoewer/go-strcase"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func GetStackInfo(stack, database string) map[string]interface{} {

	var webDir, mobileDir, backendDir string
	currentDir := fileutils.CurrentDirectory()
	splitDirs := strings.Split(currentDir, "/")
	projectName := splitDirs[len(splitDirs)-1]
	projectName = strcase.SnakeCase(projectName)

	dirNames, err := fileutils.ReadAllContents(currentDir)
	errorhandler.CheckNilErr(err)
	var splitDirName []string
	var service string
	for _, dirName := range dirNames {
		splitDirName = strings.Split(dirName, "-")
		if len(splitDirName) > 0 {
			service = splitDirName[len(splitDirName)-1]
			switch service {
			case constants.Web:
				webDir = dirName
			case constants.Mobile:
				mobileDir = dirName
			case constants.Backend:
				backendDir = dirName
			}
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

	return stackInfo
}
