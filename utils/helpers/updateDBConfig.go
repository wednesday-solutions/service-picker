package helpers

import (
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/hbs"
)

func UpdateDBConfig(stack, dbFile, database, projectName string) error {

	switch stack {
	case constants.NODE_HAPI_TEMPLATE:
		postgresSource := `const pg = require('pg');

module.exports = {
	url: process.env.DB_URI,
	host: process.env.POSTGRES_HOST,
	dialectModule: pg,
	logging: true,
	dialect: 'postgres',
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

		mysqlSource := `const mysql2 = require('mysql2');

module.exports = {
	url: process.env.DB_URI,
	host: process.env.MYSQL_HOST,
	dialectModule: mysql2,
	logging: true,
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
