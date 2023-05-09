package prompt

import (
	"fmt"

	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
	"github.com/wednesday-solutions/picky/pickyhelpers"
)

func PromptCICD() {
	var p PromptInput
	platform := p.PromptPlatform()
	p.Label = "Select option"
	p.Items = []string{constants.CreateCI, constants.CreateCD}
	p.GoBack = PromptHome
	selectedOptions, _ := p.PromptMultiSelect()
	services := PromptSelectExistingStacks()

	if platform == constants.GitHub {
		for _, option := range selectedOptions {
			if option == constants.CreateCI {

				err := pickyhelpers.CreateCI(services)
				errorhandler.CheckNilErr(err)

			} else if option == constants.CreateCD {

				err := CreateCD(services)
				errorhandler.CheckNilErr(err)
			}
		}
		fmt.Printf("%s", errorhandler.DoneMessage)
	}
	PromptHome()
}

func (p PromptInput) PromptPlatform() string {
	p.Label = "Choose a platform"
	p.Items = []string{constants.GitHub}
	p.GoBack = PromptHome
	platform, _ := p.PromptSelect()
	return platform
}

func CreateCD(directories []string) error {
	for _, dirName := range directories {
		service := utils.FindService(dirName)
		stack, database := utils.FindStackAndDatabase(dirName)
		err := pickyhelpers.CreateCDFile(service, stack, database, dirName)
		if err != nil {
			if err.Error() != errorhandler.ErrExist.Error() {
				errorhandler.CheckNilErr(err)
			}
		}
	}
	return nil
}
