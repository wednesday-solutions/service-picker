package helpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/hbs"
)

func ConvertQueries(stack, database, projectName string) error {

	var queries []string
	var files []string

	switch stack {
	case constants.NODE_HAPI_TEMPLATE:

		oauthClients := `create sequence oauth_clients_seq;

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
CREATE INDEX oauth_clients_client_secret_idx ON oauth_clients(client_secret);`

		users := `CREATE SEQUENCE users_seq;

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
	);`

		oauthAccessTokens := `create sequence oauth_access_tokens_seq;

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
);`

		oauthClientResources := `create sequence oauth_client_resources_seq;

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
CREATE INDEX oauth_client_resources_resource_id ON oauth_client_resources(resource_id);`

		oauthClientScopes := `create sequence oauth_client_scopes_seq;

create table oauth_client_scopes (
	id INT NOT NULL DEFAULT NEXTVAL ('oauth_client_scopes_seq') PRIMARY KEY, 
	oauth_client_id INT NOT NULL, 
	scope VARCHAR (36) NOT NULL, 
	created_at TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP NOT NULL, 
	updated_at TIMESTAMP(0) NULL, 
	constraint oauth_client_scopes_uindex UNIQUE (oauth_client_id),
	constraint oauth_client_scopes_oauth_clients_id_fk FOREIGN KEY (oauth_client_id) REFERENCES oauth_clients (id) ON UPDATE CASCADE
);`

		queries = []string{oauthClients,
			users,
			oauthAccessTokens,
			oauthClientResources,
			oauthClientScopes,
		}

		files = []string{"/backend/resources/v1/01_oauth_clients.sql",
			"/backend/resources/v1/02_users.sql",
			"/backend/resources/v1/03_oauth_access_tokens.sql",
			"/backend/resources/v1/04_oauth_client_resources.sql",
			"/backend/resources/v1/05_oauth_client_scopes.sql",
		}

	case constants.NODE_EXPRESS_GRAPHQL_TEMPLATE:

		products := `CREATE TABLE products (
	id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	category VARCHAR(255) NOT NULL,
	amount BIGINT NOT NULL,
	created_at DATETIME DEFAULT NOW(),
	updated_at DATETIME NULL on UPDATE NOW(),
	deleted_at DATETIME,
	INDEX(name),
	INDEX(category)
);`

		addresses := `CREATE TABLE addresses (
	id INT AUTO_INCREMENT PRIMARY KEY,
	address_1 VARCHAR(255) NOT NULL,
	address_2 VARCHAR(255) NOT NULL,
	city VARCHAR(255) NOT NULL,
	country VARCHAR(255) NOT NULL,
	latitude FLOAT NOT NULL,
	longitude FLOAT NOT NULL,
	created_at DATETIME DEFAULT NOW(),
	updated_at DATETIME NULL on UPDATE NOW(),
	deleted_at DATETIME,
	INDEX(latitude),
	INDEX(longitude)
);`

		stores := `CREATE TABLE stores (
	id INT AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	address_id INT NOT NULL,
	created_at DATETIME DEFAULT NOW(),
	updated_at DATETIME NULL on UPDATE NOW(),
	deleted_at DATETIME,
	CONSTRAINT stores_address_id FOREIGN KEY (address_id) REFERENCES addresses (id),
	INDEX(name)
);`

		suppliers := `CREATE TABLE suppliers (
	id INT AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	address_id INT NOT NULL,
	created_at DATETIME DEFAULT NOW(),
	updated_at DATETIME NULL on UPDATE NOW(),
	deleted_at DATETIME,
	CONSTRAINT suppliers_address_id FOREIGN KEY (address_id) REFERENCES addresses (id),
	INDEX(name)
);`

		supplierProducts := `CREATE TABLE supplier_products (
	id INT AUTO_INCREMENT PRIMARY KEY,
	product_id INT NOT NULL,
	supplier_id INT NOT NULL,
	created_at DATETIME DEFAULT NOW(),
	updated_at DATETIME NULL on UPDATE NOW(),
	deleted_at DATETIME,
	CONSTRAINT suppliers_product_products_id FOREIGN KEY (product_id) REFERENCES products (id),
	CONSTRAINT suppliers_product_supplier_id FOREIGN KEY (supplier_id) REFERENCES suppliers (id)
);`

		storeProducts := `CREATE TABLE store_products (
	id INT AUTO_INCREMENT PRIMARY KEY,
	product_id INT NOT NULL,
	store_id INT NOT NULL,
	created_at DATETIME DEFAULT NOW(),
	updated_at DATETIME NULL on UPDATE NOW(),
	deleted_at DATETIME,
	CONSTRAINT store_products_product_id FOREIGN KEY (product_id) REFERENCES products (id),
	CONSTRAINT store_products_store_id FOREIGN KEY (store_id) REFERENCES stores (id)
);`

		purchasedProducts := `CREATE TABLE purchased_products (
	id INT AUTO_INCREMENT PRIMARY KEY,
	product_id INT NOT NULL,
	price INT NOT NULL,
	discount INT NOT NULL,
	store_id INT NOT NULL,
	delivery_date DATETIME NOT NULL,
	created_at DATETIME DEFAULT NOW(),
	updated_at DATETIME NULL on UPDATE NOW(),
	deleted_at DATETIME,
	CONSTRAINT purchased_products_product_id FOREIGN KEY (product_id) REFERENCES products (id),
	CONSTRAINT purchased_products_store_id FOREIGN KEY (store_id) REFERENCES stores (id),
	INDEX(delivery_date),
	INDEX(store_id)
);`

		users := `CREATE TABLE users (
	id INT AUTO_INCREMENT PRIMARY KEY,
	first_name VARCHAR(255) NOT NULL,
	last_name VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL UNIQUE,
	password VARCHAR(255) NOT NULL,
	created_at DATETIME DEFAULT NOW(),
	updated_at DATETIME NULL on UPDATE NOW(),
	deleted_at DATETIME,
	INDEX(email)
);`

		queries = []string{products,
			addresses,
			stores,
			suppliers,
			supplierProducts,
			storeProducts,
			purchasedProducts,
			users,
		}

		files = []string{"/backend/resources/v1/01_products.sql",
			"/backend/resources/v1/02_addresses.sql",
			"/backend/resources/v1/03_stores.sql",
			"/backend/resources/v1/04_supplier.sql",
			"/backend/resources/v1/05_supplier_products.sql",
			"/backend/resources/v1/06_store_products.sql",
			"/backend/resources/v1/07_purchased_products.sql",
			"/backend/resources/v1/08_users.sql",
		}

	default:
		return fmt.Errorf("Selected stack is invalid")
	}

	for idx, file := range files {
		err := hbs.ParseAndWriteToFile(queries[idx], constants.POSTGRES, projectName, file)
		errorhandler.CheckNilErr(err)
	}

	return nil
}
