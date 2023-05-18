package pickyhelpers

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/utils"
)

func (s StackDetails) GetStackInfo() map[string]interface{} {

	var webDir, mobileDir, backendDir string
	var webDirectories, backendPgDirectories, backendMysqlDirectories []string
	currentDir := utils.CurrentDirectory()
	projectName := strcase.ToSnake(utils.GetProjectName())

	_, databases, directories := utils.GetExistingStacksDatabasesAndDirectories()
	for i, dirName := range directories {
		s.Service = utils.FindService(dirName)
		switch s.Service {
		case constants.Web:
			webDir = dirName
			webDirectories = append(webDirectories, dirName)
		case constants.Mobile:
			mobileDir = dirName
		case constants.Backend:
			backendDir = dirName
			if databases[i] == constants.PostgreSQL {
				backendPgDirectories = append(backendPgDirectories, dirName)
			} else if databases[i] == constants.MySQL {
				backendMysqlDirectories = append(backendMysqlDirectories, dirName)
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
			stackInfo[status], _ = utils.IsExists(destination)
		} else {
			stackInfo[status] = false
		}
	}
	stackInfo[constants.Stack] = s.Stack
	stackInfo[constants.Database] = s.Database
	stackInfo[constants.ProjectName] = projectName
	stackInfo[constants.WebDirName] = webDir
	stackInfo[constants.MobileDirName] = mobileDir
	stackInfo[constants.BackendDirName] = backendDir
	stackInfo[constants.ExistingDirectories] = directories
	stackInfo[constants.Environment] = s.Environment
	stackInfo[constants.WebDirectories] = webDirectories
	stackInfo[constants.BackendPgDirectories] = backendPgDirectories
	stackInfo[constants.BackendMysqlDirectories] = backendMysqlDirectories

	return stackInfo
}
