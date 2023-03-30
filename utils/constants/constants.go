package constants

var Repos = func() map[string]string {
	return map[string]string{
		"ReactJS":                     "https://github.com/wednesday-solutions/react-template",
		"NextJS":                      "https://github.com/wednesday-solutions/nextjs-template",
		"Node (Hapi- REST API)":       "https://github.com/wednesday-solutions/nodejs-hapi-template",
		"Node (Express- GraphQL API)": "https://github.com/wednesday-solutions/node-express-graphql-template",
		"Node (Express- REST API)":    "https://github.com/wednesday-solutions/node-mongo-express",
		"Golang-PostgreSQL":           "https://github.com/wednesday-solutions/go-template",
		"Golang-MySQL":                "https://github.com/wednesday-solutions/go-template-mysql",
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
	ReactNative = "React Native"
	Android     = "Android"
	IOS         = "IOS"
	Flutter     = "Flutter"
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
	SstConfigJsFile    = "sst.config.js"
	WebStackJsFile     = "WebStack.js"
	BackendStackJsFile = "BackendStack.js"
	EnvFile            = ".env"
	EnvDevFile         = ".env.development"
)

// helpers
const (
	Yes                      = "Yes"
	No                       = "No"
	Stack                    = "stack"
	Stacks                   = "stacks"
	Mysql                    = "mysql"
	Redis                    = "redis"
	Postgres                 = "postgres"
	Frontend                 = "frontend"
	Database                 = "database"
	ProjectName              = "projectName"
	WebStatus                = "webStatus"
	MobileStatus             = "mobileStatus"
	BackendStatus            = "backendStatus"
	GolangPostgreSQLTemplate = "Golang-PostgreSQL"
	GolangMySQLTemplate      = "Golang-MySQL"
	WebDirName               = "webDirName"
	MobileDirName            = "mobileDirName"
	BackendDirName           = "backendDirName"
	SizeOfPromptSelect       = 8
)

// Docker related files
const (
	DockerComposeFile = "docker-compose.yml"
	DockerFile        = "Dockerfile"
	DockerEnvFile     = ".env.docker"
	DockerIgnoreFile  = ".dockerignore"
)

// Template directory name
const (
	NodeHapiPgTemplate       = "node-hapi-pg"
	NodeHapiMySqlTemplate    = "node-hapi-mysql"
	NodeGraphqlPgTemplate    = "node-graphql-pg"
	NodeGraphqlMySqlTemplate = "node-graphql-mysql"
	NodeExpressMongoTemplate = "node-express-mongo"
	GolangPgTemplate         = "golang-pg"
	GolangMySqlTemplate      = "golang-mysql"
	ReactTemplate            = "react-web"
	NextTemplate             = "next-web"
	ReactNativeTemplate      = "reactnative-mobile"
	AndroidTemplate          = "android-mobile"
	IOSTemplate              = "ios-mobile"
	FlutterTemplate          = "flutter-mobile"
)
