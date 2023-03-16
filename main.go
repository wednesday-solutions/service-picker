package main

import (
	"log"

	"github.com/wednesday-solutions/picky/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("Something error happened: %v", err)
	}
}
