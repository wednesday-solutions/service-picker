package prompt

import (
	"fmt"

	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
)

func PromptSelectStackDatabase(service, stack, dirName string) {
	database := PromptSelectDatabase(service, stack)
	PromptSelectStackConfig(service, stack, database, dirName)
}

func PromptSelectDatabase(service, stack string) string {
	label := "Choose a database"
	var database string
	var items []string

	if service == constants.Backend {
		switch stack {
		case constants.NodeHapiTemplate, constants.NodeExpressGraphqlTemplate, constants.GolangEchoTemplate:
			items = []string{constants.PostgreSQL, constants.MySQL}
		case constants.NodeExpressTemplate:
			items = []string{constants.MongoDB}
		default:
			errorhandler.CheckNilErr(fmt.Errorf("\nSelected stack is invalid%s\n", errorhandler.Exclamation))
		}
	} else {
		switch stack {
		case constants.ReactJS, constants.NextJS:
			items = []string{constants.PostgreSQL, constants.MySQL, constants.MongoDB}
		case constants.ReactNativeTemplate, constants.AndroidTemplate,
			constants.IOSTemplate, constants.FlutterTemplate:
			items = []string{constants.PostgreSQL, constants.MySQL, constants.MongoDB}
		default:
			errorhandler.CheckNilErr(fmt.Errorf("\nSelected stack is invalid%s\n", errorhandler.Exclamation))
		}
	}
	database = PromptSelect(label, items)
	return database
}

func PromptAllDatabases() string {
	label := "Please select initialized database"
	items := []string{constants.PostgreSQL, constants.MySQL, constants.MongoDB}
	return PromptSelect(label, items)
}
