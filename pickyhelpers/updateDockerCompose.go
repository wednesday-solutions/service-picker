package pickyhelpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/hbs"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func UpdateDockerCompose(stack, dirName string, stackInfo map[string]interface{}) error {
	var updateDockerCompose bool
	source := `version: '3'
services:
  {{dbServiceName stack database}}:
    image: {{dbVersion database}}
    ports:
      - {{portConnection database}}
    restart: always
    env_file:
      - .env.docker

  redis:
    image: 'redis:6-alpine'
    ports:
      - {{portConnection redis}}
    command: ['redis-server', '--bind', 'redis', '--port', '6379']

  app:
    build:
      context: .
      args:
        ENVIRONMENT_NAME: docker
    depends_on:
      - redis
      - {{dbServiceName stack database}}
    restart: always
    ports:
      - {{portConnection backend}}
    env_file:
      - .env.docker
`
	switch stack {
	case constants.NodeExpressGraphqlTemplate, constants.NodeHapiTemplate:
		updateDockerCompose = true
	default:
		updateDockerCompose = false
	}
	if updateDockerCompose {
		path := fmt.Sprintf("%s/%s/%s", fileutils.CurrentDirectory(),
			dirName,
			constants.DockerComposeFile,
		)
		err := hbs.ParseAndWriteToFile(source, path, stackInfo)
		errorhandler.CheckNilErr(err)
	}
	return nil
}
