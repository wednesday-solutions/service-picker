package prompt

import (
	"fmt"

	"github.com/iancoleman/strcase"
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
		services, all := PromptSelectExistingServices()
		environment := PromptEnvironment()
		err := CreateInfra(services, cloudProvider, all, environment)
		errorhandler.CheckNilErr(err)
		PromptDeployAfterInfra(services)
	}
	PromptHome()
}

func PromptCloudProvider() string {
	label := "Choose a cloud provider"
	items := []string{constants.AWS}
	return PromptSelect(label, items)
}

func PromptEnvironment() string {
	label := "Choose an environment"
	items := []string{constants.Development, constants.QA, constants.Production}
	return PromptSelect(label, items)
}

func CreateInfra(directories []string, cloudProvider string, all bool, environment string) error {
	if cloudProvider == constants.AWS {
		status := pickyhelpers.IsInfraFilesExist()
		var (
			stack, database string
			stackInfo       map[string]interface{}
			err             error
		)
		done := make(chan bool)
		go pickyhelpers.ProgressBar(20, "Generating", done)

		if !status {
			err = pickyhelpers.CreateInfraSetup()
			errorhandler.CheckNilErr(err)
		}
		var camelCaseDirName string
		for _, dirName := range directories {
			service := utils.FindService(dirName)
			stack, database = utils.FindStackAndDatabase(dirName)
			stackInfo = pickyhelpers.GetStackInfo(stack, database, environment)

			camelCaseDirName = strcase.ToCamel(dirName)
			err = pickyhelpers.CreateInfraStacks(service, stack, database, dirName, environment)
			if err != nil {
				if err.Error() != errorhandler.ErrExist.Error() {
					errorhandler.CheckNilErr(err)
				}
			}
			if !all {
				err = pickyhelpers.CreateSstConfigFile(stackInfo, all, camelCaseDirName, directories)
				errorhandler.CheckNilErr(err)
			}
			if service == constants.Backend {
				err = pickyhelpers.UpdateEnvDevelopment(dirName, environment)
				errorhandler.CheckNilErr(err)
			}
		}
		if all {
			err = pickyhelpers.CreateSstConfigFile(stackInfo, all, constants.All, directories)
			errorhandler.CheckNilErr(err)
		}
		<-done
		fmt.Printf("\n%s %s", "Generating", errorhandler.CompleteMessage)
	}
	return nil
}
