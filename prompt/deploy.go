package prompt

import (
	"fmt"

	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
	"github.com/wednesday-solutions/picky/pickyhelpers"
)

func PromptDeploy() {
	var p PromptInput
	p.Label = "Do you want to deploy your project"
	p.GoBack = PromptHome
	response := p.PromptYesOrNoSelect()
	var stackDirectories []string
	if response {
		_ = PromptCloudProvider()
		configStacks := pickyhelpers.SstConfigExistStacks()
		if len(configStacks) == 0 {
			err := utils.PrintWarningMessage(fmt.Sprintf("Please setup infrastructure first%s\n",
				errorhandler.Exclamation,
			))
			errorhandler.CheckNilErr(err)
			PromptSetupInfra()
		}
		message := "Selected stacks are,\n\n"
		for i, stack := range configStacks {
			message = fmt.Sprintf("%s  %d.%s\n", message, i+1, stack)
		}
		fmt.Printf("%s\n", message)

		p.Label = "Do you want to change the selected stacks"
		p.GoBack = PromptDeploy
		response = p.PromptYesOrNoSelect()
		if response {
			stackDirectories = PromptSelectExistingStacks()
			nonExistingStacks := pickyhelpers.IsInfraStacksExist(stackDirectories)
			if len(nonExistingStacks) > 0 {
				message := "Didn't setup Infrastructure for the following stacks,\n\n"
				for i, stack := range nonExistingStacks {
					message = fmt.Sprintf("%s %d. %s\n", message, i+1, stack)
				}
				fmt.Printf("%s\n", message)
				err := utils.PrintWarningMessage(fmt.Sprintf("Please setup infrastructure first%s\n",
					errorhandler.Exclamation,
				))
				errorhandler.CheckNilErr(err)
				PromptSetupInfra()
			}
		}
		environment := PromptEnvironment()
		if response {
			stackInfo := pickyhelpers.GetStackInfo("", "", environment)
			err := pickyhelpers.CreateSstConfigFile(stackInfo, stackDirectories)
			errorhandler.CheckNilErr(err)
		}
		// Let's deploy..
		err := PromptDeployUtils(environment)
		errorhandler.CheckNilErr(err)
	}
	PromptHome()
}

func PromptDeployAfterInfra(services []string, environment string) {
	var p PromptInput
	p.Label = "Do you want to deploy your project"
	p.GoBack = PromptHome
	response := p.PromptYesOrNoSelect()
	if response {
		_ = PromptCloudProvider()
		if len(services) > 0 {
			message := "Selected stacks are,\n\n"
			for i, service := range services {
				message = fmt.Sprintf("%s  %d. %s\n", message, i+1, service)
			}
			fmt.Printf("%s\n", message)
			p.Label = "Are you sure to deploy the above stacks"
			response = p.PromptYesOrNoSelect()
			if response {
				// Let's deploy...
				err := PromptDeployUtils(environment)
				errorhandler.CheckNilErr(err)
			}
		} else {
			err := utils.PrintWarningMessage(fmt.Sprintf("Please setup infrastructure first%s\n",
				errorhandler.Exclamation,
			))
			errorhandler.CheckNilErr(err)
			PromptSetupInfra()
		}
	}
	PromptHome()
}

func PromptDeployUtils(environment string) error {

	// Install dependencies ('yarn install' or 'npm install')
	err := PromptInstallDependencies()
	errorhandler.CheckNilErr(err)

	// Build sst
	err = PromptBuildSST()
	errorhandler.CheckNilErr(err)

	// Let's deploy...
	err = PromptDeploySST(environment)
	errorhandler.CheckNilErr(err)

	err = utils.PrintWarningMessage("Work in progress. Please stay tuned..!")
	errorhandler.CheckNilErr(err)

	return err
}

func PromptInstallDependencies() error {
	var p PromptInput
	p.Label = "Can we install dependencies"
	p.GoBack = PromptHome
	count, pkgManager := 0, ""
	for {
		response := p.PromptYesOrNoSelect()
		count++
		if count == 1 {
			pkgManager = utils.IsYarnOrNpmInstalled()
		}
		if response {
			err := pickyhelpers.InstallDependencies(pkgManager)
			return err
		}
		if count == 2 {
			PromptHome()
		}
		err := utils.PrintWarningMessage("You can't deploy without installing dependencies")
		errorhandler.CheckNilErr(err)
	}
}

func PromptBuildSST() error {
	var p PromptInput
	p.Label = "Can we build"
	p.GoBack = PromptDeploy
	count := 0
	for {
		response := p.PromptYesOrNoSelect()
		count++
		if response {
			err := pickyhelpers.BuildSST()
			return err
		}
		if count == 2 {
			PromptHome()
		}
		err := utils.PrintWarningMessage("You can't deploy without build.")
		errorhandler.CheckNilErr(err)
	}
}

func PromptDeploySST(environment string) error {
	var p PromptInput
	p.Label = "Can we deploy now"
	p.GoBack = PromptDeploy
	count := 0
	for {
		response := p.PromptYesOrNoSelect()
		count++
		if response {
			err := pickyhelpers.DeploySST(environment)
			errorhandler.CheckNilErr(err)
		}
		if count == 2 {
			PromptHome()
		}
	}
}
