package errorhandler

import "log"

func CheckNilErr(err error) {
	if err != nil {
		// log.Fatal(err) will stop execution and throw error
		log.Fatal(err)
	}
}
