package main

import (
	"fmt"
	"os"

	"github.com/wednesday-solutions/service-picker/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Println("Something error happened: ", err)
		os.Exit(1)
	}
}