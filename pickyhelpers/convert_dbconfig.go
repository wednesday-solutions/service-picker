package pickyhelpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/utils"
)

func ConvertDBConfig(stack, dirName string) error {

	if stack == constants.NodeExpressGraphqlTemplate {

		file := fmt.Sprintf("%s/%s/%s/%s/%s/%s",
			utils.CurrentDirectory(),
			dirName,
			"server",
			"utils",
			"testUtils",
			"dbConfig.js",
		)
		dbUri := "mysql://user:password@host:3306/table"
		source := fmt.Sprintf(`export const DB_ENV = {
	DB_URI: '%s',
	%s: 'host',
	%s: 'user',
	%s: 'password',
	%s: 'table'
};
`,
			dbUri,
			constants.MysqlHost,
			constants.MysqlUser,
			"MYSQL_PASSWORD",
			constants.MysqlDatabase,
		)
		err := utils.WriteToFile(file, source)
		return err
	}
	return nil
}
