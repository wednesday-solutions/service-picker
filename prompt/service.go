package prompt

import "github.com/wednesday-solutions/picky/internal/constants"

func PromptSelectService() {
	var p PromptInput
	p.Label = "Pick a service"
	p.Items = AllServices()
	p.GoBack = PromptAlertMessage

	var i InitInfo
	i.Service = p.PromptSelect()
	i.PromptSelectStack()
}

func (i *InitInfo) PromptSelectStack() {

	i.Stack = i.PromptStack()
	if i.Service == constants.Backend {
		i.PromptSelectStackDatabase()
	}
	i.DirName = i.PromptGetDirectoryName()
	i.PromptSelectInit()
}
