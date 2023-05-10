package prompt

import (
	"fmt"

	"github.com/wednesday-solutions/picky/internal/constants"
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
	response, afterInfra := true, false
	if len(stacks) == 0 {
		stacks = utils.GetExistingInfraStacks()
	} else {
		// afterInfra will become true if the DeployStacks function is called from DeployAfterInfra
		afterInfra = true
	}
	if len(stacks) > 0 {
		if !afterInfra {
			stacks = utils.FindStackDirectoriesByConfigStacks(stacks)
		}
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
	var s pickyhelpers.StackDetails
	s.Environment = environment
	s.StackInfo = s.GetStackInfo()
	err := pickyhelpers.CreateSstConfigFile(s.StackInfo, stacks)
	errorhandler.CheckNilErr(err)

	// Deploy infrastructure
	err = PromptInstallDependenciesAndDeploy(stacks, environment)
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

// PromptBuildSST runs 'yarn build'
func PromptBuildSST(pkgManager string) error {
	var p PromptInput
	p.Label = "Do you want to build"
	p.GoBack = PromptDeploy
	response := p.PromptYesOrNoSelect()
	if response {
		err := pickyhelpers.BuildSST(pkgManager)
		return err
	} else {
		PromptDeploy()
	}
	return nil
}

// InstallDependenciesAndDeploy install dependencies of each file, then deploy.
func PromptInstallDependenciesAndDeploy(configStacks []string, environment string) error {
	var p PromptInput
	p.Label = "Can we deploy now"
	p.GoBack = PromptDeploy
	response := p.PromptYesOrNoSelect()
	pkgManager := utils.GetPackageManagerOfUser()
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
		err = utils.PrintInfoMessage("Deploying...")
		errorhandler.CheckNilErr(err)
		err = pickyhelpers.DeploySST(pkgManager, environment)
		errorhandler.CheckNilErr(err)

		err = pickyhelpers.ParseDeployOutputs()
		errorhandler.CheckNilErr(err)

		err = utils.CreateSstOutputsFile()
		errorhandler.CheckNilErr(err)

	} else {
		PromptDeploy()
	}
	return nil
}

func ShowRemoveDeploy() bool {
	path := fmt.Sprintf("%s/%s", utils.CurrentDirectory(), constants.DotSstDirectory)
	status, _ := utils.IsExists(path)
	return status
}

// PromptRemoveDeploy is the prompt for the remove deploy option of Home prompt.
func PromptRemoveDeploy() {
	var p PromptInput
	p.Label = "Do you want to remove the deployed infrastructure"
	p.GoBack = PromptHome
	response := p.PromptYesOrNoSelect()
	if response {
		environment := PromptEnvironment()
		err := RemoveDeploy(environment)
		errorhandler.CheckNilErr(err)
	}
	PromptHome()
}

func RemoveDeploy(environment string) error {
	pkgManager := utils.GetPackageManagerOfUser()
	environment = utils.GetEnvironment(environment)

	err := utils.PrintInfoMessage("Removing deployed infrastructure..")
	errorhandler.CheckNilErr(err)
	arg := fmt.Sprintf("%s:%s", "remove", environment)
	err = utils.RunCommandWithLogs("", pkgManager, "run", arg)
	return err
}
