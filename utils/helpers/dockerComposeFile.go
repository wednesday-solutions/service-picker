package helpers

import (
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
	"github.com/wednesday-solutions/picky/utils/hbs"
)

func CreateDockerComposeFile(database, projectName string, stackStatus map[string]bool) error {

	dockerComposeFile := "docker-compose.yml"
	status, _ := fileutils.IsExists(fileutils.CurrentDirectory() + dockerComposeFile)
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
  {{projectName}}_db:
    image: '{{dbVersion database}}' 
    ports:
      - {{portConnection database}} 
    restart: always # This will make sure that the container comes up post unexpected shutdowns
    env_file:
      - ./backend/.env.docker
    volumes:
      - {{projectName}}_db_volume:/var/lib/{{databaseName database}}/data

  # Setup Redis
  {{projectName}}_redis:
    image: 'redis'
    ports:
      - {{portConnection redis}}
    # Default command that redis will execute at start
    command: [ 'redis-server' ]

  # Setup {{projectName}} API
  {{projectName}}_api:
    build:
      context: './backend'
      args:
        ENVIRONMENT_NAME: docker
    ports:
      - {{portConnection backend}}
    env_file:
      - ./backend/.env.docker
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

	sourceValues := map[string]interface{}{
		"database":      database,
		"projectName":   projectName,
		"webStatus":     stackStatus["webStatus"],
		"mobileStatus":  stackStatus["mobileStatus"],
		"backendStatus": stackStatus["backendStatus"],
	}

	err = hbs.ParseAndWriteToFile(source, dockerComposeFile, sourceValues)
	errorhandler.CheckNilErr(err)

	return nil
}
