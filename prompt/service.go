package prompt

import (
	"github.com/wednesday-solutions/picky/internal/constants"
)

var services = []string{
	constants.Web,
	// constants.Mobile,
	constants.Backend,
}

func PromptSelectService() {
	var i InitInfo
	var p PromptInput
	p.Label = "Pick a service"
	p.Items = services
	p.GoBack = PromptAlertMessage

	var index int
	i.Service, index = p.PromptSelect()
	services = ResetItems(p.Items, &index)
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
