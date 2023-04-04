package prompt

import (
	"github.com/wednesday-solutions/picky/pickyhelpers"
	"github.com/wednesday-solutions/picky/utils"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
)

func PromptCreateCD() {
	label := "Do you want to create CD file"
	response := PromptYesOrNoSelect(label)
	if response {
		directories, _ := PromptSelectExistingDirectories()
		err := CreateCD(directories)
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
