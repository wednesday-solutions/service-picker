package errorhandler

import (
	"fmt"
	"log"
)

func CheckNilErr(err error) {
	if err != nil {
		// Check if the user clicks the Control + C button for exiting.
		if err.Error() == "^C" {
			err = fmt.Errorf("Program exited")
		}
		log.Fatalln(err) // it will throw error and stop execution.
	}
}
