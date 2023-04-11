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
	if response {
		err := DeployStacks([]string{}, "")
		errorhandler.CheckNilErr(err)
	}
	PromptHome()
}

// DeployStacks will deploy the infrastructure.
func DeployStacks(stacks []string, environment string) error {
	var p PromptInput
	if environment == "" {
		environment = PromptEnvironment()
	}
	response := true
	if len(stacks) == 0 {
		stacks = utils.GetInfraStacksExist()
	}
	if len(stacks) > 0 {
		stacks = utils.FindStackDirectoriesByConfigStacks(stacks)
		// Prints the existing infra stacks.
		message := "Existing stacks are,\n\n"
		for i, stack := range stacks {
			message = fmt.Sprintf("%s  %d. %s\n", message, i+1, stack)
		}
		fmt.Printf("%s\n", message)
		p.Label = "Do you want to change the existing stacks"
		p.GoBack = PromptDeploy
		response = p.PromptYesOrNoSelect()
	}
	if response {
		stacks = PromptSelectExistingStacks()
		nonExistingStacks := pickyhelpers.GetNonExistingInfraStacks(stacks)
		if len(nonExistingStacks) > 0 {
			message := "Didn't setup Infrastructure for the following stacks,\n\n"
			for i, stack := range nonExistingStacks {
				message = fmt.Sprintf("%s %d. %s\n", message, i+1, stack)
			}
			fmt.Printf("%s\n", message)
			// create infra stacks for non existing stacks.
			err := PromptCreateInfraStacksWhenDeploy(nonExistingStacks, environment)
			errorhandler.CheckNilErr(err)
		}
	}
	stackInfo := pickyhelpers.GetStackInfo("", "", environment)
	err := pickyhelpers.CreateSstConfigFile(stackInfo, stacks)
	errorhandler.CheckNilErr(err)
	err = PromptDeployUtils(stacks, environment)
	return err
}

// PromptDeployAfterInfra will come up after setting up the infrastructure.
func PromptDeployAfterInfra(configStacks []string, environment string) error {
	var p PromptInput
	p.Label = "Do you want to deploy your project"
	p.GoBack = PromptHome
	response := p.PromptYesOrNoSelect()
	if response {
		err := DeployStacks(configStacks, environment)
		return err
	}
	return nil
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

	// Deploy infrastructure
	err = PromptDeploySST(environment)
	return err
}

// PromptInstallDependencies will install dependencies
func PromptInstallDependencies(configStacks []string) error {
	var p PromptInput
	p.Label = "Can we install dependencies"
	p.GoBack = PromptHome
	pkgManager := ""
	response := p.PromptYesOrNoSelect()
	pkgManager = utils.GetPackageManagerOfUser()
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
	} else {
		PromptDeploy()
	}
	return nil
}

// PromptBuildSST runs 'yarn build'
func PromptBuildSST() error {
	var p PromptInput
	p.Label = "Can we build"
	p.GoBack = PromptDeploy
	response := p.PromptYesOrNoSelect()
	if response {
		err := pickyhelpers.BuildSST()
		return err
	} else {
		PromptDeploy()
	}
	return nil
}

// PromptDeploySST runs 'yarn deploy:environment'
func PromptDeploySST(environment string) error {
	var p PromptInput
	p.Label = "Can we deploy now"
	p.GoBack = PromptDeploy
	response := p.PromptYesOrNoSelect()
	if response {
		err := pickyhelpers.DeploySST(environment)
		errorhandler.CheckNilErr(err)
	} else {
		PromptDeploy()
	}
	return nil
}
