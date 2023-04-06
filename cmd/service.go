package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/prompt"
)

var ServiceSelection = ServiceSelectionFn()

func RunService(*cobra.Command, []string) error {
	prompt.PromptHome()
	return nil
}

// ServiceSelectionFn represents the ServiceSelection command
func ServiceSelectionFn() *cobra.Command {

	var ServiceSelection = &cobra.Command{
		Use:   constants.Service,
		Short: "Pick a Service",
		Long: `Pick a service for your:

		1. Web
		2. Mobile
		3. Backend

		from the list of @wednesday-solutions's open source projects.
`,
		RunE: RunService,
	}
	return ServiceSelection
}

func init() {

	RootCmd.AddCommand(ServiceSelection)

	ServiceSelection.Flags().BoolP("help", "h", false, "Help for service selection")
}
