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
	return err
}
