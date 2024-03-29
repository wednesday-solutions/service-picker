package pickyhelpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
	"github.com/wednesday-solutions/picky/pickyhelpers/sources"
)

func (s StackDetails) CreateCDFile() error {

	utils.CreateGithubWorkflowDir()
	cdDestination := fmt.Sprintf("%s/%s/cd-%s.yml",
		utils.CurrentDirectory(),
		constants.GithubWorkflowsDir,
		s.DirName,
	)
	status, _ := utils.IsExists(cdDestination)
	if !status {

		done := make(chan bool)
		go ProgressBar(20, "Generating", done)

		// Create CD File
		err := utils.CreateFile(cdDestination)
		errorhandler.CheckNilErr(err)

		var cdSource string
		if s.Service == constants.Backend {
			cdSource = sources.CDBackendSource(s.Stack, s.DirName, s.Environment)
		} else if s.Service == constants.Web {
			cdSource = sources.CDWebSource(s.Stack, s.DirName)
		}

		// Write CDFileData to CD File
		err = utils.WriteToFile(cdDestination, cdSource)
		errorhandler.CheckNilErr(err)

		<-done
		fmt.Printf("\n%s %s", "Generating", errorhandler.CompleteMessage)

	} else {
		fmt.Println("The", s.Service, s.Stack, "CD you are looking to create already exists")
		return errorhandler.ErrExist
	}
	return nil
}

func CreateTaskDefinition(stackDir, environment string) error {
	environment = utils.GetShortEnvName(environment)

	file := fmt.Sprintf("%s/%s/%s-%s.json",
		utils.CurrentDirectory(),
		stackDir,
		"task-definition",
		environment,
	)
	source := sources.TaskDefinitionSource(environment, stackDir)
	var err error
	if source != "" {
		err = utils.WriteToFile(file, source)
	}
	return err
}
