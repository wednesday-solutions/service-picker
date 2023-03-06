package errorhandler

func CheckNilErr(err error) {
	if err != nil {
		// panic(err) will stop execution and throw error
		panic(err)
	}
}
