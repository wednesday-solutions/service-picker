package helpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func UpdateDBConfig(stack, dbFile string) error {

	var dbConfigSource string

	switch stack {
	case constants.NodeHapiTemplate:
		dbConfigSource = `const pg = require('pg');

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

	case constants.NodeExpressGraphqlTemplate:
		dbConfigSource = fmt.Sprintf(`const Sequelize = require('sequelize');
const mysql2 = require('mysql2');
const dotenv = require('dotenv');

dotenv.config({ path: %s.env.${process.env.ENVIRONMENT_NAME}%s });

module.exports = {
	url: process.env.DB_URI,
	host: process.env.MYSQL_HOST,
	logging: true,
	dialectModule: mysql2,
	dialect: 'mysql',
	pool: {
		min: 0,
		max: 10,
		idle: 10000
	},
	define: {
		userscored: true,
		timestamps: false
	},
	retry: {
		match: [
			'unknown timed out',
			Sequelize.TimeoutError,
			'timed',
			'timeout',
			'TimeoutError',
			'Operation timeout',
			'refuse',
			'SQLITE_BUSY'
		],
		max: 10 // maximum amount of tries
	}
};
`, "`", "`")

	default:
		return fmt.Errorf("Selected stack is invalid")
	}

	err := fileutils.WriteToFile(fileutils.CurrentDirectory()+dbFile, dbConfigSource)
	errorhandler.CheckNilErr(err)

	return nil
}
