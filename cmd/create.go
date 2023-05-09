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

	var (
		flagDockercompose bool
		flagCI            bool
		flagCD            bool
		flagPlatform      string
		flagStacks        []string
		flagEnv           string
		err               error
	)

	flagDockercompose, err = cmd.Flags().GetBool(constants.DockerComposeFlag)
	errorhandler.CheckNilErr(err)

	flagCI, err = cmd.Flags().GetBool(constants.CIFlag)
	errorhandler.CheckNilErr(err)

	flagCD, err = cmd.Flags().GetBool(constants.CDFlag)
	errorhandler.CheckNilErr(err)

	flagPlatform, err = cmd.Flags().GetString(constants.Platform)
	errorhandler.CheckNilErr(err)

	flagStacks, err = cmd.Flags().GetStringSlice(constants.Stacks)
	errorhandler.CheckNilErr(err)

	flagEnv, err = cmd.Flags().GetString(constants.Environment)
	errorhandler.CheckNilErr(err)

	err = utils.CheckStacksExist(flagStacks)
	if err != nil {
		return err
	}

	fmt.Printf("\nDocker Compose: %v\nCI: %v\nCD: %v\nPlatform: %s\nStacks: %v\n",
		flagDockercompose, flagCI, flagCD, flagPlatform, flagStacks,
	)

	if !strings.EqualFold(flagPlatform, constants.GitHub) {
		return fmt.Errorf("Only GitHub is available now.\n")
	}
	if flagDockercompose {
		err = prompt.GenerateDockerCompose()
		errorhandler.CheckNilErr(err)
	}
	if flagCI {
		err = pickyhelpers.CreateCI(flagStacks)
		errorhandler.CheckNilErr(err)
	}
	if flagCD {
		err = prompt.CreateCD(flagStacks, flagEnv)
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
	ci              bool
	cd              bool
	dockerCompose   bool
	platform        string
	flagEnvironment string
)

func init() {
	ServiceCmd.AddCommand(CreateCmd)
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
	CreateCmd.Flags().StringVarP(
		&flagEnvironment, constants.Environment, "e", constants.Development, utils.UseEnvironment(),
	)
}
