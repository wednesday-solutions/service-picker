package constants

var Repos = func() map[string]string {
	return map[string]string{
		"ReactJS":                     "https://github.com/wednesday-solutions/react-template",
		"NextJS":                      "https://github.com/wednesday-solutions/nextjs-template",
		"Node (Hapi- REST API)":       "https://github.com/wednesday-solutions/nodejs-hapi-template",
		"Node (Express- GraphQL API)": "https://github.com/wednesday-solutions/node-express-graphql-template",
		"Node (Express- REST API)":    "https://github.com/wednesday-solutions/node-mongo-express",
		"GolangPostgreSQL":            "https://github.com/wednesday-solutions/go-template",
		"GolangMySQL":                 "https://github.com/wednesday-solutions/go-template-mysql",
		"React Native":                "https://github.com/wednesday-solutions/react-native-template",
		"Android":                     "https://github.com/wednesday-solutions/android-template",
		"IOS":                         "https://github.com/wednesday-solutions/ios-template",
		"Flutter":                     "https://github.com/wednesday-solutions/flutter_template",
	}
}

// CLI options
const (
	Picky   = "picky"
	Service = "service"
	Test    = "test"
)

// Services
const (
	Web     = "web"
	Mobile  = "mobile"
	Backend = "backend"
)

// Frontend stacks
const (
	ReactJS = "ReactJS"
	NextJS  = "NextJS"
)

// Backend stacks
const (
	NodeHapiTemplate           = "Node (Hapi- REST API)"
	NodeExpressGraphqlTemplate = "Node (Express- GraphQL API)"
	NodeExpressTemplate        = "Node (Express- REST API)"
	GolangEchoTemplate         = "Golang (Echo- GraphQL API)"
)

// Mobile stacks
const (
	ReactNativeTemplate = "React Native"
	AndroidTemplate     = "Android"
	IOSTemplate         = "IOS"
	FlutterTemplate     = "Flutter"
)

// Databases
const (
	PostgreSQL = "PostgreSQL"
	MySQL      = "MySQL"
	MongoDB    = "MongoDB"
)

// Features
const (
	Init        = "init"
	CloudNative = "cloud native"
	AWS         = "AWS"
	CreateCD    = "create CD pipeline"
	CreateInfra = "create Infra"
)

// Github Repo download URL
const (
	GitHubBaseURL                  = "https://raw.githubusercontent.com/wednesday-solutions/"
	CDFilePathURL                  = "/.github/workflows/cd.yml"
	NodeHapiTemplateRepo           = "nodejs-hapi-template/main"
	NodeExpressGraphqlTemplateRepo = "node-express-graphql-template/develop"
	NodeExpressTemplateRepo        = "node-mongo-express/main"
	GoEchoTemplatePostgresRepo     = "go-template/master"
	GoEchoTemplateMysqlRepo        = "go-template-mysql/main"
	ReactTemplateRepo              = "react-template/master"
	NextjsTemplateRepo             = "nextjs-template/master"
)

// Infra Files
const (
	PackageDotJsonFile = "package.json"
	EnvFile            = ".env"
	SstConfigJsFile    = "sst.config.js"
	WebStackJsFile     = "WebStack.js"
)

// helpers
const (
	Mysql                    = "mysql"
	Redis                    = "redis"
	Postgres                 = "postgres"
	Frontend                 = "frontend"
	Database                 = "database"
	GolangPostgreSQLTemplate = "GolangPostgreSQL"
	GolangMySQLTemplate      = "GolangMySQL"
)

// Files
const (
	DockerComposeFile = "docker-compose.yml"
	DockerFile        = "Dockerfile"
	DockerEnvFile     = ".env.docker"
	DockerIgnoreFile  = ".dockerignore"
)
