package cmd

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
)

var TestCommand = TestCommandFn()

func RunTest(*cobra.Command, []string) error {

	var firstName string
	var region string

	app := tview.NewApplication()
	form := tview.NewForm()
	form.AddInputField("Project Infra Name", "", 20, nil, func(text string) {
		firstName = text
	})
	form.AddDropDown("Select Infra Region", []string{"india", "usa"}, 0, func(option string, optionIndex int) {
		region = option
	})

	form.AddButton("Save", func() {
		app.Stop()
		fmt.Println("First Name: ", firstName, "\nInfra Region: ", region)
	})
	form.AddButton("Quit", func() {
		app.Stop()
	})

	form.SetBorder(true).SetTitle("Enter infra details").SetTitleAlign(tview.AlignLeft)
	if err := app.SetRoot(form, true).EnableMouse(true).Run(); err != nil {
		errorhandler.CheckNilErr(err)
	}

	return nil
}

func TestCommandFn() *cobra.Command {
	testCommand := &cobra.Command{
		Use:   "test",
		Short: "This command is for testing",
		Long:  "This command is for testing",
		RunE:  RunTest,
	}
	return testCommand
}

func init() {
	RootCmd.AddCommand(TestCommand)
}
