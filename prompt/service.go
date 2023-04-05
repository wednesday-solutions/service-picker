package prompt

import "github.com/wednesday-solutions/picky/utils/constants"

func PromptSelectService() {
	var p PromptInput
	p.Label = "Pick a service"
	p.Items = AllServices()
	p.GoBack = PromptHome
	service := p.PromptSelect()
	PromptSelectStack(service)
}

func PromptSelectStack(service string) {

	stack := PromptStack(service)

	var database string
	if service == constants.Backend {
		database = PromptSelectStackDatabase(service, stack)
	}
	dirName := PromptGetDirectoryName(stack, database)
	PromptSelectInit(service, stack, database, dirName)
}
