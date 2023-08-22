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
	for len(selectedOptions) == 0 {
		selectedOptions, _ = p.PromptMultiSelect()
	}
	// Uncomment the below line if we environment prompt.
	// environment := PromptEnvironment()
	environment := constants.Development
	stacks := PromptSelectExistingStacks()

	if platform == constants.GitHub {
		isCreateCD := false
		for _, option := range selectedOptions {
			if option == constants.CreateCI {

				err := pickyhelpers.CreateCI(stacks)
				errorhandler.CheckNilErr(err)

			} else if option == constants.CreateCD {

				isCreateCD = true
				err := CreateCD(stacks, environment)
				errorhandler.CheckNilErr(err)
			}
		}
		if isCreateCD {
			var backendExist, webExist bool
			for _, stackDir := range stacks {
				service := utils.FindService(stackDir)
				if service == constants.Backend {
					// Create task-definition.json if the stack is backend.
					err := pickyhelpers.CreateTaskDefinition(stackDir, environment)
					errorhandler.CheckNilErr(err)
					backendExist = true

					// Update existing env file of selected environment.
					err = pickyhelpers.UpdateEnvByEnvironment(stackDir, environment)
					errorhandler.CheckNilErr(err)

				} else if service == constants.Web {
					webExist = true
				}
			}
			fmt.Printf("%s", errorhandler.DoneMessage)
			PrintGitHubSecretsInfo(backendExist, webExist)
		}
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

func PrintGitHubSecretsInfo(backendExist, webExist bool) {
	err := utils.PrintInfoMessage("Save the following config data in GitHub secrets after the deployment.")
	errorhandler.CheckNilErr(err)
	secrets, count := "\n", 1

	secretKeys := []string{
		constants.AwsRegion,
		constants.AwsAccessKeyId,
		constants.AwsSecretAccessKey,
		constants.AwsEcrRepository,
	}
	if backendExist {
		for _, key := range secretKeys {
			secrets = fmt.Sprintf("%s  %d. %s\n", secrets, count, key)
			count++
		}
	}
	if webExist {
		secrets = fmt.Sprintf("%s  %d. %s\n", secrets, count, constants.DistributionId)
	}
	fmt.Printf("%s\n", secrets)
}

func CreateCD(directories []string, environment string) error {
	for _, dirName := range directories {
		var s pickyhelpers.StackDetails
		s.DirName = dirName
		s.Service = utils.FindService(dirName)
		s.Environment = environment
		s.Stack, s.Database = utils.FindStackAndDatabase(dirName)
		err := s.CreateCDFile()
		if err != nil {
			if err.Error() != errorhandler.ErrExist.Error() {
				errorhandler.CheckNilErr(err)
			}
		}
	}
	return nil
}
