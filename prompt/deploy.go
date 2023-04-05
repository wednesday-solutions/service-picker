package prompt

import (
	"fmt"

	"github.com/wednesday-solutions/picky/pickyhelpers"
	"github.com/wednesday-solutions/picky/utils"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
)

func PromptDeploy() {
	var p PromptInput
	p.Label = "Do you want to deploy your project"
	p.GoBack = PromptHome
	response := p.PromptYesOrNoSelect()
	if response {
		_ = PromptCloudProvider()
		configStacks := pickyhelpers.SstConfigExistStacks()
		if len(configStacks) > 0 {
			message := "Selected stacks are,\n\n"
			for i, stack := range configStacks {
				message = fmt.Sprintf("%s  %d.%s\n", message, i+1, stack)
			}
			fmt.Printf("%s\n", message)
		}
		var stacks []string
		for {
			stacks = PromptSelectExistingStacks()
			if len(stacks) > 0 {
				break
			}
		}
		nonExistingStacks := pickyhelpers.IsInfraStacksExist(stacks)
		if len(nonExistingStacks) > 0 {
			message := "Didn't setup Infrastructure for the following stacks,\n\n"
			for i, stack := range nonExistingStacks {
				message = fmt.Sprintf("%s %d. %s\n", message, i+1, stack)
			}
			fmt.Printf("%s\n", message)
			fmt.Printf("Please setup infrastructure first%s\n", errorhandler.Exclamation)
		} else {
			_ = PromptEnvironment()
			err := PromptRunYarnInstall()
			errorhandler.CheckNilErr(err)
			fmt.Println("Work in progress. Please stay tuned..!")
		}
	}
	PromptHome()
}

func PromptDeployAfterInfra(services []string) {
	var p PromptInput
	p.Label = "Do you want to deploy your project"
	p.GoBack = PromptHome
	response := p.PromptYesOrNoSelect()
	if response {
		_ = PromptCloudProvider()
		if len(services) > 0 {
			message := "Selected stacks are,\n\n"
			for i, service := range services {
				message = fmt.Sprintf("%s  %d. %s\n", message, i+1, service)
			}
			fmt.Printf("%s\n", message)
			p.Label = "Are you sure to deploy the above stacks"
			response = p.PromptYesOrNoSelect()
			if response {
				_ = PromptEnvironment()
				err := PromptRunYarnInstall()
				errorhandler.CheckNilErr(err)
				fmt.Println("Work in progress. Please stay tuned..!")
			}
		} else {
			err := utils.PrintWarningMessage("Please setup infrastructure first.")
			errorhandler.CheckNilErr(err)
		}
	}
	PromptHome()
}
