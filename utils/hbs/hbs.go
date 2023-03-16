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

func ParseAndWriteToFile(source, database, projectName, fileName string) error {

	ctx := map[string]interface{}{
		"database":    database,
		"redis":       "redis",
		"frontend":    "frontend",
		"web":         "web",
		"mobile":      "mobile",
		"backend":     "backend",
		"postgres":    "postgres",
		"mysql":       "mysql",
		"projectName": projectName,
	}

	// Parse the source string into template
	tpl, err := raymond.Parse(source)
	errorhandler.CheckNilErr(err)

	// Execute the template into string
	executedTemplate, err := tpl.Exec(ctx)
	errorhandler.CheckNilErr(err)

	err = fileutils.TruncateAndWriteToFile(fileutils.CurrentDirectory(), fileName, executedTemplate)
	errorhandler.CheckNilErr(err)

	return nil
}

func DatabaseName(db string) string {
	if db == constants.PostgreSQL {
		return "postgresql"
	} else if db == constants.MySQL {
		return "mysql"
	} else {
		return ""
	}
}

func DBVersion(db string) string {
	if db == constants.PostgreSQL {
		return "postgres:15"
	} else if db == constants.MySQL {
		return "mysql:8"
	} else {
		return ""
	}
}

func PortConnection(stack string) string {
	switch stack {
	case constants.PostgreSQL:
		return "5432:5432"
	case constants.MySQL:
		return "3306:3306"
	case constants.MongoDB:
		return "27017:27017"
	case constants.Web, constants.Mobile:
		return "3000:3000"
	case constants.Backend:
		return "9000:9000"
	case "redis":
		return "6379:6379"
	default:
		return ""
	}
}
