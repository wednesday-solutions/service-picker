package prompt

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
	"github.com/wednesday-solutions/picky/pickyhelpers"
)

// PromptSetupInfra is the prompt for the setup infra option of Home prompt.
func PromptSetupInfra() {
	var p PromptInput
	p.Label = "Do you want to setup infrastructure for your project"
	p.GoBack = PromptHome
	response := p.PromptYesOrNoSelect()
	if response {
		cloudProvider := PromptCloudProvider()
		environment := PromptEnvironment()
		response := true
		// stack will get infra stack files
		stacks := utils.GetExistingInfraStacks()
		if len(stacks) > 0 {
			// change infra stack files into stack directory files.
			stacks = utils.FindStackDirectoriesByConfigStacks(stacks)
			message := "Infra setup is already exist for the following stacks,\n\n"
			for i, stack := range stacks {
				message = fmt.Sprintf("%s %d. %s\n", message, i+1, stack)
			}
			fmt.Printf("%s\n", message)
			p.Label = "Do you want to change the existing stacks"
			p.GoBack = PromptSetupInfra
			response = p.PromptYesOrNoSelect()
		}
		if response {
			stacks = PromptSelectExistingStacks()
			err := CreateInfra(stacks, cloudProvider, environment)
			errorhandler.CheckNilErr(err)
		}
		if ShowPromptGitInit() {
			p.Label = "Do you want to initialize git"
			response = p.PromptYesOrNoSelect()
			if response {
				err := GitInit()
				errorhandler.CheckNilErr(err)
			}
		}
		err := PromptDeployAfterInfra(stacks, environment)
		errorhandler.CheckNilErr(err)
	}
	PromptHome()
}

// PromptCloudProvider is a prompt for selecting a cloud provider.
func PromptCloudProvider() string {
	var p PromptInput
	p.Label = "Choose a cloud provider"
	p.Items = []string{constants.AWS}
	p.GoBack = PromptHome
	cp, _ := p.PromptSelect()
	return cp
}

// PromptEnvironment is a prompt for selecting an environment.
func PromptEnvironment() string {
	var p PromptInput
	p.Label = "Choose an environment"
	p.Items = []string{constants.Development, constants.QA, constants.Production}
	p.GoBack = PromptHome
	env, _ := p.PromptSelect()
	return env
}

// CreateInfra execute all the functionalities of infra setup.
func CreateInfra(directories []string, cloudProvider string, environment string) error {
	switch cloudProvider {
	case constants.AWS:
		status := pickyhelpers.IsInfraFilesExist()

		var err error
		if !status {
			err = pickyhelpers.CreateInfraSetup()
			errorhandler.CheckNilErr(err)
		}
		var response bool
		var s pickyhelpers.StackDetails
		for _, dirName := range directories {
			var infra pickyhelpers.Infra
			infra.Service = utils.FindService(dirName)
			infra.Stack, infra.Database = utils.FindStackAndDatabase(dirName)
			infra.DirName = dirName
			infra.CamelCaseDirName = strcase.ToCamel(dirName)
			infra.Environment = environment
			infra.ForceCreate = false

			s.Stack = infra.Stack
			s.Database = infra.Database
			s.Environment = environment
			s.StackInfo = s.GetStackInfo()

			err = infra.CreateInfraStack()
			if err != nil {
				if err.Error() == errorhandler.ErrExist.Error() {
					response = PromptAlreadyExist(dirName)
					if response {
						infra.ForceCreate = true
						err = infra.CreateInfraStack()
						errorhandler.CheckNilErr(err)
					}
				} else {
					errorhandler.CheckNilErr(err)
				}
			}
			if infra.Service == constants.Backend {
				err = pickyhelpers.UpdateEnvByEnvironment(dirName, environment)
				errorhandler.CheckNilErr(err)
			}
		}
		err = pickyhelpers.CreateSstConfigFile(s.StackInfo, directories)
		errorhandler.CheckNilErr(err)

		fmt.Printf("\n%s %s", "Generating", errorhandler.CompleteMessage)
	default:
		fmt.Printf("\nWork in Progress. Please stay tuned..!\n")
	}
	return nil
}

// PromptCreateInfraStacksWhenDeploy will setup the infra of stacks which are not already set up.
func PromptCreateInfraStacksWhenDeploy(directories []string, environment string) error {
	count := 0
	for {
		var p PromptInput
		p.Label = "Do you want to setup infra for newly selected stacks"
		p.GoBack = PromptDeploy
		response := p.PromptYesOrNoSelect()
		if response {
			var err error
			status := pickyhelpers.IsInfraFilesExist()
			if !status {
				err = pickyhelpers.CreateInfraSetup()
				errorhandler.CheckNilErr(err)
			}
			var infra pickyhelpers.Infra
			for _, dirName := range directories {
				infra.DirName = dirName
				infra.CamelCaseDirName = strcase.ToCamel(dirName)
				infra.Service = utils.FindService(dirName)
				infra.Stack, infra.Database = utils.FindStackAndDatabase(dirName)
				infra.Environment = environment
				infra.ForceCreate = false
				err = infra.CreateInfraStack()
				if err != nil {
					if err.Error() == errorhandler.ErrExist.Error() {
						response := PromptAlreadyExist(infra.DirName)
						if response {
							infra.ForceCreate = true
							err = infra.CreateInfraStack()
							errorhandler.CheckNilErr(err)
						}
					}
					errorhandler.CheckNilErr(err)
				}
			}
			fmt.Println(errorhandler.DoneMessage)
			break
		}
		count++
		if count > 1 {
			break
		}
	}
	return nil
}
