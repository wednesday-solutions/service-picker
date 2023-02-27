package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wednesday-solutions/service-picker/utils"
)

var (
		FRONTEND = "frontend"
		MOBILE = "mobile"
		BACKEND = "backend"
		DATABASES = "databases"	
		REACT = "react"
	)

// ServiceCmd is the command variable of ServiceSelection.
var ServiceSelection = ServiceSelectionFn()

func RunService(*cobra.Command, []string) {
	selectedService:= utils.PromptSelect("Pick a service", []string{FRONTEND, MOBILE, BACKEND, DATABASES})
	
	switch selectedService {

	case FRONTEND: 
	
	utils.PromptSelectStack(FRONTEND, []string{REACT})

	case BACKEND: 
		

	case DATABASES: 
		
	
	case MOBILE: 
	

	default: {
		
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

	// negt gqlgen only give suggestions
	ServiceSelection.Flags().BoolP("help", "h", false, "Help for service selection")
}
