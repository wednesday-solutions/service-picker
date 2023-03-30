package prompt

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
)

func PromptSelect(label string, items []string) string {

	templates := &promptui.SelectTemplates{
		Active:   "{{ . | green }}",
		Selected: "{{ . | cyan }}",
	}

	prompt := promptui.Select{
		Label:     label,
		Items:     items,
		Templates: templates,
		IsVimMode: false,
		Pointer:   promptui.DefaultCursor,
		Size:      constants.SizeOfPromptSelect,
	}

	_, result, err := prompt.Run()
	if err != nil {
		if err.Error() == errorhandler.ErrInterrupt.Error() {
			err = errorhandler.ExitMessage
		}
		errorhandler.CheckNilErr(err)
	}
	return result
}

func PromptGetInput(label string) string {

	validate := func(input string) error {
		if len(input) <= 1 {
			return fmt.Errorf("Length should be greater than 1%s\n", errorhandler.Exclamation)
		}
		return nil
	}
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }}",
		Valid:   "{{ . | green }}",
		Invalid: "{{ . | red }}",
	}
	prompt := promptui.Prompt{
		Label:     label,
		Validate:  validate,
		Templates: templates,
		IsVimMode: true,
	}
	result, err := prompt.Run()
	if err != nil {
		if err.Error() == errorhandler.ErrInterrupt.Error() {
			err = errorhandler.ExitMessage
		}
		errorhandler.CheckNilErr(err)
	}
	return result
}

func PromptYesOrNoSelect(label string) bool {
	items := []string{constants.Yes, constants.No}

	response := PromptSelect(label, items)
	if response == constants.Yes {
		return true
	} else {
		return false
	}
}

func SelectOne(label string, items []string) string {
	return PromptSelect(label, items)
}
