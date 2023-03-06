package prompt

import (
	"github.com/manifoldco/promptui"
	"github.com/wednesday-solutions/service-picker/utils/errorhandler"
)

func PromptSelect(label string, items []string) string {

	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	_, result, err := prompt.Run()

	errorhandler.CheckNilErr(err)

	return result
}
