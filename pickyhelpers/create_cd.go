package pickyhelpers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
)

func CreateCDFile(service, stack, database, dirName string) error {

	// Need to write the CD file.
	var cdFileUrl string
	cdDestination := fmt.Sprintf("%s/%s/cd-%s.yml",
		utils.CurrentDirectory(),
		constants.GithubWorkflowsDir,
		dirName,
	)
	status, _ := utils.IsExists(cdDestination)
	if !status {

		done := make(chan bool)
		go ProgressBar(20, "Generating", done)

		// Access CD File which is present in the Github.
		resp, err := http.Get(cdFileUrl)
		errorhandler.CheckNilErr(err)
		defer resp.Body.Close()

		// Read the body of the response.
		cdFileData, err := io.ReadAll(resp.Body)
		errorhandler.CheckNilErr(err)

		// Create CD File
		cdFileExist, _ := utils.IsExists(cdDestination)
		if !cdFileExist {
			utils.CreateGithubWorkflowDir()
			err = utils.CreateFile(cdDestination)
			errorhandler.CheckNilErr(err)
		}

		// Write CDFileData to CD File
		err = utils.WriteToFile(cdDestination, string(cdFileData))
		errorhandler.CheckNilErr(err)

		<-done
		fmt.Printf("\n%s %s", "Generating", errorhandler.CompleteMessage)

	} else {
		fmt.Println("The", service, stack, "CD you are looking to create already exists")
		return errorhandler.ErrExist
	}
	return nil
}
