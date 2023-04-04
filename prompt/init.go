package prompt

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/wednesday-solutions/picky/pickyhelpers"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func PromptSelectInit(service, stack, database, dirName string) {
	response := true
	for response {
		label := fmt.Sprintf("Do you want to initialize '%s' as '%s'", stack, dirName)
		response = PromptYesOrNoSelect(label)
		if response {
			Init(service, stack, database, dirName)
		} else {
			response = PromptConfirm()
			if response {
				PromptHome()
			} else {
				response = true
			}
		}
	}
}

func Init(service, stack, database, dirName string) {

	var label string
	currentDir := fileutils.CurrentDirectory()
	if stack == constants.GolangEchoTemplate {
		stack = fmt.Sprintf("%s-%s", strings.Split(stack, " ")[0], database)
	}
	destination := filepath.Join(currentDir, dirName)
	status, _ := fileutils.IsExists(destination)
	var response bool
	for status {
		label = fmt.Sprintf("The '%s' already exists, do you want to update it", dirName)
		response = PromptYesOrNoSelect(label)
		if response {
			// Delete all contents of existing directory.
			err := fileutils.RemoveAllContents(destination)
			errorhandler.CheckNilErr(err)
		} else {
			PromptHome()
		}
	}
	if !status {
		// Create directory with directory name we got.
		err := fileutils.MakeDirectory(currentDir, dirName)
		errorhandler.CheckNilErr(err)
		response = true
	}
	if response {
		done := make(chan bool)
		go pickyhelpers.ProgressBar(100, "Downloading", done)

		// Clone the selected repo into service directory.
		err := pickyhelpers.CloneRepo(stack, dirName, currentDir)
		errorhandler.CheckNilErr(err)

		// stackInfo gives the information about the stacks which is present in the root.
		stackInfo := pickyhelpers.GetStackInfo(stack, database, constants.Environment)

		// Database conversion
		if service == constants.Backend {
			err = pickyhelpers.ConvertTemplateDatabase(stack, database, dirName, stackInfo)
			errorhandler.CheckNilErr(err)
		}
		// create and update docker files
		err = pickyhelpers.CreateDockerFiles(dirName, stackInfo)
		errorhandler.CheckNilErr(err)

		<-done
		fmt.Printf("\nDownloading %s", errorhandler.CompleteMessage)
	}
	label = "Do you want to initialize another service"
	response = PromptYesOrNoSelect(label)
	if response {
		PromptSelectService()
	} else {
		PromptHome()
	}
}
