package prompt

import (
	"fmt"
	"strings"

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

type Label struct {
	InvalidPrefix string
	ValidPrefix   string
	Question      string
	Suffix        string
}

func (p PromptInput) PromptYesOrNoConfirm() bool {
	var l Label
	l.Question = p.Label
	l.InvalidPrefix = constants.IconQuestion
	l.ValidPrefix = constants.IconSelect
	l.Suffix = "[y/N]"

	validateFn := func(input string) error {
		if len(input) < 1 {
			return fmt.Errorf("Please enter 'y' or 'n'")
		}
		return nil
	}
	templates := &promptui.PromptTemplates{
		Valid:   "{{ .ValidPrefix | bold }} {{ .Question | bold }} {{ .Suffix | faint }} ",
		Invalid: "{{ .InvalidPrefix | bold }} {{ .Question | bold }} {{ .Suffix | faint }} ",
		Success: "{{ .Question | faint }} {{ .Suffix | faint }} ",
	}
	// Refer official doc of promptui for prompt label templates.
	prompt := promptui.Prompt{
		Label:     l,
		Templates: templates,
		Validate:  validateFn,
	}
	for {
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
		result = strings.ToLower(result)
		if result == "y" || result == "yes" {
			return true
		} else if result == "n" || result == "no" {
			return false
		} else {
			err := utils.PrintWarningMessage("Please enter 'y' or 'n'.")
			errorhandler.CheckNilErr(err)
		}
	}
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

// PromptStack is prompt for selecting stack. It will come up after user selecting the service.
func (i *InitInfo) PromptStack() string {
	stacksWithDetails := utils.GetStackDetails(i.Service)
	templates := &promptui.SelectTemplates{
		Active:   fmt.Sprintf("%s {{ .Name | magenta | underline }}", constants.IconChoose),
		Inactive: "{{ .Name }}",
		Selected: fmt.Sprintf("%s {{ .Name | cyan }}", constants.IconSelect),
		Details:  i.GetDetailsTemplatesOfStacks(),
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
