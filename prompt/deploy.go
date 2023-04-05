package prompt

import (
	"fmt"

	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
	"github.com/wednesday-solutions/picky/pickyhelpers"
)

// PromptDeploy is the prompt for the deploy option of Home prompt.
func PromptDeploy() {
	var p PromptInput
	p.Label = "Do you want to deploy your project"
	p.GoBack = PromptHome
	response := p.PromptYesOrNoSelect()
	// configStacks means the existing stacks in the sst.config.js
	var configStacks []string
	if response {
		_ = PromptCloudProvider()
		configStacks = pickyhelpers.SstConfigExistStacks()
		if len(configStacks) == 0 {
			err := utils.PrintWarningMessage(fmt.Sprintf("Please setup infrastructure first%s\n",
				errorhandler.Exclamation,
			))
			errorhandler.CheckNilErr(err)
			PromptSetupInfra()
		}
		// Prints the existing stacks in the sst.config.js
		message := "Selected stacks are,\n\n"
		for i, stack := range configStacks {
			message = fmt.Sprintf("%s  %d.%s\n", message, i+1, stack)
		}
		fmt.Printf("%s\n", message)

		p.Label = "Do you want to change the selected stacks"
		p.GoBack = PromptDeploy
		response = p.PromptYesOrNoSelect()
		if response {
			configStacks = PromptSelectExistingStacks()
			// handling errors by checking the selected stacks are present in the stacks directory.
			nonExistingStacks := pickyhelpers.IsInfraStacksExist(configStacks)
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
			err := pickyhelpers.CreateSstConfigFile(stackInfo, configStacks)
			errorhandler.CheckNilErr(err)
		}
		// Let's deploy..
		err := PromptDeployUtils(configStacks, environment)
		errorhandler.CheckNilErr(err)
	}
	PromptHome()
}

// PromptDeployAfterInfra will come up after setting up the infrastructure.
func PromptDeployAfterInfra(configStacks []string, environment string) {
	var p PromptInput
	p.Label = "Do you want to deploy your project"
	p.GoBack = PromptHome
	response := p.PromptYesOrNoSelect()
	if response {
		_ = PromptCloudProvider()
		if len(configStacks) > 0 {
			message := "Selected stacks are,\n\n"
			for i, configStack := range configStacks {
				message = fmt.Sprintf("%s  %d. %s\n", message, i+1, configStack)
			}
			fmt.Printf("%s\n", message)
			p.Label = "Are you sure to deploy the above stacks"
			response = p.PromptYesOrNoSelect()
			if response {
				// Let's deploy...
				err := PromptDeployUtils(configStacks, environment)
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

// PromptDeployUtils contains all three steps of deployments such as install dependencies,
// build, and deployment.
func PromptDeployUtils(configStacks []string, environment string) error {

	// Install dependencies ('yarn install' or 'npm install')
	err := PromptInstallDependencies(configStacks)
	errorhandler.CheckNilErr(err)

	// Build sst
	err = PromptBuildSST()
	errorhandler.CheckNilErr(err)

	// Let's deploy...
	err = PromptDeploySST(environment)
	return err
}

// PromptInstallDependencies will install dependencies
func PromptInstallDependencies(configStacks []string) error {
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
			// install sst dependencies(root directory)
			err := utils.PrintInfoMessage("Installing sst dependencies")
			errorhandler.CheckNilErr(err)
			err = pickyhelpers.InstallDependencies(pkgManager)
			errorhandler.CheckNilErr(err)

			// install selected stacks dependencies(respected stack directory)
			for _, configStackDir := range configStacks {
				err := utils.PrintInfoMessage(fmt.Sprintf("Installing %s dependencies", configStackDir))
				errorhandler.CheckNilErr(err)
				err = pickyhelpers.InstallDependencies(
					pkgManager,
					utils.CurrentDirectory(),
					configStackDir,
				)
				errorhandler.CheckNilErr(err)
			}
			return err
		}
		if count == 2 {
			PromptHome()
		}
		err := utils.PrintWarningMessage("You can't deploy without installing dependencies")
		errorhandler.CheckNilErr(err)
	}
}

// PromptBuildSST runs 'yarn build'
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

// PromptDeploySST runs 'yarn deploy:environment'
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
