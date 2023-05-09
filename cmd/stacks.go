package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wednesday-solutions/picky/flagcmd"
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
)

var StacksCmd = StacksCmdFn()

func RunStacks(cmd *cobra.Command, args []string) error {

	var err error
	var flag flagcmd.StackFlag
	flag.A, err = cmd.Flags().GetBool("allstacks")
	errorhandler.CheckNilErr(err)

	flag.E, err = cmd.Flags().GetBool("existingstacks")
	errorhandler.CheckNilErr(err)

	flag.W, err = cmd.Flags().GetBool("webstacks")
	errorhandler.CheckNilErr(err)

	flag.M, err = cmd.Flags().GetBool("mobilestacks")
	errorhandler.CheckNilErr(err)

	flag.B, err = cmd.Flags().GetBool("backendstacks")
	errorhandler.CheckNilErr(err)

	userOutput := flag.FlagStacks()
	fmt.Println(userOutput)

	return nil
}

func StacksCmdFn() *cobra.Command {
	stacksCmd := &cobra.Command{
		Use:   constants.Stacks,
		Short: "See stacks",
		Long:  "See stacks",
		Args:  cobra.NoArgs,
		RunE:  RunStacks,
	}
	return stacksCmd
}

func init() {
	RootCmd.AddCommand(StacksCmd)
	// declaring flags
	StacksCmd.Flags().BoolP("allstacks", "a", false, "all available stacks")
	StacksCmd.Flags().BoolP("existingstacks", "e", true, "all existing stacks")
	StacksCmd.Flags().BoolP("webstacks", "w", false, "web stacks")
	StacksCmd.Flags().BoolP("mobilestacks", "m", false, "mobile stacks")
	StacksCmd.Flags().BoolP("backendstacks", "b", false, "backend stacks")
}
