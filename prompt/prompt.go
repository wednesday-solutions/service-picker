package prompt

import (
	"fmt"

	"github.com/spaceweasel/promptui"
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
)

type PromptInput struct {
	Label  string
	Items  []string
	GoBack func()
}

func (p PromptInput) PromptSelect() string {
	templates := &promptui.SelectTemplates{
		Active:   fmt.Sprintf("%s {{ . | magenta | underline }}", constants.IconChoose),
		Selected: fmt.Sprintf("%s {{ . | cyan }}", constants.IconSelect),
	}
	prompt := promptui.Select{
		Label:     p.Label,
		Items:     p.Items,
		Templates: templates,
		IsVimMode: false,
		Pointer:   promptui.DefaultCursor,
		Size:      constants.SizeOfPromptSelect,
	}
	_, result, err := prompt.Run()
	if err != nil {
		if err.Error() == errorhandler.ErrInterrupt.Error() {
			err = errorhandler.ExitMessage
		} else if err == promptui.ErrEOF {
			p.GoBack()
			fmt.Printf("\nSomething error happened in GoBack.\n")
		}
		errorhandler.CheckNilErr(err)
	}
	return result
}

func (p PromptInput) PromptGetInput() string {

	validate := func(input string) error {
		if len(input) <= 1 {
			return fmt.Errorf("Length should be greater than 1%s\n", errorhandler.Exclamation)
		}
		return nil
	}
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }}",
		Valid:   fmt.Sprintf("%s {{ . | green }}", constants.IconSelect),
		Invalid: fmt.Sprintf("%s {{ . | red }}", constants.IconWrong),
	}
	prompt := promptui.Prompt{
		Label:     p.Label,
		Validate:  validate,
		Templates: templates,
		IsVimMode: true,
	}
	result, err := prompt.Run()
	if err != nil {
		if err.Error() == errorhandler.ErrInterrupt.Error() {
			err = errorhandler.ExitMessage
		} else if err == promptui.ErrEOF {
			p.GoBack()
			fmt.Printf("\nSomething error happened in GoBack.\n")
		}
		errorhandler.CheckNilErr(err)
	}
	return result
}

func (p PromptInput) PromptYesOrNoSelect() bool {
	p.Items = []string{constants.Yes, constants.No}

	response := p.PromptSelect()
	if response == constants.Yes {
		return true
	} else {
		return false
	}
}

func (p PromptInput) PromptMultiSelect() ([]string, []int) {
	templates := &promptui.MultiSelectTemplates{
		Selected: fmt.Sprintf("%s {{ . | green }}", constants.IconSelect),
	}
	prompt := promptui.MultiSelect{
		Label:     p.Label,
		Items:     p.Items,
		Templates: templates,
		Size:      constants.SizeOfPromptSelect,
	}
	results, err := prompt.Run()
	if err != nil {
		if err.Error() == errorhandler.ErrInterrupt.Error() {
			err = errorhandler.ExitMessage
		} else if err == promptui.ErrEOF {
			p.GoBack()
			fmt.Printf("\nSomething error happened in GoBack.\n")
		}
		errorhandler.CheckNilErr(err)
		return nil, nil
	}
	selected := []string{}
	for _, result := range results {
		selected = append(selected, p.Items[result])
	}
	err = utils.PrintMultiSelectMessage(selected)
	errorhandler.CheckNilErr(err)
	return selected, results
}

func PromptStack(service string) string {

	stacksWithDetails := utils.GetStackDetails(service)
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

	templates := &promptui.SelectTemplates{
		Active:   fmt.Sprintf("%s {{ .Name | magenta | underline }}", constants.IconChoose),
		Inactive: "{{ .Name }}",
		Selected: fmt.Sprintf("%s {{ .Name | cyan }}", constants.IconSelect),
		Details:  details,
	}

	prompt := promptui.Select{
		Label:     "Pick a stack",
		Items:     stacksWithDetails,
		Templates: templates,
		Size:      constants.SizeOfPromptSelect,
	}
	idx, _, err := prompt.Run()
	if err != nil {
		if err.Error() == errorhandler.ErrInterrupt.Error() {
			err = errorhandler.ExitMessage
		} else if err == promptui.ErrEOF {
			PromptSelectService()
			fmt.Printf("\nSomething error happened in GoBack.\n")
		}
		errorhandler.CheckNilErr(err)
	}
	return stacksWithDetails[idx].Name
}
