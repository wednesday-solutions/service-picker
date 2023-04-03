package prompt

import (
	"fmt"
	"path/filepath"

	"github.com/wednesday-solutions/picky/pickyhelpers"
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

func PromptGetDirectoryName(stack, database string) string {
	var p PromptInput
	suffix := utils.GetSuffixOfStack(stack, database)
	suffix = fmt.Sprintf("('-%s' suffix will be added)", suffix)
	p.Label = fmt.Sprintf("Please enter a name for the '%s stack'%s", stack, suffix)
	p.GoBack = PromptSelectService
	dirName := p.PromptGetInput()
	dirName = utils.DirectoryName(dirName, stack, database)
	status := true
	var err error
	for status {
		status, err = fileutils.IsExists(filepath.Join(fileutils.CurrentDirectory(), dirName))
		errorhandler.CheckNilErr(err)
		if status {
			p.Label = "Entered name already exists. Please enter another name"
			dirName = p.PromptGetInput()
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
	var p PromptInput
	p.Label = "Are you sure"
	p.GoBack = PromptHome
	return p.PromptYesOrNoSelect()
}

func Exit() {
	errorhandler.CheckNilErr(errorhandler.ErrInterrupt)
}

func PromptSelectExistingStacks() []string {
	var p PromptInput
	p.Label = "Select existing stacks"
	p.GoBack = PromptHome
	_, _, directories := utils.ExistingStacksDatabasesAndDirectories()
	p.Items = directories
	p.Items = append(p.Items, "All")
	results, responses := p.PromptMultiSelect()
	for _, respIdx := range responses {
		if respIdx == len(p.Items)-1 {
			return directories
		}
	}
	return results
}

func PromptRunYarnInstall() error {
	var p PromptInput
	p.Label = "Can we run 'yarn install'"
	p.GoBack = PromptHome
	count := 0
	for {
		response := p.PromptYesOrNoSelect()
		count++
		if response {
			err := pickyhelpers.RunYarnInstall()
			return err
		}
		if count == 2 {
			PromptHome()
		}
		err := utils.PrintWarningMessage("You can't deploy without installing 'yarn'")
		errorhandler.CheckNilErr(err)
	}
}
