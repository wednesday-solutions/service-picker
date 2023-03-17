package helpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func CreateDockerFiles(stackStatus map[string]bool) error {

	var (
		path      string
		source    string
		fileFound bool
		err       error
	)
	dockerfile := "Dockerfile"
	dockerEnvFile := ".env.docker"

	if stackStatus["webStatus"] {

		path = fmt.Sprintf("%s/%s/%s", fileutils.CurrentDirectory(),
			constants.Web,
			dockerfile,
		)
		fileFound, _ = fileutils.IsExists(path)
		if !fileFound {
			source = `FROM node:14-alpine as baseimage
RUN mkdir app/
ADD . app/
WORKDIR app/

# RUN npm install

FROM baseimage
CMD [ "yarn", "start:docker" ]
EXPOSE 3000`

			err = fileutils.WriteToFile(path, source)
			errorhandler.CheckNilErr(err)
		}
		path = fmt.Sprintf("%s/%s/%s", fileutils.CurrentDirectory(),
			constants.Web,
			dockerEnvFile,
		)
		fileFound, _ = fileutils.IsExists(path)
		if !fileFound {
			source = `REACT_APP_GOOGLE_LOGIN_CLIENT_ID=775156772157-nlfs5bn4skllfemcvhr6kbphp61achn7.apps.googleusercontent.com
REACT_APP_GRAPHQL_URL=http://localhost:9000/graphql`

			err = fileutils.WriteToFile(path, source)
			errorhandler.CheckNilErr(err)
		}
	}

	if stackStatus["mobileStatus"] {

		path = fmt.Sprintf("%s/%s/%s", fileutils.CurrentDirectory(),
			constants.Mobile,
			dockerfile,
		)
		fileFound, _ = fileutils.IsExists(path)
		if !fileFound {
			source = `FROM baseimage
RUN mkdir app/
ADD . app/
WORKDIR app/

# RUN install

FROM baseimage
CMD [ "yarn", "start:docker" ]
EXPOSE 3000
`

			err = fileutils.WriteToFile(path, source)
			errorhandler.CheckNilErr(err)
		}
		path = fmt.Sprintf("%s/%s/%s", fileutils.CurrentDirectory(),
			constants.Mobile,
			dockerEnvFile,
		)
		fileFound, _ = fileutils.IsExists(path)
		if !fileFound {
			source = ""

			err = fileutils.WriteToFile(path, source)
			errorhandler.CheckNilErr(err)
		}
	}

	if stackStatus["backendStatus"] {

		path = fmt.Sprintf("%s/%s/%s", fileutils.CurrentDirectory(),
			constants.Backend,
			dockerfile,
		)
		fileFound, _ = fileutils.IsExists(path)
		if fileFound {
			source = ``
			err = fileutils.AppendToFile(path, source)
			errorhandler.CheckNilErr(err)
		}
		path = fmt.Sprintf("%s/%s/%s", fileutils.CurrentDirectory(),
			constants.Backend,
			dockerEnvFile,
		)
		fileFound, _ = fileutils.IsExists(path)
		if fileFound {
			source = `
APP_NAME=app
APP_API_DOCKER_CLUSTER_SECRET='{"username": "admin", "password":"password", "host": "app_db", "port": 5432, "dbname": "app_db"}'`

			err = fileutils.AppendToFile(path, source)
			errorhandler.CheckNilErr(err)
		}
	}

	return nil
}
