package main

import (
	"log"

	"github.com/wednesday-solutions/picky/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalln("Something error happened: ", err)
	}
}
