package prompt

import (
	"fmt"
	"strings"

	"github.com/wednesday-solutions/picky/pickyhelpers"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func PromptSelectInit(service, stack, database string) {

	currentDir := fileutils.CurrentDirectory()
	if stack == constants.GolangEchoTemplate {
		stack = fmt.Sprintf("%s%s", strings.Split(stack, " ")[0], database)
	}

	destination := currentDir + "/" + service
	status, _ := fileutils.IsExists(destination)
	var response bool
	if status {
		label := fmt.Sprintf("The %s service already exists, do you want to rewrite it?", service)
		response = PromptYesOrNoSelect(label)
		if response {
			// Delete existing service.
			err := fileutils.RemoveContent(destination)
			errorhandler.CheckNilErr(err)
		}
	}
	if !status {
		// Create directory in the name of selected service.
		err := fileutils.MakeDirectory(currentDir, service)
		errorhandler.CheckNilErr(err)
	}
	if !status || response {
		done := make(chan bool)
		go pickyhelpers.ProgressBar(100, "Downloading", done)

		// Clone the selected repo into service directory.
		err := pickyhelpers.CloneRepo(stack, service, currentDir)
		errorhandler.CheckNilErr(err)

		// stackInfo gives the information about the stacks which is present in the root.
		stackInfo := pickyhelpers.GetStackInfo(stack, database)

		// Database conversion
		if service == constants.Backend {
			err = ConvertTemplateDatabase(stack, database, stackInfo)
			errorhandler.CheckNilErr(err)
		}
		// create and update docker files
		err = pickyhelpers.CreateDockerFiles(stackInfo)
		errorhandler.CheckNilErr(err)

		var dockerComposeFileExist bool
		if stackInfo[constants.BackendStatus].(bool) &&
			(stackInfo[constants.WebStatus].(bool) || stackInfo[constants.MobileStatus].(bool)) {
			// create docker-compose file
			err = pickyhelpers.CreateDockerComposeFile(stackInfo, dockerComposeFileExist)
			if err == errorhandler.ErrExist {
				dockerComposeFileExist = true
			} else {
				errorhandler.CheckNilErr(err)
			}
		}
		<-done
		fmt.Printf("\nDownloading %s", errorhandler.CompleteMessage)

		if dockerComposeFileExist {
			label := fmt.Sprintf("%s already exist, do you want to rewrite it?", constants.DockerComposeFile)
			response = PromptYesOrNoSelect(label)
			if response {
				err = pickyhelpers.CreateDockerComposeFile(stackInfo, dockerComposeFileExist)
				errorhandler.CheckNilErr(err)
				fmt.Print(errorhandler.DoneMessage)
			} else {
				fmt.Print(errorhandler.ExitMessage)
			}
		}
	} else {
		fmt.Print(errorhandler.ExitMessage)
	}
}
