package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wednesday-solutions/picky/utils/prompt"
)

var (
	FRONTEND  = "frontend"
	MOBILE    = "mobile"
	BACKEND   = "backend"
	DATABASES = "databases"
	REACT     = "react"
)

var (
	NODE_HAPI       = "Node (Hapi- REST API)"
	NODE_EXPRESS    = "Node (Express- GraphQL API)"
	NODE_EXPRESS_TS = "Node (Express- TypeScript)"
	GOLANG          = "Golang (Echo- GraphQL API)"
)

// ServiceCmd is the command variable of ServiceSelection.
var ServiceSelection = ServiceSelectionFn()

func RunService(*cobra.Command, []string) {
	selectedService := prompt.PromptSelect("Pick a service", []string{FRONTEND, MOBILE, BACKEND, DATABASES})

	switch selectedService {

	case FRONTEND:
		prompt.PromptSelectStack(FRONTEND, []string{REACT})

	case BACKEND:
		prompt.PromptSelectStack(BACKEND, []string{NODE_HAPI, NODE_EXPRESS, NODE_EXPRESS_TS, GOLANG})

	case DATABASES:

	case MOBILE:

	default:
		{

		}
	}
}

// ServiceSelectionFn represents the ServiceSelection command
func ServiceSelectionFn() *cobra.Command {

	var ServiceSelection = &cobra.Command{
		Use:   "service",
		Short: "Pick a Service",
		Long: `Pick a service for your:

		1. Frontend
		2. Mobile
		3. Backend
		4. Databases

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
