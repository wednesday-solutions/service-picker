package pickyhelpers

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
)

func UpdateEnvFiles(stack, dirName string) error {

	var envFileSources []string
	snakeCaseDirName := strcase.ToSnake(dirName)
	envFiles := []string{".env.local", ".env.development", constants.DockerEnvFile}
	for idx, file := range envFiles {
		envFiles[idx] = fmt.Sprintf("%s/%s", dirName, file)
	}
	var envLocalSource, envDevSource, envDockerSource string
	backendPortNumber := utils.GetPortNumber(constants.BackendPortNumber)
	redisPortNumber := utils.GetPortNumber(constants.RedisPortNumber)
	switch stack {
	case constants.NodeHapiTemplate:

		postgresPortNumber := utils.GetDatabasePortNumber(constants.PostgreSQL)
		envLocalSource = fmt.Sprintf(`NAME=Node Template
NODE_ENV=development
ENVIRONMENT_NAME=local
PORT=%d
DB_URI=postgres://root:password@localhost:5432/temp_dev
POSTGRES_HOST=0.0.0.0
POSTGRES_PORT=%d
POSTGRES_DB=temp_dev
POSTGRES_USER=root
POSTGRES_PASSWORD=password
REDIS_HOST=localhost
REDIS_PORT=%d`,
			backendPortNumber,
			postgresPortNumber,
			redisPortNumber,
		)

		envDevSource = fmt.Sprintf(`NAME=Node Template (DEV)
NODE_ENV=development
ENVIRONMENT_NAME=development
PORT=%d
DB_URI=postgres://root:password@db_postgres:5432/temp_dev
POSTGRES_HOST=db_postgres
POSTGRES_PORT=%d
POSTGRES_DB=temp_dev
POSTGRES_USER=root
POSTGRES_PASSWORD=password
REDIS_HOST=redis
REDIS_PORT=%d`,
			backendPortNumber,
			postgresPortNumber,
			redisPortNumber,
		)

		envDockerSource = fmt.Sprintf(`NAME=Node Template
NODE_ENV=production
ENVIRONMENT_NAME=docker
PORT=%d
DB_URI=postgres://root:password@%s_db:5432/temp_dev
POSTGRES_HOST=%s_db
POSTGRES_DB=temp_dev
POSTGRES_USER=root
POSTGRES_PASSWORD=password
POSTGRES_PORT=%d
REDIS_HOST=redis
REDIS_PORT=%d`,
			backendPortNumber,
			snakeCaseDirName,
			snakeCaseDirName,
			postgresPortNumber,
			redisPortNumber,
		)

		envFileSources = []string{envLocalSource, envDevSource, envDockerSource}

	case constants.NodeExpressGraphqlTemplate:

		mysqlPortNumber := utils.GetDatabasePortNumber(constants.MySQL)
		envLocalSource = fmt.Sprintf(`DB_URI=mysql://root:password@localhost:3306/reporting_dashboard_dev
MYSQL_HOST=0.0.0.0
MYSQL_PORT=%d
MYSQL_DATABASE=reporting_dashboard_dev
MYSQL_USER=root
MYSQL_PASSWORD=password
PORT=%d
NODE_ENV=local
ACCESS_TOKEN_SECRET=4cd7234152590dcfe77e1b6fc52e84f4d30c06fddadd0dd2fb42cbc51fa14b1bb195bbe9d72c9599ba0c6b556f9bd1607a8478be87e5a91b697c74032e0ae7af
REDIS_HOST=localhost
REDIS_PORT=%d`,
			mysqlPortNumber,
			backendPortNumber,
			redisPortNumber,
		)

		envDevSource = fmt.Sprintf(`DB_URI=mysql://root:password@db_mysql:3306/reporting_dashboard_dev
MYSQL_HOST=db_mysql
MYSQL_PORT=%d
MYSQL_DATABASE=reporting_dashboard_dev
MYSQL_USER=root
MYSQL_PASSWORD=password
PORT=%d
ACCESS_TOKEN_SECRET=4cd7234152590dcfe77e1b6fc52e84f4d30c06fddadd0dd2fb42cbc51fa14b1bb195bbe9d72c9599ba0c6b556f9bd1607a8478be87e5a91b697c74032e0ae7af`,
			mysqlPortNumber,
			backendPortNumber,
		)

		envDockerSource = fmt.Sprintf(`DB_URI=mysql://root:password@%s_db:3306/reporting_dashboard_dev
MYSQL_HOST=%s_db
MYSQL_DATABASE=reporting_dashboard_dev
MYSQL_USER=reporting_dashboard_role
MYSQL_PASSWORD=password
MYSQL_ROOT_PASSWORD=password
MYSQL_PORT=%d
PORT=%d
NODE_ENV=local
ENVIRONMENT_NAME=docker
ACCESS_TOKEN_SECRET=4cd7234152590dcfe77e1b6fc52e84f4d30c06fddadd0dd2fb42cbc51fa14b1bb195bbe9d72c9599ba0c6b556f9bd1607a8478be87e5a91b697c74032e0ae7af
REDIS_HOST=redis
REDIS_PORT=%d
APP_NAME=app`,
			snakeCaseDirName,
			snakeCaseDirName,
			mysqlPortNumber,
			backendPortNumber,
			redisPortNumber,
		)

		envFileSources = []string{envLocalSource, envDevSource, envDockerSource}

	default:
		return fmt.Errorf("Selected stack is invalid")
	}

	for idx, envFile := range envFiles {
		envFile = fmt.Sprintf("%s/%s", utils.CurrentDirectory(), envFile)
		err := utils.WriteToFile(envFile, envFileSources[idx])
		errorhandler.CheckNilErr(err)
	}

	return nil
}

func UpdateEnvDockerFileForDefaultDBInTemplate(stack, dirName string) error {
	snakeCaseDirName := strcase.ToSnake(dirName)
	var envDockerSource string
	backendPortNumber := utils.GetPortNumber(constants.BackendPortNumber)
	redisPortNumber := utils.GetPortNumber(constants.RedisPortNumber)
	switch stack {
	case constants.NodeExpressGraphqlTemplate:
		postgresPortNumber := utils.GetDatabasePortNumber(constants.PostgreSQL)
		envDockerSource = fmt.Sprintf(`DB_URI=postgres://reporting_dashboard_role:reportingdashboard123@%s_db:5432/reporting_dashboard_dev
POSTGRES_HOST=%s_db
POSTGRES_DB=reporting_dashboard_dev
POSTGRES_USER=reporting_dashboard_role
POSTGRES_PASSWORD=reportingdashboard123
POSTGRES_PORT=%d
ACCESS_TOKEN_SECRET=4cd7234152590dcfe77e1b6fc52e84f4d30c06fddadd0dd2fb42cbc51fa14b1bb195bbe9d72c9599ba0c6b556f9bd1607a8478be87e5a91b697c74032e0ae7af
PORT=%d
NODE_ENV=production
ENVIRONMENT_NAME=docker
REDIS_DOMAIN=redis
REDIS_PORT=%d`,
			snakeCaseDirName,
			snakeCaseDirName,
			postgresPortNumber,
			backendPortNumber,
			redisPortNumber,
		)

	case constants.NodeHapiTemplate:
		mysqlPortNumber := utils.GetDatabasePortNumber(constants.MySQL)
		envDockerSource = fmt.Sprintf(`NAME=Node Template (DEV)
NODE_ENV=production
ENVIRONMENT_NAME=docker
PORT=%d
DB_URI=mysql://root:password@%s_db:3306/temp_dev
MYSQL_HOST=%s_db
MYSQL_DATABASE=temp_dev
MYSQL_USER=def_user
MYSQL_PASSWORD=password
MYSQL_ROOT_PASSWORD=password
MYSQL_PORT=%d
REDIS_HOST=redis
REDIS_PORT=%d`,
			backendPortNumber,
			snakeCaseDirName,
			snakeCaseDirName,
			mysqlPortNumber,
			redisPortNumber,
		)
	default:
		return fmt.Errorf("Selected stack is invalid.")
	}
	file := fmt.Sprintf("%s/%s/%s", utils.CurrentDirectory(),
		dirName, constants.DockerEnvFile,
	)
	err := utils.WriteToFile(file, envDockerSource)
	errorhandler.CheckNilErr(err)

	return nil
}
