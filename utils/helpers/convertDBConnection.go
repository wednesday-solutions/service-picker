package helpers

import (
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/hbs"
)

func ConvertDBConnection(stack, dbFile, database, projectName string) error {

	switch stack {
	case constants.NODE_HAPI:
		postgresSource := `module.exports = {
	url: process.env.DB_URI,
	logging: true,
	options: {
		dialect: 'postgres',
		pool: {
			min: 0,
			max: 10,
			idle: 10000
		},
		define: {
			userscored: true,
			timestamps: false
		}
	}
};
`

		mysqlSource := `const mysql2 = require('mysql2');

module.exports = {
	url: process.env.DB_URI,
	host: process.env.MYSQL_HOST,
	dialectModule: mysql2,
	dialect: 'mysql',
	pool: {
		min: 0,
		max: 10,
		idle: 10000,
	},
	define: {
		userscored: true,
		timestamps: false,
	},
};
`

		var err error
		if database == constants.POSTGRES {
			err = hbs.ParseAndWriteToFile(postgresSource, database, projectName, dbFile)
		} else if database == constants.MYSQL {
			err = hbs.ParseAndWriteToFile(mysqlSource, database, projectName, dbFile)
		}
		errorhandler.CheckNilErr(err)

		return nil
	default:
		return nil
	}
}
