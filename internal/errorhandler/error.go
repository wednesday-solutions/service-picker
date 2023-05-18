package errorhandler

import (
	"fmt"
	"os"
)

func CheckNilErr(err error) {
	if err != nil {
		// Check if the user clicks the Control + C button for exiting.
		if err.Error() == ErrInterrupt.Error() {
			err = ExitMessage
		}
		fmt.Print(err)
		os.Exit(1)
	}
}
