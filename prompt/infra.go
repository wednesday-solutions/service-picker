package prompt

import (
	"fmt"

	"github.com/wednesday-solutions/picky/pickyhelpers"
	"github.com/wednesday-solutions/picky/utils"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
)

func PromptSetupInfra() {
	label := "Do you want to setup infrastructure for your project"
	response := PromptYesOrNoSelect(label)
	if response {
		cloudProvider := PromptCloudProvider()
		directories := PromptSelectExistingDirectories()
		err := CreateInfra(directories, cloudProvider)
		errorhandler.CheckNilErr(err)
	}
	PromptHome()
}

func PromptCloudProvider() string {
	label := "Choose a cloud provider"
	items := []string{constants.AWS}
	return PromptSelect(label, items)
}

func CreateInfra(directories []string, cloudProvider string) error {
	if cloudProvider == constants.AWS {
		status := pickyhelpers.IsInfraFilesExist()
		var (
			stack, database string
			stackInfo       map[string]interface{}
			err             error
		)
		done := make(chan bool)
		go pickyhelpers.ProgressBar(30, "Generating", done)

		for _, dirName := range directories {
			service := utils.FindService(dirName)
			stack, database = utils.FindStackAndDatabase(dirName)
			stackInfo = pickyhelpers.GetStackInfo(stack, database)
			if !status {
				err = pickyhelpers.CreateInfraSetup(stackInfo)
				errorhandler.CheckNilErr(err)
			} else {
				err = pickyhelpers.CreateInfraStacks(service, stack, database, dirName, stackInfo)
				if err != nil {
					if err.Error() != errorhandler.ErrExist.Error() {
						errorhandler.CheckNilErr(err)
					}
				}
			}
			if service == constants.Backend {
				err = pickyhelpers.UpdateEnvDevelopment(dirName)
				errorhandler.CheckNilErr(err)
			}
		}
		<-done
		fmt.Printf("\n%s %s", "Generating", errorhandler.CompleteMessage)
	}
	return nil
}
