package helpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/hbs"
)

func UpdateEnvFiles(stack, database, projectName string) error {

	var envFileSources []string
	envFiles := []string{"/backend/.env.local", "/backend/.env.development", "/backend/.env.docker"}

	switch stack {
	case constants.NODE_HAPI_TEMPLATE:

		envFileSources = []string{`NAME=Node Template
NODE_ENV=development
ENVIRONMENT_NAME=local
PORT=9000
DB_URI=postgres://root:password@localhost:5432/temp_dev
POSTGRES_HOST=0.0.0.0
POSTGRES_DB=temp_dev
POSTGRES_USER=root
POSTGRES_PASSWORD=password
REDIS_HOST=localhost
REDIS_PORT=6379`,

			`NAME=Node Template (DEV)
NODE_ENV=development
ENVIRONMENT_NAME=development
PORT=9000
DB_URI=postgres://root:password@db_postgres:5432/temp_dev
POSTGRES_HOST=db_postgres
POSTGRES_DB=temp_dev
POSTGRES_USER=root
POSTGRES_PASSWORD=password
REDIS_HOST=redis`,

			`NAME=Node Template
NODE_ENV=production
ENVIRONMENT_NAME=docker
PORT=9000
DB_URI=postgres://root:password@db_postgres:5432/temp_dev
POSTGRES_HOST=db_postgres
POSTGRES_DB=temp_dev
POSTGRES_USER=root
POSTGRES_PASSWORD=password
POSTGRES_PORT=5432
REDIS_HOST=redis`,
		}

	case constants.NODE_EXPRESS_GRAPHQL_TEMPLATE:

		envFileSources = []string{`DB_URI=mysql://root:password@localhost:3306/reporting_dashboard_dev
MYSQL_HOST=0.0.0.0
MYSQL_DB=reporting_dashboard_dev
MYSQL_USER=root
MYSQL_PASSWORD=password
NODE_ENV=local
ACCESS_TOKEN_SECRET=4cd7234152590dcfe77e1b6fc52e84f4d30c06fddadd0dd2fb42cbc51fa14b1bb195bbe9d72c9599ba0c6b556f9bd1607a8478be87e5a91b697c74032e0ae7af
REDIS_DOMAIN=localhost
REDIS_PORT=6379`,

			`DB_URI=mysql://root:password@db_mysql:5432/reporting_dashboard_dev
MYSQL_HOST=db_mysql
MYSQL_DB=reporting_dashboard_dev
MYSQL_USER=root
MYSQL_PASSWORD=password
ACCESS_TOKEN_SECRET=4cd7234152590dcfe77e1b6fc52e84f4d30c06fddadd0dd2fb42cbc51fa14b1bb195bbe9d72c9599ba0c6b556f9bd1607a8478be87e5a91b697c74032e0ae7af`,

			`DB_URI=mysql://root:password@db_mysql:5432/reporting_dashboard_dev
MYSQL_HOST=db_mysql
MYSQL_DB=reporting_dashboard_dev
MYSQL_USER=root
MYSQL_PASSWORD=password
MYSQL_PORT=5432
NODE_ENV=production
ENVIRONMENT_NAME=docker
ACCESS_TOKEN_SECRET=4cd7234152590dcfe77e1b6fc52e84f4d30c06fddadd0dd2fb42cbc51fa14b1bb195bbe9d72c9599ba0c6b556f9bd1607a8478be87e5a91b697c74032e0ae7af
REDIS_DOMAIN=redis
REDIS_PORT=6379`,
		}

	default:
		return fmt.Errorf("Selected stack is invalid")
	}

	for idx, envFile := range envFiles {
		err := hbs.ParseAndWriteToFile(envFileSources[idx], database, projectName, envFile)
		errorhandler.CheckNilErr(err)
	}

	return nil
}
