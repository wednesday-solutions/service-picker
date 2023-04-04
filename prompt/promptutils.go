package prompt

import (
	"fmt"
	"path/filepath"

	"github.com/wednesday-solutions/picky/utils"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func AllServices() []string {
	return []string{constants.Web, constants.Mobile, constants.Backend}
}

func AllStacksOfService(service string) []string {
	var items []string
	switch service {
	case constants.Web:
		items = []string{constants.ReactJS, constants.NextJS}
	case constants.Backend:
		items = []string{constants.NodeHapiTemplate,
			constants.NodeExpressGraphqlTemplate,
			constants.NodeExpressTemplate,
			constants.GolangEchoTemplate,
		}
	case constants.Mobile:
		items = []string{constants.ReactNative,
			constants.Android,
			constants.IOS,
			constants.Flutter,
		}
	default:
		errorhandler.CheckNilErr(fmt.Errorf("Selected stack is invalid"))
	}
	return items
}

func GetDirectoryName(stack, database string) string {
	suffix := "(will add suffix)"
	label := fmt.Sprintf("Please enter a name for the '%s'%s", stack, suffix)
	dirName := PromptGetInput(label)
	dirName = utils.DirectoryName(dirName, stack, database)
	status := true
	var err error
	for status {
		status, err = fileutils.IsExists(filepath.Join(fileutils.CurrentDirectory(), dirName))
		errorhandler.CheckNilErr(err)
		if status {
			label = "Entered name already exists. Please enter another name"
			dirName = PromptGetInput(label)
			dirName = utils.DirectoryName(dirName, stack, database)
		}
	}
	return dirName
}

func PromptExit() {
	response := PromptConfirm()
	if response {
		Exit()
	} else {
		PromptHome()
	}
}

func PromptConfirm() bool {
	label := "Are you sure?"
	return PromptYesOrNoSelect(label)
}

func Exit() {
	errorhandler.CheckNilErr(errorhandler.ErrInterrupt)
}

func PromptSelectExistingDirectories() ([]string, bool) {
	label := "Select existing service"
	_, _, directories := utils.ExistingStacksDatabasesAndDirectories()
	items := directories
	items = append(items, "All")
	response := PromptSelect(label, items)
	if response == "All" {
		return directories, true
	} else {
		return []string{response}, false
	}
}
