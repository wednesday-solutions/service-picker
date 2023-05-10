package pickyhelpers

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/wednesday-solutions/picky/hbs"
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
)

func CreateDockerComposeFile(stackInfo map[string]interface{}) error {

	filePath := fmt.Sprintf("%s/%s", utils.CurrentDirectory(),
		constants.DockerComposeFile,
	)
	var (
		backendMysqlDirectories []string
		backendPgDirectories    []string
		backendMysqlSnakeCased  []string
		backendPgSnakeCased     []string
	)
	var snakeCaseDirName string
	_, databases, directories := utils.GetExistingStacksDatabasesAndDirectories()
	for index, directory := range directories {
		service := utils.FindService(directory)
		if service == constants.Backend {
			if databases[index] == constants.MySQL {
				backendMysqlDirectories = append(backendMysqlDirectories, directory)
				snakeCaseDirName = strcase.ToSnake(directory)
				backendMysqlSnakeCased = append(backendMysqlSnakeCased, snakeCaseDirName)
			} else if databases[index] == constants.PostgreSQL {
				backendPgDirectories = append(backendPgDirectories, directory)
				snakeCaseDirName = strcase.ToSnake(directory)
				backendPgSnakeCased = append(backendPgSnakeCased, snakeCaseDirName)
			}
		}
	}

	// Don't make any changes in the below source string.
	source := `version: '3'
services:`

	for i, d := range backendPgDirectories {
		source = fmt.Sprintf(`%s
  # Setup {{PostgreSQL}}
  %s_db:
    image: '{{dbVersion PostgreSQL}}' 
    ports:
      - {{portConnection PostgreSQL}} 
    restart: always # This will make sure that the container comes up post unexpected shutdowns
    env_file:
      - ./%s/.env.docker
    volumes:
      - %s{{databaseVolumeConnection PostgreSQL}}
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
  # Setup %s api
  %s:
    build:
      context: './%s'
      args:
        ENVIRONMENT_NAME: docker
    ports:
      - {{portConnection backend}}
    env_file:
      - ./%s/.env.docker
    environment:
      ENVIRONMENT_NAME: docker
{{dependsOnFieldOfGo stack}}
`, source, backendPgSnakeCased[i], d, d, d, d, d, d)
	}

	for i, d := range backendMysqlDirectories {
		source = fmt.Sprintf(`%s
  # Setup {{MySQL}}
  %s_db:
    image: '{{dbVersion MySQL}}' 
    ports:
      - {{portConnection MySQL}} 
    restart: always # This will make sure that the container comes up post unexpected shutdowns
    env_file:
      - ./%s/.env.docker
    volumes:
      - %s{{databaseVolumeConnection MySQL}}
{{#equal stack GolangMySQL}}
    environment:
      MYSQL_DATABASE: ${MYSQL_DBNAME}
      MYSQL_PASSWORD: ${MYSQL_PASS}
      MYSQL_ROOT_PASSWORD: ${MYSQL_ROOT_PASSWORD}
{{/equal}}

{{#equal stack GolangMySQL}}
{{{waitForDBService MySQL}}}

{{/equal}}
  # Setup %s api
  %s:
    build:
      context: './%s'
      args:
        ENVIRONMENT_NAME: docker
    ports:
      - {{portConnection backend}}
    env_file:
      - ./%s/.env.docker
    environment:
      ENVIRONMENT_NAME: docker
{{dependsOnFieldOfGo stack}}
`, source, backendMysqlSnakeCased[i], d, d, d, d, d, d)
	}

	source = fmt.Sprintf(`%s
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

{{/each}}{{#if backendStatus}}
# Setup Volumes
volumes:
{{#each backendPgDirectories}}
  {{this}}-db-volume:
{{/each}}
{{#each backendMysqlDirectories}}
  {{this}}-db-volume:
{{/each}}{{/if}}
`, source)

	err := hbs.ParseAndWriteToFile(source, filePath, stackInfo)
	errorhandler.CheckNilErr(err)

	return nil
}
