package helpers

import (
	"github.com/aymerick/raymond"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func WriteDockerFile(fileName, db, projectName string) error {

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

	ctx := map[string]interface{}{
		"database":    db,
		"redis":       "redis",
		"frontend":    "frontend",
		"backend":     "backend",
		"projectName": projectName,
	}

	tpl, err := raymond.Parse(source)
	errorhandler.CheckNilErr(err)

	result, err := tpl.Exec(ctx)
	errorhandler.CheckNilErr(err)

	err = fileutils.WriteToFile(fileutils.CurrentDirectory(), fileName, result)
	errorhandler.CheckNilErr(err)

	return nil
}

func init() {
	raymond.RegisterHelper("databaseName", DatabaseName)
	raymond.RegisterHelper("dbVersion", DBVersion)
	raymond.RegisterHelper("portConnection", PortConnection)
}

func DatabaseName(db string) string {
	if db == constants.POSTGRES {
		return "postgresql"
	} else if db == constants.MYSQL {
		return "mysql"
	}
	return ""
}

func DBVersion(db string) string {
	if db == constants.POSTGRES {
		return "postgres:15"
	} else if db == constants.MYSQL {
		return "mysql:8"
	}
	return ""
}

func PortConnection(stack string) string {
	switch stack {
	case constants.POSTGRES:
		return "5432:5432"
	case constants.MYSQL:
		return "3306:3306"
	case constants.MONGODB:
		return "27017:27017"
	case constants.FRONTEND:
		return "3000:3000"
	case constants.BACKEND:
		return "9000:9000"
	case "redis":
		return "6379:6379"
	}
	return ""
}
