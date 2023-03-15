package helpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/hbs"
)

func ConvertMysqlToPostgres(stack, database, projectName string) error {

	var queries []string
	var files []string

	switch stack {
	case constants.NODE_HAPI_TEMPLATE:

		queries = []string{`create sequence oauth_clients_seq;

create type grant_type_enum as ENUM('CLIENT_CREDENTIALS');

create table oauth_clients (
	id INT NOT NULL PRIMARY KEY DEFAULT NEXTVAL ('oauth_clients_seq'), 
	client_id VARCHAR(320) NOT NULL, 
	client_secret VARCHAR(36) NOT NULL, 
	grant_type grant_type_enum, 
	created_at TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP(0) NULL,
	CONSTRAINT oauth_clients_client_id UNIQUE (client_id)
);

CREATE INDEX oauth_clients_client_id_idx ON oauth_clients(client_id);
CREATE INDEX oauth_clients_client_secret_idx ON oauth_clients(client_secret);`,

			`CREATE SEQUENCE users_seq;

CREATE TABLE users 
	( 
		id              INT NOT NULL DEFAULT NEXTVAL ('users_seq') PRIMARY KEY, 
		oauth_client_id INT NOT NULL, 
		first_name      VARCHAR (32) NOT NULL, 
		last_name       VARCHAR(32) NOT NULL, 
		email           VARCHAR(32) NOT NULL, 
		created_at      TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP NOT NULL, 
		CONSTRAINT users_oauth_clients_id_fk FOREIGN KEY (oauth_client_id) 
		REFERENCES oauth_clients (id) ON UPDATE CASCADE 
	);`,

			`create sequence oauth_access_tokens_seq;

CREATE OR REPLACE FUNCTION public.is_json_valid(json_data json)
RETURNS boolean
LANGUAGE plpgsql
AS $function$
DECLARE
	json_type TEXT;
BEGIN
	json_type := json_typeof(json_data);
	IF json_type = 'array' AND json_array_length(json_data) > 0 THEN
		RETURN TRUE;
	ELSIF json_type = 'object' AND json_data::text <> '{}'::TEXT THEN
		RETURN TRUE;
	END IF;
	RETURN FALSE;
END;
$function$;

create table oauth_access_tokens (
	id INT NOT NULL PRIMARY KEY DEFAULT NEXTVAL ('oauth_access_tokens_seq'), 
	oauth_client_id INT NOT NULL, 
	access_token VARCHAR(64) NOT NULL, 
	expires_in INTEGER CHECK (expires_in > 0) NOT NULL, 
	expires_on TIMESTAMP(0) NOT NULL, 
	metadata JSON NOT NULL, 
	created_at TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP(0) NULL, 
	CONSTRAINT oauth_access_tokens_access_token_uindex UNIQUE (access_token), 
	CONSTRAINT oauth_access_tokens_oauth_clients_id_fk FOREIGN KEY (oauth_client_id) REFERENCES oauth_clients (id) ON UPDATE CASCADE,
	CONSTRAINT oauth_access_tokens_check_metadata CHECK(is_json_valid(metadata))
);`,

			`create sequence oauth_client_resources_seq;

create table oauth_client_resources (
	id INT NOT NULL PRIMARY KEY DEFAULT NEXTVAL ('oauth_client_resources_seq'), 
	oauth_client_id INT NOT NULL, 
	resource_type VARCHAR(36) NOT NULL, 
	resource_id VARCHAR(36) NOT NULL, 
	created_at TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP(0) NULL, 
	CONSTRAINT oauth_client_resources_oauth_client_id_resource_uindex UNIQUE (
		oauth_client_id, resource_type, resource_id
	), 
	CONSTRAINT oauth_client_resources_oauth_clients_id_fk FOREIGN KEY (oauth_client_id) REFERENCES oauth_clients (id) ON UPDATE CASCADE
);

CREATE INDEX oauth_client_resources_resource_type ON oauth_client_resources(resource_type);
CREATE INDEX oauth_client_resources_resource_id ON oauth_client_resources(resource_id);`,

			`create sequence oauth_client_scopes_seq;

create table oauth_client_scopes (
	id INT NOT NULL DEFAULT NEXTVAL ('oauth_client_scopes_seq') PRIMARY KEY, 
	oauth_client_id INT NOT NULL, 
	scope VARCHAR (36) NOT NULL, 
	created_at TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP NOT NULL, 
	updated_at TIMESTAMP(0) NULL, 
	constraint oauth_client_scopes_uindex UNIQUE (oauth_client_id),
	constraint oauth_client_scopes_oauth_clients_id_fk FOREIGN KEY (oauth_client_id) REFERENCES oauth_clients (id) ON UPDATE CASCADE
);`,
		}

		files = []string{"/backend/resources/v1/01_oauth_clients.sql",
			"/backend/resources/v1/02_users.sql",
			"/backend/resources/v1/03_oauth_access_tokens.sql",
			"/backend/resources/v1/04_oauth_client_resources.sql",
			"/backend/resources/v1/05_oauth_client_scopes.sql",
		}

	case constants.NODE_EXPRESS_GRAPHQL_TEMPLATE:
		queries = []string{}

		files = []string{}

	default:
		return fmt.Errorf("Selected stack is invalid")
	}

	for idx, file := range files {
		err := hbs.ParseAndWriteToFile(queries[idx], constants.POSTGRES, projectName, file)
		errorhandler.CheckNilErr(err)
	}

	return nil
}
