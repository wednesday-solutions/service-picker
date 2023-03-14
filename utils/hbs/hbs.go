package hbs

import (
	"github.com/aymerick/raymond"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func init() {
	raymond.RegisterHelper("databaseName", DatabaseName)
	raymond.RegisterHelper("dbVersion", DBVersion)
	raymond.RegisterHelper("portConnection", PortConnection)
}

func ParseAndWriteToFile(source, db, projectName, fileName string) error {

	ctx := map[string]interface{}{
		"database":    db,
		"redis":       "redis",
		"frontend":    "frontend",
		"backend":     "backend",
		"postgres":    constants.POSTGRES,
		"mysql":       constants.MYSQL,
		"projectName": projectName,
	}

	tpl, err := raymond.Parse(source)
	errorhandler.CheckNilErr(err)

	executedTemplate, err := tpl.Exec(ctx)
	errorhandler.CheckNilErr(err)

	err = fileutils.TruncateAndWriteToFile(fileutils.CurrentDirectory(), fileName, executedTemplate)
	errorhandler.CheckNilErr(err)

	return nil
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
