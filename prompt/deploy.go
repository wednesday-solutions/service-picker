package prompt

import "fmt"

func PromptDeploy() {
	label := "Do you want to deploy your project"
	response := PromptYesOrNoSelect(label)
	if response {
		_ = PromptCloudProvider()
	}
	fmt.Println("Work in progress. Please stay tuned..!")
	PromptHome()
}
