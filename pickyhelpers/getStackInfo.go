package pickyhelpers

import (
	"strings"

	"github.com/stoewer/go-strcase"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func GetStackInfo(stack, database string) map[string]interface{} {

	currentDir := fileutils.CurrentDirectory()
	splitDirs := strings.Split(currentDir, "/")
	projectName := splitDirs[len(splitDirs)-1]
	projectName = strcase.SnakeCase(projectName)

	stackDestination := map[string]string{
		constants.WebStatus:     currentDir + "/" + constants.Web,
		constants.MobileStatus:  currentDir + "/" + constants.Mobile,
		constants.BackendStatus: currentDir + "/" + constants.Backend,
	}
	stackInfo := make(map[string]interface{})

	for status, destination := range stackDestination {
		stackInfo[status], _ = fileutils.IsExists(destination)
	}
	stackInfo[constants.Stack] = stack
	stackInfo[constants.Database] = database
	stackInfo[constants.ProjectName] = projectName

	return stackInfo
}
