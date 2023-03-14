package constants

var Repos = func() map[string]string {
	return map[string]string{
		"react":                       "https://github.com/wednesday-solutions/react-template",
		"next":                        "https://github.com/wednesday-solutions/nextjs-template",
		"Node (Hapi- REST API)":       "https://github.com/wednesday-solutions/nodejs-hapi-template",
		"Node (Express- GraphQL API)": "https://github.com/wednesday-solutions/node-express-graphql-template",
		"Node (Express- TypeScript)":  "",
		"Golang-postgres":             "https://github.com/wednesday-solutions/go-template",
		"Golang-mysql":                "https://github.com/wednesday-solutions/go-template-mysql",
	}
}

// CLI options
const (
	PICKY   = "picky"
	SERVICE = "service"
)

// Services
const (
	WEB      = "web"
	MOBILE   = "mobile"
	BACKEND  = "backend"
	FRONTEND = "frontend"
	DATABASE = "database"
)

// Frontend stacks
const (
	REACT = "react"
	NEXT  = "next"
)

// Backend stacks
const (
	NODE_HAPI       = "Node (Hapi- REST API)"
	NODE_EXPRESS    = "Node (Express- GraphQL API)"
	NODE_EXPRESS_TS = "Node (Express- TypeScript)"
	GOLANG          = "Golang (Echo- GraphQL API)"
)

// Databases
const (
	POSTGRES = "postgres"
	MYSQL    = "mysql"
	MONGODB  = "mongoDB"
)

// Features
const (
	INIT         = "init"
	CLOUD_NATIVE = "cloud native"
	AWS          = "AWS"
	CREATE_CD    = "Create CD pipeline"
	CREATE_INFRA = "Create Infra"
)
