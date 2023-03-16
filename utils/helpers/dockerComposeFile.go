package helpers

import (
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
	"github.com/wednesday-solutions/picky/utils/hbs"
)

func WriteDockerFile(database, projectName string) error {

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

  # Setup {{projectName}} frontend
  {{projectName}}_web:
    build:
      context: './frontend'
    ports:
      - {{portConnection frontend}}
    env_file:
      - ./frontend/.env.docker

# Setup Volumes
volumes:
  {{projectName}}_db_volume:
`

	err = hbs.ParseAndWriteToFile(source, database, projectName, dockerComposeFile)
	errorhandler.CheckNilErr(err)

	return nil
}
