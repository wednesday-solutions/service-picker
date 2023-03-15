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
	NODE_HAPI_TEMPLATE            = "Node (Hapi- REST API)"
	NODE_EXPRESS_GRAPHQL_TEMPLATE = "Node (Express- GraphQL API)"
	NODE_EXPRESS_TS               = "Node (Express- TypeScript)"
	GOLANG_ECHO_TEMPLATE          = "Golang (Echo- GraphQL API)"
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

// Github download URL
const (
	GitHubBaseURL                  = "https://raw.githubusercontent.com/wednesday-solutions/"
	CDFilePathURL                  = "/.github/workflows/cd.yml"
	NodeHapiTemplateRepo           = "nodejs-hapi-template/main"
	NodeExpressGraphqlTemplateRepo = "node-express-graphql-template/develop"
	GoEchoTemplatePostgresRepo     = "go-template/master"
	GoEchoTemplateMysqlRepo        = "go-template-mysql/main"
	ReactTemplateRepo              = "react-template/master"
	NextjsTemplateRepo             = "nextjs-template/master"
)
