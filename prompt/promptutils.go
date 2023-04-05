package prompt

import (
	"fmt"
	"path/filepath"

	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
)

func AllServices() []string {
	return []string{constants.Web, constants.Mobile, constants.Backend}
}

func AllStacksOfService(service string) []string {
	var items []string
	switch service {
	case constants.Web:
		items = []string{constants.ReactJS, constants.NextJS}
	case constants.Backend:
		items = []string{constants.NodeHapiTemplate,
			constants.NodeExpressGraphqlTemplate,
			constants.NodeExpressTemplate,
			constants.GolangEchoTemplate,
		}
	case constants.Mobile:
		items = []string{constants.ReactNative,
			constants.Android,
			constants.IOS,
			constants.Flutter,
		}
	default:
		errorhandler.CheckNilErr(fmt.Errorf("Selected stack is invalid"))
	}
	return items
}

func PromptGetDirectoryName(stack, database string) string {
	var p PromptInput
	suffix := utils.GetSuffixOfStack(stack, database)
	exampleLabel := fmt.Sprintf("('-%s' suffix will be added). Eg: test-%s ", suffix, suffix)
	p.Label = fmt.Sprintf("Please enter a name for the '%s' stack %s", stack, exampleLabel)
	p.GoBack = PromptSelectService
	dirName := p.PromptGetInput()
	dirName = utils.DirectoryName(dirName, stack, database)
	status := true
	var err error
	for status {
		status, err = utils.IsExists(filepath.Join(utils.CurrentDirectory(), dirName))
		errorhandler.CheckNilErr(err)
		if status {
			p.Label = "Entered name already exists. Please enter another name"
			dirName = p.PromptGetInput()
			dirName = utils.DirectoryName(dirName, stack, database)
		}
	}
	return dirName
}

func PromptExit() {
	response := PromptConfirm()
	if response {
		Exit()
	} else {
		PromptHome()
	}
}

func PromptConfirm() bool {
	var p PromptInput
	p.Label = "Are you sure"
	p.GoBack = PromptHome
	return p.PromptYesOrNoSelect()
}

func Exit() {
	errorhandler.CheckNilErr(errorhandler.ErrInterrupt)
}

func PromptSelectExistingStacks() []string {
	var p PromptInput
	p.Label = "Select existing stacks"
	p.GoBack = PromptHome
	_, _, directories := utils.ExistingStacksDatabasesAndDirectories()
	p.Items = directories
	p.Items = append(p.Items, "All")
	var results []string
	var responses []int
	for {
		if len(responses) == 0 {
			results, responses = p.PromptMultiSelect()
		} else {
			break
		}
	}
	for _, respIdx := range responses {
		if respIdx == len(p.Items)-1 {
			return directories
		}
	}
	return results
}

// GetDetailsTemplatesOfStacks return the details template for the given service.
func GetDetailsTemplatesOfStacks(service string) string {
	details := fmt.Sprintf(`
-------- %s --------
{{ "Name:" | faint }}       {{ .Name }}
{{ "Language:" | faint }}   {{ .Language }}
{{ "Framework:" | faint }}  {{ .Framework }}`, service)

	if service == constants.Backend {
		details = fmt.Sprintf(`%s
{{ "Databases:" | faint }}  {{ .Databases }}
{{ "Type:" | faint }}       {{ .Type }}
`, details)
	}
	return details
}
