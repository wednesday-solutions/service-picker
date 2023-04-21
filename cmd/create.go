package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
	"github.com/wednesday-solutions/picky/pickyhelpers"
	"github.com/wednesday-solutions/picky/prompt"
)

func CreateService(cmd *cobra.Command, args []string) error {

	dockerCompose, err := cmd.Flags().GetBool(constants.DockerComposeFlag)
	errorhandler.CheckNilErr(err)

	ci, err := cmd.Flags().GetBool(constants.CIFlag)
	errorhandler.CheckNilErr(err)

	cd, err := cmd.Flags().GetBool(constants.CDFlag)
	errorhandler.CheckNilErr(err)

	platform, err := cmd.Flags().GetString(constants.Platform)
	errorhandler.CheckNilErr(err)

	stacks, err = cmd.Flags().GetStringSlice(constants.Stacks)
	errorhandler.CheckNilErr(err)

	err = utils.CheckStacksExist(stacks)
	if err != nil {
		return err
	}

	fmt.Printf("\nDocker Compose: %v\nCI: %v\nCD: %v\nPlatform: %s\nStacks: %v\n",
		dockerCompose, ci, cd, platform, stacks,
	)

	if !strings.EqualFold(platform, constants.GitHub) {
		return fmt.Errorf("Only GitHub is available now.\n")
	}
	if dockerCompose {
		err = prompt.GenerateDockerCompose()
		errorhandler.CheckNilErr(err)
	}
	if ci {
		err = pickyhelpers.CreateCI(stacks)
		errorhandler.CheckNilErr(err)
	}
	if cd {
		err = prompt.CreateCD(stacks)
		errorhandler.CheckNilErr(err)
	}
	return err
}

func CreateCmdFn() *cobra.Command {
	var CreateCommand = &cobra.Command{
		Use:  constants.Create,
		RunE: CreateService,
	}
	return CreateCommand
}

var CreateCmd = CreateCmdFn()
var (
	ci            bool
	cd            bool
	dockerCompose bool
	platform      string
)

func init() {
	ServiceSelection.AddCommand(CreateCmd)
	CreateCmd.Flags().BoolVarP(&ci, constants.CIFlag, "i", false, utils.UseCI())
	CreateCmd.Flags().BoolVarP(&cd, constants.CDFlag, "d", false, utils.UseCD())
	CreateCmd.Flags().BoolVarP(
		&dockerCompose, constants.DockerComposeFlag, "c", false, utils.UseDockerCompose(),
	)
	CreateCmd.Flags().StringVarP(
		&platform, constants.Platform, "p", constants.Github, utils.UsePlatform(),
	)
	CreateCmd.Flags().StringSliceVarP(
		&stacks, constants.Stacks, "t", utils.ExistingStacks(), utils.UseInfraStacks(),
	)
}
