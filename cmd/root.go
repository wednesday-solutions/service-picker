package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
)

// RootCmd is the command variable of root command picky.
var RootCmd = RootCmdFn()
var version = "0.0.8"

// RootCmd represents the base command when called without any subcommands
func RootCmdFn() *cobra.Command {

	var cmd = &cobra.Command{
		Use:     constants.Picky,
		Version: version,
		Short:   "Service Picker",
		Long: fmt.Sprintf(`
Hello%s
Welcome to Service Picker.

It contains a number of @wednesday-solutions's open source projects, connected and working together. Pick whatever you need and build your own ecosystem.

This repo will have support for production applications using the following tech stacks
- Frontend
  - react
  - next
- Backend
  - Node (Hapi - REST API)
  - Node (Express - GraphQL API)
- Databases
  - MySQL
  - PostgreSQL
- Cache
  - Redis
- Infrastructure Provider
  - AWS

Wednesday Solutions`, errorhandler.WaveMessage),
	}
	return cmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	RootCmd.CompletionOptions.DisableDefaultCmd = true

	err := RootCmd.Execute()
	if err != nil {
		if err.Error() == errorhandler.ErrInterrupt.Error() {
			err = errorhandler.ExitMessage
		}
		errorhandler.CheckNilErr(err)
	}
	return nil
}
