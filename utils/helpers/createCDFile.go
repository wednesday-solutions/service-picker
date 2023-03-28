package helpers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func CreateCDFile(stack, service, database string) error {
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
	default:
		return fmt.Errorf("Selected stack is invalid")
	}

	cdDestination := fileutils.CurrentDirectory() + "/" + service + constants.CDFilePathURL
	status, _ := fileutils.IsExists(cdDestination)
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
		err = fileutils.CreateFile(cdDestination)
		errorhandler.CheckNilErr(err)

		// Write CDFileData to CD File
		err = fileutils.WriteToFile(cdDestination, string(cdFileData))
		errorhandler.CheckNilErr(err)

		<-done
		fmt.Printf("\n%s%s", "Generating", errorhandler.CompleteMessage)

	} else {
		fmt.Println("The", service, stack, "CD you are looking to create already exists")
	}
	return nil
}
