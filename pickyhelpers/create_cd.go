package pickyhelpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
)

func CreateCDFile(service, stack, database, dirName string) error {

	err := utils.PrintWarningMessage(fmt.Sprintf(
		"CD of %s is in work in progress..!", stack),
	)
	errorhandler.CheckNilErr(err)
	if err == nil {
		return nil
	}

	utils.CreateGithubWorkflowDir()
	cdDestination := fmt.Sprintf("%s/%s/cd-%s.yml",
		utils.CurrentDirectory(),
		constants.GithubWorkflowsDir,
		dirName,
	)
	status, _ := utils.IsExists(cdDestination)
	if !status {

		done := make(chan bool)
		go ProgressBar(20, "Generating", done)

		// Create CD File
		err := utils.CreateFile(cdDestination)
		errorhandler.CheckNilErr(err)

		// Need to write the CD source to cdSource.
		var cdSource string
		// Write CDFileData to CD File
		err = utils.WriteToFile(cdDestination, cdSource)
		errorhandler.CheckNilErr(err)

		<-done
		fmt.Printf("\n%s %s", "Generating", errorhandler.CompleteMessage)

	} else {
		fmt.Println("The", service, stack, "CD you are looking to create already exists")
		return errorhandler.ErrExist
	}
	return nil
}
