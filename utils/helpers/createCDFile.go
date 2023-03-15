package helpers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func CreateCDFile(stack, dirName, database string) error {
	var cdFileUrl string

	switch stack {
	case constants.NODE_HAPI_TEMPLATE:
		cdFileUrl = fmt.Sprintf("%s%s%s", constants.GitHubBaseURL,
			constants.NodeHapiTemplateRepo,
			constants.CDFilePathURL,
		)
	case constants.NODE_EXPRESS_GRAPHQL_TEMPLATE:
		cdFileUrl = fmt.Sprintf("%s%s%s", constants.GitHubBaseURL,
			constants.NodeExpressGraphqlTemplateRepo,
			constants.CDFilePathURL,
		)
	case constants.GOLANG_ECHO_TEMPLATE:
		if database == constants.POSTGRES {
			cdFileUrl = fmt.Sprintf("%s%s%s", constants.GitHubBaseURL,
				constants.GoEchoTemplatePostgresRepo,
				constants.CDFilePathURL,
			)
		} else if database == constants.MYSQL {
			cdFileUrl = fmt.Sprintf("%s%s%s", constants.GitHubBaseURL,
				constants.GoEchoTemplateMysqlRepo,
				constants.CDFilePathURL,
			)
		}
	case constants.REACT:
		cdFileUrl = fmt.Sprintf("%s%s%s", constants.GitHubBaseURL,
			constants.ReactTemplateRepo,
			constants.CDFilePathURL,
		)
	case constants.NEXT:
		cdFileUrl = fmt.Sprintf("%s%s%s", constants.GitHubBaseURL,
			constants.NextjsTemplateRepo,
			constants.CDFilePathURL,
		)
	default:
		return fmt.Errorf("Selected stack is invalid")
	}

	cdDestination := fileutils.CurrentDirectory() + "/" + dirName + constants.CDFilePathURL
	status, _ := fileutils.IsExists(cdDestination)
	if !status {

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

	} else {
		fmt.Println("The", dirName, stack, "CD you are looking to create already exists")
	}
	return nil
}
