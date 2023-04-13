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

type InitInfo struct {
	Service  string
	Stack    string
	Database string
	DirName  string
}

func (i *InitInfo) PromptSelectInit() {
	var p PromptInput
	p.GoBack = PromptSelectService
	response := true
	for response {
		p.Label = fmt.Sprintf("Do you want to initialize '%s' as '%s'", i.Stack, i.DirName)
		response = p.PromptYesOrNoSelect()
		if response {
			i.Init()
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

func (i *InitInfo) Init() {
	var p PromptInput
	p.GoBack = PromptSelectService
	currentDir := utils.CurrentDirectory()
	if i.Stack == constants.GolangEchoTemplate {
		i.Stack = fmt.Sprintf("%s-%s", strings.Split(i.Stack, " ")[0], i.Database)
	}
	destination := filepath.Join(currentDir, i.DirName)
	status, _ := utils.IsExists(destination)
	var response bool
	for status {
		p.Label = fmt.Sprintf("The '%s' already exists, do you want to update it", i.DirName)
		response = p.PromptYesOrNoSelect()
		if response {
			// Delete all contents of existing directory.
			err := utils.RemoveAllContents(destination)
			errorhandler.CheckNilErr(err)
			break
		} else {
			PromptHome()
		}
	}
	if !status {
		// Create directory with directory name we got.
		err := utils.MakeDirectory(currentDir, i.DirName)
		errorhandler.CheckNilErr(err)
		response = true
	}
	if response {
		err := i.StackInitialize()
		errorhandler.CheckNilErr(err)
	}
	p.Label = "Do you want to initialize another service"
	response = p.PromptYesOrNoSelect()
	if response {
		PromptSelectService()
	} else {
		PromptHome()
	}
}

func (i InitInfo) StackInitialize() error {

	done := make(chan bool)
	go pickyhelpers.ProgressBar(100, "Downloading", done)

	// Clone the selected repo into service directory.
	var s pickyhelpers.StackDetails
	s.Stack = i.Stack
	s.DirName = i.DirName
	s.CurrentDir = utils.CurrentDirectory()
	s.Database = i.Database
	err := s.CloneRepo()
	errorhandler.CheckNilErr(err)

	// Delete .git folder inside the cloned repo.
	err = s.DeleteDotGitFolder()
	errorhandler.CheckNilErr(err)

	// stackInfo gives the information about the stacks which is present in the root.
	s.Environment = constants.Environment
	s.StackInfo = s.GetStackInfo()
	// Database conversion
	if i.Service == constants.Backend {
		err = s.ConvertTemplateDatabase()
		errorhandler.CheckNilErr(err)
	}
	// create and update docker files
	err = s.CreateDockerFiles()
	errorhandler.CheckNilErr(err)

	<-done
	fmt.Printf("\nDownloading %s", errorhandler.CompleteMessage)

	return err
}
