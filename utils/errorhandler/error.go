package errorhandler

import (
	"log"
)

func CheckNilErr(err error) {
	if err != nil {
		// Check if the user clicks the Control + C button for exiting.
		if err.Error() == ErrInterrupt.Error() {
			err = ExitMessage
		}
		log.Fatal(err) // it will throw error and stop execution.
	}
}
