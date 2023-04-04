package prompt

import "github.com/wednesday-solutions/picky/utils/constants"

func PromptSelectService() {
	label := "Pick a service"
	services := AllServices()
	service := PromptSelect(label, services)
	PromptSelectStack(service)
}

func PromptSelectStack(service string) {
	label := "Pick a stack"
	items := AllStacksOfService(service)
	stack := PromptSelect(label, items)

	var database string
	if service == constants.Backend {
		database = PromptSelectStackDatabase(service, stack)
	}
	dirName := GetDirectoryName(stack, database)
	PromptSelectInit(service, stack, database, dirName)
}
