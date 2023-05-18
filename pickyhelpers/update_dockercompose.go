package pickyhelpers

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/wednesday-solutions/picky/hbs"
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
)

func UpdateDockerCompose(stack, dirName string, stackInfo map[string]interface{}) error {
	var updateDockerCompose bool
	snakeCaseDirName := strcase.ToSnake(dirName)
	source := fmt.Sprintf(`version: '3'
services:
  %s_db:
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
      - %s_db
    restart: always
    ports:
      - {{portConnection backend}}
    env_file:
      - .env.docker
`, snakeCaseDirName, snakeCaseDirName)
	switch stack {
	case constants.NodeExpressGraphqlTemplate, constants.NodeHapiTemplate:
		updateDockerCompose = true
	default:
		updateDockerCompose = false
	}
	if updateDockerCompose {
		path := fmt.Sprintf("%s/%s/%s", utils.CurrentDirectory(),
			dirName,
			constants.DockerComposeFile,
		)
		err := hbs.ParseAndWriteToFile(source, path, stackInfo)
		errorhandler.CheckNilErr(err)
	}
	return nil
}
