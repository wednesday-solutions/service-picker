package helpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func CreateDockerEnvFile(webStatus, mobileStatus bool) error {

	dockerEnvFile := ".env.docker"
	var dockerEnvFound bool
	var err error
	var envPath, envSource string
	if webStatus {
		envPath = fmt.Sprintf("%s/%s/%s", fileutils.CurrentDirectory(),
			constants.Web,
			dockerEnvFile,
		)
		dockerEnvFound, err = fileutils.IsExists(envPath)
		errorhandler.CheckNilErr(err)
		if !dockerEnvFound {
			envSource = `REACT_APP_GOOGLE_LOGIN_CLIENT_ID=775156772157-nlfs5bn4skllfemcvhr6kbphp61achn7.apps.googleusercontent.com
REACT_APP_BACKTRACK_GRAPHQL_URL=http://localhost:9000/graphql`

			err := fileutils.WriteToFile(envPath, envSource)
			errorhandler.CheckNilErr(err)
		}
	}

	if mobileStatus {
		envPath = fmt.Sprintf("%s/%s/%s", fileutils.CurrentDirectory(),
			constants.Mobile,
			dockerEnvFile,
		)
		dockerEnvFound, err = fileutils.IsExists(envPath)
		errorhandler.CheckNilErr(err)
		if !dockerEnvFound {
			envSource = ""

			err := fileutils.WriteToFile(envPath, envSource)
			errorhandler.CheckNilErr(err)
		}
	}

	return nil
}
