package prompt

import (
	"fmt"

	"github.com/wednesday-solutions/picky/pickyhelpers"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
)

func PromptDeploy() {
	label := "Do you want to deploy your project"
	response := PromptYesOrNoSelect(label)
	if response {
		_ = PromptCloudProvider()
		services, _ := PromptSelectExistingServices()
		nonExistingServices := pickyhelpers.IsInfraStacksExist(services)
		if len(nonExistingServices) > 0 {
			message := "Didn't setup Infrastructure for the following services,\n\n"
			for i, service := range nonExistingServices {
				message = fmt.Sprintf("%s %d. %s\n", message, i+1, service)
			}
			fmt.Printf("%s\n", message)
			fmt.Printf("Please setup infrastructure first%s\n", errorhandler.Exclamation)
		} else {
			_ = PromptEnvironment()
			fmt.Println("Work in progress. Please stay tuned..!")
		}
	}
	PromptHome()
}

func PromptDeployAfterInfra(services []string) {
	label := "Do you want to deploy your project"
	response := PromptYesOrNoSelect(label)
	if response {
		_ = PromptCloudProvider()
		if len(services) > 0 {
			message := "Selected services are,\n\n"
			for i, service := range services {
				message = fmt.Sprintf("%s  %d. %s\n", message, i+1, service)
			}
			fmt.Printf("%s\n", message)
			label := "Are you sure to deploy the above services"
			response = PromptYesOrNoSelect(label)
			if response {
				_ = PromptEnvironment()
				fmt.Println("Work in progress. Please stay tuned..!")
			}
		} else {
			fmt.Printf("\nPlease setup infrastructure first%s\n", errorhandler.Exclamation)
		}
	}
	PromptHome()
}
