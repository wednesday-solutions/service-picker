package helpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func UpdateDBIndexToPostgres(stack, indexFile, database, projectName string) error {

	indexFile = fmt.Sprintf("%s%s", fileutils.CurrentDirectory(), indexFile)
	var postgresSource string
	switch stack {
	case constants.NODE_HAPI:
		postgresSource = fmt.Sprintf(`import Sequelize from 'sequelize';
import SequelizeMock from 'sequelize-mock';
import * as pg from 'pg';
import oauthAccessTokens from './oauthAccessTokens';
import oauthClientResources from './oauthClientResources';
import oauthClientScopes from './oauthClientScopes';
import oauthClients from './oauthClients';
import users from './users';
import { isTestEnv, logger } from '@utils';
const cls = require('cls-hooked');

let client;
let namespace;
export const getClient = (force) => {
	if (!namespace) {
		namespace = cls.createNamespace(
			%s${process.env.ENVIRONMENT_NAME}-namespace%s
		);
	}
	if (force || !client) {
		try {
			if (!isTestEnv()) {
				Sequelize.useCLS(namespace);
			}
			client = new Sequelize(
				process.env.POSTGRES_DB,
				process.env.POSTGRES_USER,
				process.env.POSTGRES_PASSWORD,
				{
					host: process.env.POSTGRES_HOST,
					dialect: 'postgres',
					dialectModule: pg,
					pool: {
						min: 0,
						max: 10,
						idle: 10000,
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
							'SQLITE_BUSY',
						],
						max: 10, // maximum amount of tries
					},
				}
			);
		} catch (err) {
			logger().info({ err });
			throw err;
		}
	}
	return client;
};

export const connect = async () => {
	client = getClient();
	try {
		await client.authenticated();
		console.log('Connection has been established successfully. \n', {
			db: process.env.POSTGRES_DB,
			user: process.env.POSTGRES_USER,
			host: process.env.POSTGRES_HOST,
		});
	} catch (error) {
		console.error('Unable to connect to the database:', error);
		throw error;
	}
};

let sequelize;
if (process.env.NODE_ENV === 'test') {
	sequelize = new SequelizeMock();
} else {
	sequelize = getClient();
}

export const models = {
	oauthAccessTokens: oauthAccessTokens(sequelize, Sequelize.DataTypes),
	oauthClientResources: oauthClientResources(sequelize, Sequelize.DataTypes),
	oauthClientScopes: oauthClientScopes(sequelize, Sequelize.DataTypes),
	oauthClients: oauthClients(sequelize, Sequelize.DataTypes),
	users: users(sequelize, Sequelize.DataTypes),
};

Object.keys(models).forEach((modelName) => {
	if (models[modelName].associate) {
		models[modelName].associate(models);
	}
});

models.sequelize = sequelize;
models.Sequelize = Sequelize;

export default models;
`, "`", "`")

		if database == constants.POSTGRES {
			err := fileutils.WriteToFile(indexFile, postgresSource)
			errorhandler.CheckNilErr(err)
		}
	}

	return nil
}
