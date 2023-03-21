package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/prompt"
)

// ServiceCmd is the command variable of ServiceSelection.
var ServiceSelection = ServiceSelectionFn()

func RunService(*cobra.Command, []string) error {
	selectedService := prompt.PromptSelect("Pick a service", []string{constants.Web,
		constants.Mobile,
		constants.Backend},
	)

	switch selectedService {

	case constants.Web:
		prompt.PromptSelectStack(constants.Web, []string{constants.ReactJS, constants.NextJS})

	case constants.Backend:
		prompt.PromptSelectStack(constants.Backend, []string{constants.NodeHapiTemplate,
			constants.NodeExpressGraphqlTemplate,
			constants.NodeExpressTemplate,
			constants.GolangEchoTemplate},
		)

	case constants.Mobile:
		prompt.PromptSelectStack(constants.Mobile, []string{constants.ReactNativeTemplate,
			constants.AndroidTemplate,
			constants.IOSTemplate,
			constants.FlutterTemplate,
		})

	default:
		errorhandler.CheckNilErr(fmt.Errorf("Selected stack is invalid"))
	}
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
