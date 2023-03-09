package helpers

import (
	"github.com/aymerick/raymond"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func WriteDockerFile(fileName, db, projectName string) error {

	source := `version: '3'
services:
	# Setup {{database}}
	{{projectName}}_db:
		image: {{{dbVersionAndPort database}}}
		restart: always # This will make sure that the container comes up post unexpected shutdowns
		env_file:
			- ./backend/.env.docker
		volumes:
			- {{projectName}}_db_volume:/var/lib/{{databaseName database}}/data

	# Setup Redis
	{{projectName}}_redis:
		image: 'redis'
		ports:
			- 6379:6379
		# Default command that redis will execute at start
		command: [ 'redis-server' ]

	# Setup {{projectName}} API
	{{projectName}}_api:
		build:
			context: './backend'
			args:
				ENVIRONMENT_NAME: docker
		ports:
			- 9000:9000
		env_file:
			- ./backend/.env.docker

	# Setup {{projectName}} frontend
	{{projectName}}_web:
		build:
			context: './frontend'
		ports:
			- 3000:3000
		env_file:
			- ./frontend/.env.docker

# Setup Volumes
volumes:
	{{projectName}}_db_volume:
`

	ctx := map[string]interface{}{
		"database":    db,
		"projectName": projectName,
	}

	tpl, err := raymond.Parse(source)
	if err != nil {
		return err
	}
	result, err := tpl.Exec(ctx)
	if err != nil {
		return err
	}
	err = fileutils.WriteToFile(fileutils.CurrentDirectory(), fileName, result)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	raymond.RegisterHelper("databaseName", DatabaseName)
	raymond.RegisterHelper("dbVersionAndPort", DBVersionAndPort)
}

func DatabaseName(db string) string {
	if db == "postgres" {
		return "postgresql"
	} else if db == "mysql" {
		return "mysql"
	}
	return ""
}

func DBVersionAndPort(db string) string {
	if db == "postgres" {
		return `'postgres:15'
		ports:
			- 5432:5432`
	} else if db == "mysql" {
		return `'mysql:8'
		ports:
			- 3306:3306`
	}
	return ""
}
