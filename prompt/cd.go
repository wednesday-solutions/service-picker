package prompt

import (
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
	"github.com/wednesday-solutions/picky/pickyhelpers"
)

func PromptCreateCD() {
	var p PromptInput
	p.Label = "Do you want to create CD file"
	p.GoBack = PromptHome
	response := p.PromptYesOrNoSelect()
	if response {
		services := PromptSelectExistingStacks()
		err := CreateCD(services)
		errorhandler.CheckNilErr(err)
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
