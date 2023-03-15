package helpers

import (
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
)

func ProgressBar(max int64, description string, done chan bool) {

	bar := progressbar.Default(max, description)
	for i := 0; i < int(max); i++ {
		err := bar.Add(1)
		errorhandler.CheckNilErr(err)

		time.Sleep(40 * time.Millisecond)
	}
	done <- true
}
