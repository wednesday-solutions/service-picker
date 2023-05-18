package constants

import "github.com/spaceweasel/promptui"

var Repos = func() map[string]string {
	return map[string]string{
		ReactJS:                    "https://github.com/wednesday-solutions/react-template",
		NextJS:                     "https://github.com/wednesday-solutions/nextjs-template",
		ReactGraphqlTS:             "https://github.com/wednesday-solutions/react-graphql-ts-template",
		NodeHapiTemplate:           "https://github.com/wednesday-solutions/nodejs-hapi-template",
		NodeExpressGraphqlTemplate: "https://github.com/wednesday-solutions/node-express-graphql-template",
		NodeExpressTemplate:        "https://github.com/wednesday-solutions/node-mongo-express",
		GolangPostgreSQLTemplate:   "https://github.com/wednesday-solutions/go-template",
		GolangMySQLTemplate:        "https://github.com/wednesday-solutions/go-template-mysql",
		ReactNative:                "https://github.com/wednesday-solutions/react-native-template",
		Android:                    "https://github.com/wednesday-solutions/android-template",
		IOS:                        "https://github.com/wednesday-solutions/ios-template",
		Flutter:                    "https://github.com/wednesday-solutions/flutter_template",
	}
}

// CLI options
const (
	Picky   = "picky"
	Service = "service"
	Test    = "test"
	Init    = "init"
	Create  = "create"
	Infra   = "infra"
)

// Home options
const (
	InitService   = "Init Service"
	DockerCompose = "Docker Compose"
	CICD          = "CI/CD"
	SetupInfra    = "Setup Infra"
	Deploy        = "Deploy"
	RemoveDeploy  = "Remove Deploy"
	GitInit       = "Git Init"
	Exit          = "Exit"
)

// Services
const (
	Web     = "web"
	Mobile  = "mobile"
	Backend = "backend"
)

// Frontend stacks
const (
	ReactJS        = "React JS"
	NextJS         = "Next JS"
	ReactGraphqlTS = "React GraphQL TS"
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

// Docker compose actions
const (
	CreateDockerCompose = "Create Docker Compose"
	RunDockerCompose    = "Run Docker Compose"
)

// CICD actions
const (
	CreateCI = "Create CI"
	CreateCD = "Create CD"
)

// Infra Files
const (
	PackageDotJsonFile = "package.json"
	SstConfigFile      = "sst.config.js"
	EnvFile            = ".env"
	EnvDevFile         = ".env.development"
	ParseSstOutputs    = "parseSstOutputs.js"
)

// Cloud Providers
const (
	AWS = "AWS"
)

// helpers
const (
	Yes                      = "Yes"
	No                       = "No"
	Stack                    = "stack"
	Stacks                   = "stacks"
	DB                       = "db"
	Pg                       = "pg"
	Mysql                    = "mysql"
	Postgresql               = "postgresql"
	Mongodb                  = "mongodb"
	Mongo                    = "mongo"
	Redis                    = "redis"
	Postgres                 = "postgres"
	Frontend                 = "frontend"
	Database                 = "database"
	ProjectName              = "projectName"
	PgNative                 = "pg-native"
	Mysql2                   = "mysql2"
	WebStatus                = "webStatus"
	MobileStatus             = "mobileStatus"
	BackendStatus            = "backendStatus"
	GolangPostgreSQLTemplate = "Golang-PostgreSQL"
	GolangMySQLTemplate      = "Golang-MySQL"
	WebDirName               = "webDirName"
	MobileDirName            = "mobileDirName"
	BackendDirName           = "backendDirName"
	SizeOfPromptSelect       = 8
	All                      = "All"
	SstConfigStack           = "sstConfigStack"
	ExistingDirectories      = "existingDirectories"
	Yarn                     = "yarn"
	Npm                      = "npm"
	WebDirectories           = "webDirectories"
	BackendPgDirectories     = "backendPgDirectories"
	BackendMysqlDirectories  = "backendMysqlDirectories"
	Zero                     = 0
	One                      = 1
	Two                      = 2
	Three                    = 3
	BackendSuffixSize        = 3
	WebSuffixSize            = 3
	MobileSuffixSize         = 2
	DotSstDirectory          = ".sst"
	PostgresqlData           = "postgresql/data"
	GithubWorkflowsDir       = ".github/workflows"
	DotGitFolder             = ".git"
	CDFilePathURL            = "/.github/workflows/cd.yml"
	GitHub                   = "GitHub"
	Graphql                  = "graphql"
	Git                      = "git"
	DotGitIgnore             = ".gitignore"
	NodeModules              = "node_modules"
	CloudProvider            = "cloudprovider"
	Directory                = "directory"
	DockerComposeFlag        = "dockercompose"
	CIFlag                   = "ci"
	CDFlag                   = "cd"
	Platform                 = "platform"
	Github                   = "github"
	DotSst                   = ".sst"
	OutputsJson              = "outputs.json"
	DBUsername               = "username"
	DeployNow                = "Deploy Now"
	ChangeStacks             = "Change Stacks"
	GoBack                   = "Go Back"
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
	GolangPgTemplate         = "golang-graphql-pg"
	GolangMySqlTemplate      = "golang-graphql-mysql"
	ReactTemplate            = "react-js-web"
	NextTemplate             = "next-js-web"
	ReactGraphqlTemplate     = "react-graphql-web"
	ReactNativeTemplate      = "reactnative-mobile"
	AndroidTemplate          = "android-mobile"
	IOSTemplate              = "ios-mobile"
	FlutterTemplate          = "flutter-mobile"
)

// Environments
const (
	Environment = "environment"
	Dev         = "dev"
	QA          = "qa"
	Prod        = "prod"
	Development = "development"
	Production  = "production"
	Develop     = "develop"
)

// UI icons
var (
	IconChoose   = promptui.Styler(promptui.FGBold)("▸")
	IconSelect   = promptui.Styler(promptui.FGGreen)("✔")
	IconWarn     = promptui.Styler(promptui.FGYellow)("⚠")
	IconWrong    = promptui.Styler(promptui.FGRed)("✗")
	IconQuestion = promptui.Styler(promptui.FGMagenta)("?")
)

// Stack short name for flags
const (
	ReactjsLower      = "reactjs"
	NextjsLower       = "nextjs"
	ReactGraphqlLower = "reactgraphql"
	ReactNativeLower  = "reactnative"
	AndroidLower      = "android"
	IOSLower          = "ios"
	FlutterLower      = "flutter"
	NodeHapi          = "nodehapi"
	NodeGraphql       = "nodegraphql"
	NodeExpress       = "nodeexpress"
	Golang            = "golang"
)

// GitHub secrets
const (
	DistributionId     = "DISTRIBUTION_ID"
	AwsRegion          = "AWS_REGION"
	AwsAccessKeyId     = "AWS_ACCESS_KEY_ID"
	AwsSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
	AwsEcrRepository   = "AWS_ECR_REPOSITORY"
)

// Port numbers for infra setup
const (
	WebPortNumber      = 3000
	BackendPortNumber  = 9000
	PostgresPortNumber = 5432
	MysqlPortNumber    = 3306
	RedisPortNumber    = 6379
)

// Env values
const (
	PostgresUser  = "POSTGRES_USER"
	PostgresHost  = "POSTGRES_HOST"
	PostgresDB    = "POSTGRES_DB"
	MysqlUser     = "MYSQL_USER"
	MysqlHost     = "MYSQL_HOST"
	MysqlDatabase = "MYSQL_DATABASE"
	RedisHost     = "REDIS_HOST"
	RedisPort     = "REDIS_PORT"
)

// Port number names
const (
	BackendPort  = "PORT"
	PostgresPort = "POSTGRES_PORT"
	MysqlPort    = "MYSQL_PORT"
)
