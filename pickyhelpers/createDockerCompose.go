package pickyhelpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/hbs"
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
)

func CreateDockerComposeFile(stackInfo map[string]interface{}) error {

	filePath := fmt.Sprintf("%s/%s", utils.CurrentDirectory(),
		constants.DockerComposeFile,
	)
	// Don't make any changes in the below source string.
	_ = `version: '3'
services:
{{#if backendStatus}}
  # Setup {{database}}
  {{dbServiceName stack database}}:
    image: '{{dbVersion database}}' 
    ports:
      - {{portConnection database}} 
    restart: always # This will make sure that the container comes up post unexpected shutdowns
    env_file:
      - ./{{backendDirName}}/.env.docker
    volumes:
      - {{projectName}}_db_volume:/var/lib/{{databaseVolume database}}
{{#equal stack GolangPostgreSQL}}
    environment:
      POSTGRES_USER: ${PSQL_USER}
      POSTGRES_PASSWORD: ${PSQL_PASS}
      POSTGRES_DB: ${PSQL_DBNAME}
      POSTGRES_PORT: ${PSQL_PORT}
{{/equal}}
{{#equal stack GolangMySQL}}
    environment:
      MYSQL_DATABASE: ${MYSQL_DBNAME}
      MYSQL_PASSWORD: ${MYSQL_PASS}
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
{{/equal}}

  # Setup Redis
  redis:
    image: 'redis:6-alpine'
    ports:
      - {{portConnection redis}}
    # Default command that redis will execute at start
    command: ['redis-server']

{{#equal stack GolangPostgreSQL}}
{{{waitForDBService database}}}

{{/equal}}
{{#equal stack GolangMySQL}}
{{{waitForDBService database}}}

{{/equal}}
  # Setup {{projectName}} API
  {{projectName}}_api:
    build:
      context: './{{backendDirName}}'
      args:
        ENVIRONMENT_NAME: docker
    ports:
      - {{portConnection backend}}
    env_file:
      - ./{{backendDirName}}/.env.docker
    environment:
      ENVIRONMENT_NAME: docker
{{dependsOnFieldOfGo stack}}
{{/if}}
{{#each webDirectories}} 
  # Setup {{projectName}} {{this}}
  {{projectName}}_{{this}}:
    build:
      context: './{{this}}'
    ports:
      - {{portConnection web}}
    env_file:
      - ./{{this}}/.env.docker

{{else}}
# No web directories

{{/each}}
# Setup Volumes
volumes:
  {{projectName}}_db_volume:
`

	source := `version: '3'
services:
{{#each backendPgDirectories}}
  # Setup {{PostgreSQL}}
  {{dbServiceName stack PostgreSQL}}:
    image: '{{dbVersion PostgreSQL}}' 
    ports:
      - {{portConnection PostgreSQL}} 
    restart: always # This will make sure that the container comes up post unexpected shutdowns
    env_file:
      - ./{{this}}/.env.docker
    volumes:
      - {{this}}{{databaseVolumeConnection PostgreSQL}}
{{#equal stack GolangPostgreSQL}}
    environment:
      POSTGRES_USER: ${PSQL_USER}
      POSTGRES_PASSWORD: ${PSQL_PASS}
      POSTGRES_DB: ${PSQL_DBNAME}
      POSTGRES_PORT: ${PSQL_PORT}
{{/equal}}

{{#equal stack GolangPostgreSQL}}
{{{waitForDBService PostgreSQL}}}

{{/equal}}
  # Setup {{this}} api
  {{this}}:
    build:
      context: './{{this}}'
      args:
        ENVIRONMENT_NAME: docker
    ports:
      - {{portConnection backend}}
    env_file:
      - ./{{this}}/.env.docker
    environment:
      ENVIRONMENT_NAME: docker
{{dependsOnFieldOfGo stack}}
{{/each}}

{{#each backendMysqlDirectories}}
  # Setup {{MySQL}}
  {{dbServiceName stack MySQL}}:
    image: '{{dbVersion MySQL}}' 
    ports:
      - {{portConnection MySQL}} 
    restart: always # This will make sure that the container comes up post unexpected shutdowns
    env_file:
      - ./{{this}}/.env.docker
    volumes:
      - {{this}}{{databaseVolumeConnection MySQL}}
{{#equal stack GolangMySQL}}
    environment:
      MYSQL_DATABASE: ${MYSQL_DBNAME}
      MYSQL_PASSWORD: ${MYSQL_PASS}
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
{{/equal}}

{{#equal stack GolangMySQL}}
{{{waitForDBService MySQL}}}

{{/equal}}
  # Setup {{this}} api
  {{this}}:
    build:
      context: './{{this}}'
      args:
        ENVIRONMENT_NAME: docker
    ports:
      - {{portConnection backend}}
    env_file:
      - ./{{this}}/.env.docker
    environment:
      ENVIRONMENT_NAME: docker
{{dependsOnFieldOfGo stack}}
{{/each}}
{{#if backendStatus}}
  # Setup Redis
  redis:
    image: 'redis:6-alpine'
    ports:
      - {{portConnection redis}}
    # Default command that redis will execute at start
    command: ['redis-server']

{{/if}}
{{#each webDirectories}} 
  # Setup {{this}} web
  {{this}}:
    build:
      context: './{{this}}'
    ports:
      - {{portConnection web}}
    env_file:
      - ./{{this}}/.env.docker

{{else}}
# No web directories

{{/each}}
# Setup Volumes
volumes:
{{#each backendPgDirectories}}
  {{this}}-db-volume:
{{/each}}
{{#each backendMysqlDirectories}}
  {{this}}-db-volume:
{{/each}}
`

	err := hbs.ParseAndWriteToFile(source, filePath, stackInfo)
	errorhandler.CheckNilErr(err)

	return nil
}
