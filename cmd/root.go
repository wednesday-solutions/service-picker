package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
)

// RootCmd is the command variable of root command picky.
var RootCmd = RootCmdFn()
var version = "0.0.1"

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
- mobile
  - android app
  - iOS app
  - react-native app
  - flutter app
- frontend
  - react
  - next
- backend
  - Node (Hapi - REST API)
  - Node (Express - GraphQL API)
  - Node (Express - TypeScript)
  - Golang (Echo - GraphQL API)
- Databases
  - MySQL
  - PostgreSQL
  - MongoDB
  - DynamoDB
  - Neo4j
- Infrastructure
  - Redis
  - Kafka

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
		if err == errorhandler.ErrInterrupt {
			err = errorhandler.ExitMessage
		}
		errorhandler.CheckNilErr(err)
	}
	return nil
}
