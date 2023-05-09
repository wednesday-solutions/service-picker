package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wednesday-solutions/picky/flagcmd"
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
)

func InitStack(cmd *cobra.Command, args []string) error {

	var (
		i            flagcmd.InitInfo
		err          error
		errorMessage string
	)

	i.Service, err = cmd.Flags().GetString(constants.Service)
	errorhandler.CheckNilErr(err)
	i.Stack, err = cmd.Flags().GetString(constants.Stack)
	errorhandler.CheckNilErr(err)
	i.Database, err = cmd.Flags().GetString(constants.Database)
	errorhandler.CheckNilErr(err)
	i.Directory, err = cmd.Flags().GetString(constants.Directory)
	errorhandler.CheckNilErr(err)

	allFlagsExist := true
	if i.Service == "" {
		allFlagsExist = false
		errorMessage = fmt.Sprintf("%s%s\n", errorMessage,
			"add service with the flag of '--service'",
		)
	}
	if i.Stack == "" {
		allFlagsExist = false
		errorMessage = fmt.Sprintf("%s%s\n", errorMessage,
			"add stack with the flag of '--stack'",
		)
	}
	if i.Database == "" && service == constants.Backend {
		allFlagsExist = false
		errorMessage = fmt.Sprintf("%s%s\n", errorMessage,
			"add database with the flag of '--database'",
		)
	}
	if i.Directory == "" {
		allFlagsExist = false
		errorMessage = fmt.Sprintf("%s%s\n", errorMessage,
			"add directory with the flag of '--directory'",
		)
	}
	if !allFlagsExist {
		return fmt.Errorf("%s", errorMessage)
	}
	err = i.FlagInit()
	return err
}

func InitCmdFn() *cobra.Command {
	var InitCommand = &cobra.Command{
		Use:  constants.Init,
		RunE: InitStack,
	}
	return InitCommand
}

var InitCmd = InitCmdFn()
var (
	service   string
	stack     string
	database  string
	directory string
)

func init() {
	ServiceSelection.AddCommand(InitCmd)
	InitCmd.Flags().StringVarP(&service, constants.Service, "s", "", utils.UseService())
	InitCmd.Flags().StringVarP(&stack, constants.Stack, "t", "", utils.UseStack())
	InitCmd.Flags().StringVarP(&database, constants.Database, "d", "", utils.UseDatabase())
	InitCmd.Flags().StringVarP(&directory, constants.Directory, "f", "", utils.UseDirectory())
}
