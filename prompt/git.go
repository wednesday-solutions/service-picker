package prompt

import (
	"fmt"

	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
)

func ShowPromptGitInit() bool {
	status, _ := utils.IsExists(fmt.Sprintf("%s/%s",
		utils.CurrentDirectory(),
		constants.DotGitFolder),
	)
	if status {
		return false
	} else {
		return true
	}
}

func PromptGitInit() {
	var p PromptInput
	p.Label = "Do you want to initialize git"
	p.GoBack = PromptHome
	response := p.PromptYesOrNoSelect()
	if response {
		err := GitInit()
		errorhandler.CheckNilErr(err)
	}
	PromptHome()
}

func GitInit() error {
	err := utils.RunCommandWithLogs("", constants.Git, constants.Init)
	if err != nil {
		return err
	}
	// create .gitignore file
	file := fmt.Sprintf("%s/%s",
		utils.CurrentDirectory(),
		constants.DotGitIgnore,
	)
	err = utils.CreateFile(file)
	errorhandler.CheckNilErr(err)

	err = WriteDotGitignoreFile(file)
	return err
}

func WriteDotGitignoreFile(file string) error {
	source := constants.NodeModules
	_, _, directories := utils.GetExistingStacksDatabasesAndDirectories()
	for _, dir := range directories {
		source = fmt.Sprintf("%s\n%s/%s", source, dir, constants.NodeModules)
	}
	source = fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n",
		source,
		constants.DotSst,
		constants.OutputsJson,
		"cdk.context.json",
		"build",
		"out",
	)
	err := utils.WriteToFile(file, source)
	return err
}
