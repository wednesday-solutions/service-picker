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

func (i *InitInfo) PromptGetDirectoryName() string {
	var p PromptInput
	suffix := utils.GetSuffixOfStack(i.Stack, i.Database)
	exampleLabel := fmt.Sprintf("('-%s' suffix will be added). Eg: test-%s ", suffix, suffix)
	p.Label = fmt.Sprintf("Please enter a name for the '%s' stack %s", i.Stack, exampleLabel)
	p.GoBack = PromptSelectService
	i.DirName = p.PromptGetInput()
	i.DirName = utils.CreateStackDirectory(i.DirName, i.Stack, i.Database)
	status := true
	var err error
	for status {
		status, err = utils.IsExists(filepath.Join(utils.CurrentDirectory(), i.DirName))
		errorhandler.CheckNilErr(err)
		if status {
			p.Label = "Entered name already exists. Please enter another name"
			i.DirName = p.PromptGetInput()
			i.DirName = utils.CreateStackDirectory(i.DirName, i.Stack, i.Database)
		}
	}
	return i.DirName
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

// PromptSelectExistingStacks is a prompt function will ask for selecting available stacks.
func PromptSelectExistingStacks() []string {
	var p PromptInput
	p.Label = "Select available stacks"
	p.GoBack = PromptHome
	_, _, directories := utils.ExistingStacksDatabasesAndDirectories()
	p.Items = directories
	if len(directories) > 1 {
		p.Items = append(p.Items, "All")
	}
	var results []string
	var responses []int
	count := 0
	for {
		if len(responses) == 0 {
			results, responses = p.PromptMultiSelect()
		} else {
			break
		}
		count++
		if count > 2 {
			PromptHome()
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
func (i *InitInfo) GetDetailsTemplatesOfStacks() string {
	details := fmt.Sprintf(`
  -------- %s --------
{{ "Name:" | faint }}       {{ .Name }}
{{ "Language:" | faint }}   {{ .Language }}
{{ "Framework:" | faint }}  {{ .Framework }}
{{ "Type:" | faint }}       {{ .Type }}`, i.Service)

	if i.Service == constants.Backend {
		details = fmt.Sprintf(`%s
{{ "Databases:" | faint }}  {{ .Databases }}
`, details)
	}
	return details
}

func PromptAlreadyExist(existingFile string) bool {
	var p PromptInput
	p.Label = fmt.Sprintf("'%s' already exists, do you want to rewrite it", existingFile)
	p.GoBack = PromptHome
	return p.PromptYesOrNoSelect()
}

func PromptAlertMessage() {
	err := utils.PrintWarningMessage("Click Ctrl+C to exit.")
	errorhandler.CheckNilErr(err)
	PromptHome()
}
