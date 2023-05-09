package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
	"github.com/wednesday-solutions/picky/prompt"
)

var InfraCmd = InfraCmdFn()
var (
	stacks        []string
	cloudProvider string
	environment   string
)

func InfraSetup(cmd *cobra.Command, args []string) error {
	var (
		err           error
		stacks        []string
		cloudProvider string
		environment   string
	)
	stacks, err = cmd.Flags().GetStringSlice(constants.Stacks)
	errorhandler.CheckNilErr(err)

	err = utils.CheckStacksExist(stacks)
	if err != nil {
		return err
	}
	cloudProvider, err = cmd.Flags().GetString(constants.CloudProvider)
	errorhandler.CheckNilErr(err)
	environment, err = cmd.Flags().GetString(constants.Environment)
	errorhandler.CheckNilErr(err)

	cloudProvider = utils.GetCloudProvider(cloudProvider)
	environment = utils.GetEnvironmentValue(environment)

	userOutput := "\nSelected stacks:"
	for idx, stack := range stacks {
		userOutput = fmt.Sprintf("%s\n  %d. %s", userOutput, idx+1, stack)
	}
	fmt.Printf("%s\nCloud Provider: %s\nEnvironment: %s\n\n",
		userOutput, cloudProvider, environment,
	)
	err = prompt.CreateInfra(stacks, cloudProvider, environment)
	errorhandler.CheckNilErr(err)
	return err
}

func InfraCmdFn() *cobra.Command {
	var InfraCmd = &cobra.Command{
		Use:  constants.Infra,
		RunE: InfraSetup,
	}
	return InfraCmd
}

func init() {
	ServiceCmd.AddCommand(InfraCmd)
	InfraCmd.Flags().StringSliceVarP(
		&stacks, constants.Stacks, "t", utils.GetExistingStacks(), utils.UseInfraStacks(),
	)
	InfraCmd.Flags().StringVarP(
		&cloudProvider, constants.CloudProvider, "p", constants.AWS, utils.UseCloudProvider(),
	)
	InfraCmd.Flags().StringVarP(
		&environment, constants.Environment, "e", constants.Development, utils.UseEnvironment(),
	)
}
