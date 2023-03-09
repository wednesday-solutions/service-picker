package helpers

import (
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/hbs"
)

func UpdateEnvFiles(database, projectName string) error {

	postgresDockerSource := `NAME=Node Template
NODE_ENV=production
ENVIRONMENT_NAME=docker
PORT=9000
DB_URI=postgres://reporting_dashboard_role:reportingdashboard123@db_postgres:5432/reporting_dashboard_dev
POSTGRES_HOST=db_postgres
POSTGRES_DB=reporting_dashboard_dev
POSTGRES_USER=reporting_dashboard_role
POSTGRES_PASSWORD=reportingdashboard123
POSTGRES_PORT=5432
REDIS_HOST=redis`

	mysqlDockerSource := `NAME=Node Template
NODE_ENV=production
ENVIRONMENT_NAME=docker
PORT=9000
DB_URI=mysql://root:password@db_mysql:3306/temp_dev
MYSQL_HOST=db_mysql
MYSQL_DATABASE=temp_dev
MYSQL_USER=def_user
MYSQL_PASSWORD=password
MYSQL_ROOT_PASSWORD=password
ACCESS_TOKEN_SECRET=4cd7234152590dcfe77e1b6fc52e84f4d30c06fddadd0dd2fb42cbc51fa14b1bb195bbe9d72c9599ba0c6b556f9bd1607a8478be87e5a91b697c74032e0ae7af
REDIS_DOMAIN=redis
REDIS_PORT=6379`

	postgresLocalSource := `NAME=Node Template
NODE_ENV=development
ENVIRONMENT_NAME=local
PORT=9000
DB_URI=postgres://reporting_dashboard_role:reportingdashboard123@localhost:5432/reporting_dashboard_dev
POSTGRES_HOST=0.0.0.0
POSTGRES_DB=reporting_dashboard_dev
POSTGRES_USER=reporting_dashboard_role
POSTGRES_PASSWORD=reportingdashboard123
REDIS_HOST=localhost
REDIS_PORT=6379`

	mysqlLocalSource := `NAME=Node Template
NODE_ENV=local
ENVIRONMENT_NAME=local
PORT=9000
DB_URI=mysql://root:password@localhost:3306/temp_dev
MYSQL_HOST=localhost
MYSQL_DATABASE=temp_dev
MYSQL_USER=root
MYSQL_PASSWORD=password
REDIS_DOMAIN=localhost
ACCESS_TOKEN_SECRET=4cd7234152590dcfe77e1b6fc52e84f4d30c06fddadd0dd2fb42cbc51fa14b1bb195bbe9d72c9599ba0c6b556f9bd1607a8478be87e5a91b697c74032e0ae7af`

	postgresDevSource := `NAME=Node Template (DEV)
NODE_ENV=development
ENVIRONMENT_NAME=development
PORT=9000
DB_URI=postgres://reporting_dashboard_role:reportingdashboard123@db_postgres:5432/reporting_dashboard_dev
POSTGRES_HOST=db_postgres
POSTGRES_DB=reporting_dashboard_dev
POSTGRES_USER=reporting_dashboard_role
POSTGRES_PASSWORD=reportingdashboard123
REDIS_HOST=redis`

	mysqlDevSource := `NAME=Node Template (DEV)
NODE_ENV=development
ENVIRONMENT_NAME=development
PORT=9000
DB_URI=mysql://root:password@db_mysql:3306/temp_dev
MYSQL_HOST=db_mysql
MYSQL_DATABASE=temp_dev
MYSQL_USER=root
MYSQL_PASSWORD=password
MYSQL_ROOT_PASSWORD=password
REDIS_HOST=redis
ACCESS_TOKEN_SECRET=4cd7234152590dcfe77e1b6fc52e84f4d30c06fddadd0dd2fb42cbc51fa14b1bb195bbe9d72c9599ba0c6b556f9bd1607a8478be87e5a91b697c74032e0ae7af`

	var envFile string
	var err error
	if database == constants.POSTGRES {
		envFile = "/backend/.env.docker"
		err = hbs.ParseAndWriteToFile(postgresDockerSource, database, projectName, envFile)
		errorhandler.CheckNilErr(err)

		envFile = "/backend/.env.local"
		err = hbs.ParseAndWriteToFile(postgresLocalSource, database, projectName, envFile)
		errorhandler.CheckNilErr(err)

		envFile = "/backend/.env.development"
		err = hbs.ParseAndWriteToFile(postgresDevSource, database, projectName, envFile)

	} else if database == constants.MYSQL {
		envFile = "/backend/.env.docker"
		err = hbs.ParseAndWriteToFile(mysqlDockerSource, database, projectName, envFile)
		errorhandler.CheckNilErr(err)

		envFile = "/backend/.env.local"
		err = hbs.ParseAndWriteToFile(mysqlLocalSource, database, projectName, envFile)
		errorhandler.CheckNilErr(err)

		envFile = "/backend/.env.development"
		err = hbs.ParseAndWriteToFile(mysqlDevSource, database, projectName, envFile)
	}
	errorhandler.CheckNilErr(err)

	return nil
}
