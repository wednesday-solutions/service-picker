package pickyhelpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/hbs"
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
)

func UpdateDBConfig(stack, dbFile string, stackInfo map[string]interface{}) error {

	var dbConfigSource string

	switch stack {
	case constants.NodeHapiTemplate:
		dbConfigSource = fmt.Sprintf(`const pg = require('pg');

module.exports = {
	url:
		process.env.DB_URI ||
		%spostgres://${process.env.POSTGRES_USER}:${process.env.POSTGRES_PASSWORD}@${process.env.POSTGRES_HOST}/${
			process.env.POSTGRES_DB
		}%s,
	host: process.env.POSTGRES_HOST,
	dialectModule: pg,
	dialect: 'postgres',
	pool: {
		min: 0,
		max: 10,
		idle: 10000,
	},
	define: {
		underscored: true,
		timestamps: false,
	},
};
`, "`", "`")

	case constants.NodeExpressGraphqlTemplate:
		dbConfigSource = fmt.Sprintf(`const Sequelize = require('sequelize');
const mysql2 = require('mysql2');
const dotenv = require('dotenv');

dotenv.config({ path: {{envEnvironmentName}} });

module.exports = {
  url:
    process.env.DB_URI ||
    %smysql://${process.env.MYSQL_USER}:${process.env.MYSQL_PASSWORD}@${process.env.MYSQL_HOST}/${
      process.env.MYSQL_DATABASE
    }%s,
  host: process.env.MYSQL_HOST,
  dialectModule: mysql2,
  dialect: 'mysql',
  pool: {
    min: 0,
    max: 10,
    idle: 10000
  },
  define: {
    underscored: true,
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
    max: 10 // maximum amount of tries.
  }
};
`, "`", "`")

	default:
		return fmt.Errorf("Selected stack is invalid")
	}

	err := hbs.ParseAndWriteToFile(dbConfigSource, dbFile, stackInfo)
	errorhandler.CheckNilErr(err)

	return nil
}
