package main

import (
	"github.com/wednesday-solutions/picky/cmd"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		if err.Error() == errorhandler.ErrInterrupt.Error() {
			err = errorhandler.ExitMessage
		}
		errorhandler.CheckNilErr(err)
	}
}
