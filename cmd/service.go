package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/prompt"
)

// ServiceCmd is the command variable of ServiceSelection.
var ServiceSelection = ServiceSelectionFn()

func RunService(*cobra.Command, []string) {
	selectedService := prompt.PromptSelect("Pick a service", []string{constants.WEB, constants.MOBILE, constants.BACKEND})

	switch selectedService {

	case constants.WEB:
		prompt.PromptSelectStack(constants.WEB, []string{constants.REACT, constants.NEXT})

	case constants.BACKEND:
		prompt.PromptSelectStack(constants.BACKEND, []string{constants.NODE_HAPI, constants.NODE_EXPRESS, constants.NODE_EXPRESS_TS, constants.GOLANG})

	case constants.MOBILE:

	default:

	}
}

// ServiceSelectionFn represents the ServiceSelection command
func ServiceSelectionFn() *cobra.Command {

	var ServiceSelection = &cobra.Command{
		Use:   constants.SERVICE,
		Short: "Pick a Service",
		Long: `Pick a service for your:

		1. Web
		2. Mobile
		3. Backend

		from the list of @wednesday-solutions's open source projects.
`,
		Run: RunService,
	}
	return ServiceSelection
}

func init() {

	RootCmd.AddCommand(ServiceSelection)

	ServiceSelection.Flags().BoolP("help", "h", false, "Help for service selection")
}
