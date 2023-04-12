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
	var cdFileUrl string

	switch stack {
	case constants.NodeHapiTemplate:
		cdFileUrl = fmt.Sprintf("%s%s%s", constants.GitHubBaseURL,
			constants.NodeHapiTemplateRepo,
			constants.CDFilePathURL,
		)
	case constants.NodeExpressGraphqlTemplate:
		cdFileUrl = fmt.Sprintf("%s%s%s", constants.GitHubBaseURL,
			constants.NodeExpressGraphqlTemplateRepo,
			constants.CDFilePathURL,
		)
	case constants.GolangEchoTemplate:
		if database == constants.PostgreSQL {
			cdFileUrl = fmt.Sprintf("%s%s%s", constants.GitHubBaseURL,
				constants.GoEchoTemplatePostgresRepo,
				constants.CDFilePathURL,
			)
		} else if database == constants.MySQL {
			cdFileUrl = fmt.Sprintf("%s%s%s", constants.GitHubBaseURL,
				constants.GoEchoTemplateMysqlRepo,
				constants.CDFilePathURL,
			)
		}
	case constants.ReactJS:
		cdFileUrl = fmt.Sprintf("%s%s%s", constants.GitHubBaseURL,
			constants.ReactTemplateRepo,
			constants.CDFilePathURL,
		)
	case constants.NextJS:
		cdFileUrl = fmt.Sprintf("%s%s%s", constants.GitHubBaseURL,
			constants.NextjsTemplateRepo,
			constants.CDFilePathURL,
		)
	case constants.ReactGraphqlTS:
		cdFileUrl = fmt.Sprintf("%s%s%s", constants.GitHubBaseURL,
			constants.ReactGraphqlTemplateRepo,
			constants.CDFilePathURL,
		)
	case constants.ReactNative:
		cdFileUrl = fmt.Sprintf("%s%s%s", constants.GitHubBaseURL,
			constants.ReactNativeTemplateRepo,
			constants.CDFilePathURL, // build.yml
		)
	case constants.Android:
		cdFileUrl = fmt.Sprintf("%s%s%s", constants.GitHubBaseURL,
			constants.AndroidTemplateRepo,
			constants.CDFilePathURL,
		)
	case constants.IOS:
		cdFileUrl = fmt.Sprintf("%s%s%s", constants.GitHubBaseURL,
			constants.IOSTemplateRepo,
			constants.CDFilePathURL, // ci.yml
		)
	case constants.Flutter:
		cdFileUrl = fmt.Sprintf("%s%s%s", constants.GitHubBaseURL,
			constants.FlutterTemplateRepo,
			constants.CDFilePathURL,
		)
	default:
		return fmt.Errorf("Selected stack is invalid")
	}

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