package utils

import (
	"fmt"
	"os"
	"text/template"

	"github.com/fatih/color"
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
)

// CreateMessageTemplate creates new text template for printing colorful logs.
func CreateMessageTemplate(name, text string) *template.Template {
	tpl, err := template.New(name).Parse(text)
	errorhandler.CheckNilErr(err)
	tpl = template.Must(tpl, err)
	return tpl
}

// PrintMultiSelectMessage prints multi selected options.
func PrintMultiSelectMessage(messages []string) error {
	var message, coloredMessage string
	var tpl *template.Template
	if len(messages) > constants.Zero {
		var templateText string
		if len(messages) == constants.One {
			templateText = fmt.Sprintf("%s %d option selected: {{ . }}\n",
				constants.IconSelect,
				len(messages))
		} else {
			templateText = fmt.Sprintf("%s %d options selected: {{ . }}\n",
				constants.IconSelect,
				len(messages))
		}
		for _, option := range messages {
			message = fmt.Sprintf("%s%s ", message, option)
		}
		coloredMessage = color.GreenString("%s", message)
		tpl = CreateMessageTemplate("message", templateText)
	} else {
		message = "No options selected, please select atleast one."
		coloredMessage = color.YellowString("%s", message)
		tpl = CreateMessageTemplate("responseMessage", fmt.Sprintf("%s {{ . }}\n", constants.IconWarn))
	}
	err := tpl.Execute(os.Stdout, coloredMessage)
	return err
}

// PrintWarningMessage prints given message in yellow color as warning message in terminal.
func PrintWarningMessage(message string) error {
	tpl := CreateMessageTemplate("warningMessage", fmt.Sprintf("\n%s {{ . }}\n", constants.IconWarn))
	message = color.YellowString("%s", message)
	err := tpl.Execute(os.Stdout, message)
	return err
}

// PrintInfoMessage prints given message in cyan color as info message in terminal.
func PrintInfoMessage(message string) error {
	tpl := CreateMessageTemplate("InfoMessage", fmt.Sprintf("\n%s {{ . }}\n", constants.IconChoose))
	message = color.CyanString("%s", message)
	err := tpl.Execute(os.Stdout, message)
	return err
}

// PrintErrorMessage prints the given message in red color as error message
func PrintErrorMessage(message string) error {
	tpl := CreateMessageTemplate("errorMessage", fmt.Sprintf("\n{{ . }}%s\n", errorhandler.Exclamation))
	message = color.RedString("%s", message)
	err := tpl.Execute(os.Stdout, message)
	return err
}
