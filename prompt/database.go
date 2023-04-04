package prompt

import (
	"fmt"

	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
)

func PromptSelectStackDatabase(service, stack string) string {
	database := SelectDatabase(stack)
	return database
}

func SelectDatabase(stack string) string {
	label := "Choose a database"
	var database string
	var items []string
	switch stack {
	case constants.NodeHapiTemplate, constants.NodeExpressGraphqlTemplate, constants.GolangEchoTemplate:
		items = []string{constants.PostgreSQL, constants.MySQL}
	case constants.NodeExpressTemplate:
		items = []string{constants.MongoDB}
	default:
		errorhandler.CheckNilErr(fmt.Errorf("\nSelected stack is invalid%s\n", errorhandler.Exclamation))
	}
	database = PromptSelect(label, items)
	return database
}

func PromptAllDatabases() string {
	label := "Choose a database"
	items := []string{constants.PostgreSQL, constants.MySQL, constants.MongoDB}
	return PromptSelect(label, items)
}
