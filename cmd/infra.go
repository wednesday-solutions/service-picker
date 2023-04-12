package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
)

func InfraSetup(cmd *cobra.Command, args []string) error {
	var (
		err           error
		stacks        []string
		stackExist    bool
		cloudProvider string
		environment   string
	)
	stacks, err = cmd.Flags().GetStringSlice("stacks")
	errorhandler.CheckNilErr(err)

	if len(stacks) == 0 {
		return fmt.Errorf("Provide existing stacks\n")
	}
	_, _, directories := utils.ExistingStacksDatabasesAndDirectories()
	for _, stack := range stacks {
		for _, dir := range directories {
			if stack == dir {
				stackExist = true
			}
		}
		if !stackExist {
			return fmt.Errorf("Entered stack '%s' not exists\n", stack)
		}
	}
	cloudProvider, err = cmd.Flags().GetString("cloudprovider")
	errorhandler.CheckNilErr(err)
	environment, err = cmd.Flags().GetString("environment")
	errorhandler.CheckNilErr(err)

	cloudProvider = utils.GetCloudProvider(cloudProvider)
	environment = utils.GetEnvironmentValue(environment)

	fmt.Printf("\nStacks: %v\nCloud Provider: %s\nEnvironment: %s\n\n",
		stacks, cloudProvider, environment,
	)
	return err
}

func InfraCmdFn() *cobra.Command {
	var InfraCmd = &cobra.Command{
		Use:  constants.Infra,
		RunE: InfraSetup,
	}
	return InfraCmd
}

var InfraCmd = InfraCmdFn()
var (
	stacks        []string
	cloudProvider string
	environment   string
)

func init() {
	ServiceSelection.AddCommand(InfraCmd)
	InfraCmd.Flags().StringSliceVarP(
		&stacks, "stacks", "t", utils.ExistingStacks(), utils.UsageInfraStacks(),
	)
	InfraCmd.Flags().StringVarP(
		&cloudProvider, "cloudprovider", "p", constants.AWS, utils.UsageCloudProvider(),
	)
	InfraCmd.Flags().StringVarP(
		&environment, "environment", "e", constants.Development, utils.UsageEnvironment(),
	)
}
