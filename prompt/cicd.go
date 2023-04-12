package prompt

import (
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
	"github.com/wednesday-solutions/picky/pickyhelpers"
)

func PromptCICD() {
	var p PromptInput
	p.Label = "Select an option"
	p.Items = []string{constants.CreateCI, constants.CreateCD}
	p.GoBack = PromptHome
	selectedOptions, _ := p.PromptMultiSelect()
	services := PromptSelectExistingStacks()

	for _, option := range selectedOptions {
		if option == constants.CreateCI {

			err := pickyhelpers.CreateCI(services)
			errorhandler.CheckNilErr(err)

		} else if option == constants.CreateCD {

			err := CreateCD(services)
			errorhandler.CheckNilErr(err)
		}
	}
	PromptHome()
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
