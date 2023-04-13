package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
)

func CreateService(cmd *cobra.Command, args []string) error {

	dockerCompose, err := cmd.Flags().GetString("dockercompose")
	errorhandler.CheckNilErr(err)

	ci, err := cmd.Flags().GetString("ci")
	errorhandler.CheckNilErr(err)

	cd, err := cmd.Flags().GetString("cd")
	errorhandler.CheckNilErr(err)

	fmt.Printf("\nDocker Compose: %s\nCI: %s\nCD: %s\n", dockerCompose, ci, cd)

	return nil
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
	dockerCompose string
	ci            string
	cd            string
)

func init() {
	ServiceSelection.AddCommand(CreateCmd)
	CreateCmd.Flags().StringVarP(&dockerCompose, "dockercompose", "c", "", "Write docker compose")
	CreateCmd.Flags().StringVarP(&ci, "ci", "i", "", "Write ci")
	CreateCmd.Flags().StringVarP(&cd, "cd", "d", "", "Write cd")
}
