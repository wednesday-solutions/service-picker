package helpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
	"github.com/wednesday-solutions/picky/utils/hbs"
)

func CreateDockerComposeFile(stackData map[string]interface{}) error {

	dockerComposeFile := "docker-compose.yml"
	filePath := fmt.Sprintf("%s/%s", fileutils.CurrentDirectory(), dockerComposeFile)
	status, _ := fileutils.IsExists(filePath)
	if status {
		return nil
	}

	// create Docker File
	err := fileutils.MakeFile(fileutils.CurrentDirectory(), dockerComposeFile)
	errorhandler.CheckNilErr(err)

	// Don't make any changes in the below source string.
	source := `version: '3'
services:
  # Setup {{database}}
  {{dbServiceName database}}:
    image: '{{dbVersion database}}' 
    ports:
      - {{portConnection database}} 
    restart: always # This will make sure that the container comes up post unexpected shutdowns
    env_file:
      - {{envFileBackend database}}
    volumes:
      - {{projectName}}_db_volume:/var/lib/{{databaseVolume database}}

  # Setup Redis
  redis:
    image: 'redis'
    ports:
      - {{portConnection redis}}
    # Default command that redis will execute at start
    command: ['redis-server']

  # Setup {{projectName}} API
  {{projectName}}_api:
    build:
      context: './backend'
      args:
        ENVIRONMENT_NAME: docker
    ports:
      - {{portConnection backend}}
    env_file:
      - {{envFileBackend database}}
{{#if webStatus}}
 
  # Setup {{projectName}} web
  {{projectName}}_web:
    build:
      context: './web'
    ports:
      - {{portConnection web}}
    env_file:
      - ./web/.env.docker
{{/if}}
{{#if mobileStatus}}

  # Setup {{projectName}} mobile
  {{projectName}}_mobile:
    build:
      context: './mobile'
    ports:
      - {{portConnection mobile}}
    env_file:
      - ./mobile/.env.docker
{{/if}}

# Setup Volumes
volumes:
  {{projectName}}_db_volume:
`

	err = hbs.ParseAndWriteToFile(source, filePath, stackData)
	errorhandler.CheckNilErr(err)

	return nil
}
