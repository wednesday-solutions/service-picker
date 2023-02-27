package utils

func checkNilErr(err error) {
	if err != nil {
		// panic(err) will stop exucution and throw error
		panic(err)
	}
}