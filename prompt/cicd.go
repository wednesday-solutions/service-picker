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
	environment := PromptEnvironment()
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
		fmt.Printf("%s", errorhandler.DoneMessage)
		if isCreateCD {

			for _, stackDir := range stacks {
				service := utils.FindService(stackDir)
				if service == constants.Backend {
					// Create task-definition.json if the stack is backend.
					err := pickyhelpers.CreateTaskDefinition(stackDir, environment)
					errorhandler.CheckNilErr(err)
				}
			}
			PrintGitHubSecretsInfo()
		}
	}
	PromptHome()
}

func (p PromptInput) PromptPlatform() string {
	p.Label = "Choose a platform"
	p.Items = []string{constants.GitHub}
	p.GoBack = PromptHome
	return p.PromptSelect()
}

func PrintGitHubSecretsInfo() {
	err := utils.PrintInfoMessage("Save the following config data in GitHub secrets after the deployment.")
	errorhandler.CheckNilErr(err)
	secrets := fmt.Sprintf("\n  %d. %s\n  %d. %s\n  %d. %s\n  %d. %s\n\n",
		1, "AWS_ACCESS_KEY_ID",
		2, "AWS_SECRET_ACCESS_KEY",
		3, "AWS_REGION",
		4, "AWS_ECR_REPOSITORY",
	)
	fmt.Printf("%s", secrets)
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
