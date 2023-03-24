package helpers

import (
	"fmt"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
)

func ProgressBar(max int, description string, done chan bool) {

	bar := progressbar.NewOptions(max,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(75),
		progressbar.OptionSetDescription(fmt.Sprintf("[cyan][1/1][reset] %s...", description)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]-[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	for i := 0; i < max; i++ {
		err := bar.Add(1)
		errorhandler.CheckNilErr(err)

		time.Sleep(200 * time.Millisecond)
	}
	fmt.Printf("\nCompleted.\n")
	done <- true
}
