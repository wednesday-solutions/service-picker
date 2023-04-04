package prompt

import (
	"fmt"

	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
)

func PromptSelectStackDatabase(service, stack string) string {
	database := SelectDatabase(stack)
	return database
}

func SelectDatabase(stack string) string {
	var p PromptInput
	p.Label = "Choose a database"
	p.GoBack = PromptSelectService
	var database string
	switch stack {
	case constants.NodeHapiTemplate, constants.NodeExpressGraphqlTemplate, constants.GolangEchoTemplate:
		p.Items = []string{constants.PostgreSQL, constants.MySQL}
	case constants.NodeExpressTemplate:
		p.Items = []string{constants.MongoDB}
	default:
		errorhandler.CheckNilErr(fmt.Errorf("\nSelected stack is invalid%s\n", errorhandler.Exclamation))
	}
	database = p.PromptSelect()
	return database
}

func PromptAllDatabases() string {
	var p PromptInput
	p.Label = "Choose a database"
	p.GoBack = PromptSelectService
	p.Items = []string{constants.PostgreSQL, constants.MySQL, constants.MongoDB}
	return p.PromptSelect()
}
