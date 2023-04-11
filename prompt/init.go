package prompt

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
	"github.com/wednesday-solutions/picky/pickyhelpers"
)

func PromptSelectInit(service, stack, database, dirName string) {
	var p PromptInput
	p.GoBack = PromptSelectService
	response := true
	for response {
		p.Label = fmt.Sprintf("Do you want to initialize '%s' as '%s'", stack, dirName)
		response = p.PromptYesOrNoSelect()
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
	var p PromptInput
	p.GoBack = PromptSelectService
	currentDir := utils.CurrentDirectory()
	if stack == constants.GolangEchoTemplate {
		stack = fmt.Sprintf("%s-%s", strings.Split(stack, " ")[0], database)
	}
	destination := filepath.Join(currentDir, dirName)
	status, _ := utils.IsExists(destination)
	var response bool
	for status {
		p.Label = fmt.Sprintf("The '%s' already exists, do you want to update it", dirName)
		response = p.PromptYesOrNoSelect()
		if response {
			// Delete all contents of existing directory.
			err := utils.RemoveAllContents(destination)
			errorhandler.CheckNilErr(err)
		} else {
			PromptHome()
		}
	}
	if !status {
		// Create directory with directory name we got.
		err := utils.MakeDirectory(currentDir, dirName)
		errorhandler.CheckNilErr(err)
		response = true
	}
	if response {
		done := make(chan bool)
		go pickyhelpers.ProgressBar(100, "Downloading", done)

		// Clone the selected repo into service directory.
		err := pickyhelpers.CloneRepo(stack, dirName, currentDir)
		errorhandler.CheckNilErr(err)

		// Delete .git folder inside the cloned repo.
		err = pickyhelpers.DeleteDotGitFolder(dirName)
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
	p.Label = "Do you want to initialize another service"
	response = p.PromptYesOrNoSelect()
	if response {
		PromptSelectService()
	} else {
		PromptHome()
	}
}
