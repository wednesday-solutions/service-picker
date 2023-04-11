package prompt

import (
	"fmt"

	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
)

func (i *InitInfo) PromptSelectStackDatabase() {
	i.SelectDatabase()
}

func (i *InitInfo) SelectDatabase() {
	var p PromptInput
	p.Label = "Choose a database"
	p.GoBack = PromptSelectService
	switch i.Stack {
	case constants.NodeHapiTemplate, constants.NodeExpressGraphqlTemplate, constants.GolangEchoTemplate:
		p.Items = []string{constants.PostgreSQL, constants.MySQL}
	case constants.NodeExpressTemplate:
		p.Items = []string{constants.MongoDB}
	default:
		errorhandler.CheckNilErr(fmt.Errorf("\nSelected stack is invalid%s\n", errorhandler.Exclamation))
	}
	i.Database = p.PromptSelect()
}

func PromptAllDatabases() string {
	var p PromptInput
	p.Label = "Choose a database"
	p.GoBack = PromptSelectService
	p.Items = []string{constants.PostgreSQL, constants.MySQL, constants.MongoDB}
	return p.PromptSelect()
}
