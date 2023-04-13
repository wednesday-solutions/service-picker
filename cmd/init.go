package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wednesday-solutions/picky/flagcommand"
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
)

func InitService(cmd *cobra.Command, args []string) error {

	var (
		i            flagcommand.InitInfo
		err          error
		errorMessage string
	)

	i.Service, err = cmd.Flags().GetString("service")
	errorhandler.CheckNilErr(err)
	i.Stack, err = cmd.Flags().GetString("stack")
	errorhandler.CheckNilErr(err)
	i.Database, err = cmd.Flags().GetString("database")
	errorhandler.CheckNilErr(err)
	i.Directory, err = cmd.Flags().GetString("directory")
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
		RunE: InitService,
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
	InitCmd.Flags().StringVarP(&service, "service", "s", "", utils.UseService())
	InitCmd.Flags().StringVarP(&stack, "stack", "t", "", utils.UseStack())
	InitCmd.Flags().StringVarP(&database, "database", "d", "", utils.UseDatabase())
	InitCmd.Flags().StringVarP(&directory, "directory", "f", "", utils.UseDirectory())
}