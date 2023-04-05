package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/fatih/color"
	"github.com/iancoleman/strcase"
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
)

// DirectoryName will create a default suffix for the stack which is selected by the user.
// The directoryName depends on user input, stack, and the database which are selected or provided by the user.
func DirectoryName(dirName, stack, database string) string {
	switch stack {
	case constants.NodeHapiTemplate:
		if database == constants.PostgreSQL {
			dirName = fmt.Sprintf("%s-%s", dirName, constants.NodeHapiPgTemplate)
		} else if database == constants.MySQL {
			dirName = fmt.Sprintf("%s-%s", dirName, constants.NodeHapiMySqlTemplate)
		}
	case constants.NodeExpressGraphqlTemplate:
		if database == constants.PostgreSQL {
			dirName = fmt.Sprintf("%s-%s", dirName, constants.NodeGraphqlPgTemplate)
		} else if database == constants.MySQL {
			dirName = fmt.Sprintf("%s-%s", dirName, constants.NodeGraphqlMySqlTemplate)
		}
	case constants.NodeExpressTemplate:
		if database == constants.MongoDB {
			dirName = fmt.Sprintf("%s-%s", dirName, constants.NodeExpressMongoTemplate)
		}
	case constants.GolangEchoTemplate:
		if database == constants.PostgreSQL {
			dirName = fmt.Sprintf("%s-%s", dirName, constants.GolangPgTemplate)
		} else if database == constants.MySQL {
			dirName = fmt.Sprintf("%s-%s", dirName, constants.GolangMySqlTemplate)
		}
	case constants.ReactJS:
		dirName = fmt.Sprintf("%s-%s", dirName, constants.ReactTemplate)
	case constants.NextJS:
		dirName = fmt.Sprintf("%s-%s", dirName, constants.NextTemplate)
	case constants.ReactNative:
		dirName = fmt.Sprintf("%s-%s", dirName, constants.ReactNativeTemplate)
	case constants.Android:
		dirName = fmt.Sprintf("%s-%s", dirName, constants.AndroidTemplate)
	case constants.IOS:
		dirName = fmt.Sprintf("%s-%s", dirName, constants.IOSTemplate)
	case constants.Flutter:
		dirName = fmt.Sprintf("%s-%s", dirName, constants.FlutterTemplate)
	}
	return dirName
}

// ExistingStacksDatabasesAndDirectories return existing stacks details.
// Stack details contain stack name, database, and the directory name.
func ExistingStacksDatabasesAndDirectories() ([]string, []string, []string) {
	var stacks, databases, dirNames []string
	var stack, database string
	directories, err := ReadAllContents(CurrentDirectory())
	errorhandler.CheckNilErr(err)

	for _, dirName := range directories {
		stack, database = FindStackAndDatabase(dirName)
		if stack != "" {
			stacks = append(stacks, stack)
			databases = append(databases, database)
			dirNames = append(dirNames, dirName)
		}
	}
	return stacks, databases, dirNames
}

// FindStackAndDatabase return stack and database of given directory name.
func FindStackAndDatabase(dirName string) (string, string) {
	var stack, database string
	_, stackSuffix, lastSuffix := SplitStackDirectoryName(dirName)

	switch lastSuffix {
	case constants.Pg:
		database = constants.PostgreSQL
		if stackSuffix == "hapi" {
			stack = constants.NodeHapiTemplate
		} else if stackSuffix == "graphql" {
			stack = constants.NodeExpressGraphqlTemplate
		} else if stackSuffix == "golang" {
			stack = constants.GolangEchoTemplate
		}
	case constants.Mysql:
		database = constants.MySQL
		if stackSuffix == "hapi" {
			stack = constants.NodeHapiTemplate
		} else if stackSuffix == "graphql" {
			stack = constants.NodeExpressGraphqlTemplate
		} else if stackSuffix == "golang" {
			stack = constants.GolangEchoTemplate
		}
	case constants.Mongo:
		database = constants.MongoDB
		if stackSuffix == "express" {
			stack = constants.NodeExpressTemplate
		}
	case constants.Web:
		if stackSuffix == "react" {
			stack = constants.ReactJS
		} else if stackSuffix == "next" {
			stack = constants.NextJS
		}
	case constants.Mobile:
		if stackSuffix == "reactnative" {
			stack = constants.ReactNative
		} else if stackSuffix == "android" {
			stack = constants.Android
		} else if stackSuffix == "ios" {
			stack = constants.IOS
		} else if stackSuffix == "flutter" {
			stack = constants.Flutter
		}
	}
	return stack, database
}

// SplitStackDirectoryName returns user-input, stack-suffix and last-suffix of the given stack directory name.
func SplitStackDirectoryName(dirName string) (string, string, string) {
	var userInput, stackSuffix, lastSuffix string
	var splitDirName []string
	var isBackendStack bool
	splitDirName = strings.Split(dirName, "-")
	if len(splitDirName) > 2 {
		lastSuffix = splitDirName[len(splitDirName)-1]
		stackSuffix = splitDirName[len(splitDirName)-2]
		if lastSuffix == constants.Pg || lastSuffix == constants.Mysql || lastSuffix == constants.Mongo {
			isBackendStack = true
		}
		var suffixSize int
		if isBackendStack {
			if len(splitDirName) > 3 {
				suffixSize = 3
			}
		} else {
			suffixSize = 2
		}
		userInput = splitDirName[0]
		for _, split := range splitDirName[1 : len(splitDirName)-suffixSize] {
			userInput = fmt.Sprintf("%s_%s", userInput, split)
		}
	}
	return userInput, stackSuffix, lastSuffix
}

// FindUserInputStackName return user-input of the given stack directory name.
func FindUserInputStackName(dirName string) string {
	userInput, _, _ := SplitStackDirectoryName(dirName)
	return userInput
}

// ExistingStackAndDatabase return stack and database of the given stack directory name.
func ExistingStackAndDatabase(dirName string) (string, string) {
	stack, database := FindStackAndDatabase(dirName)
	return stack, database
}

// FindService return service of the given stack directory name.
func FindService(dirName string) string {
	_, _, lastSuffix := SplitStackDirectoryName(dirName)
	switch lastSuffix {
	case constants.Pg, constants.Mysql, constants.Mongo:
		return constants.Backend
	default:
		return lastSuffix
	}
}

// CovertToCamelCase return camel cased array of string of the given array of string.
func ConvertToCamelCase(slice []string) []string {
	camelSlice := []string{}
	for _, str := range slice {
		camelSlice = append(camelSlice, strcase.ToCamel(str))
	}
	return camelSlice
}

// CreateMessageTemplate creates new text template for printing colorful logs.
func CreateMessageTemplate(name, text string) *template.Template {
	tpl, err := template.New(name).Parse(text)
	errorhandler.CheckNilErr(err)
	tpl = template.Must(tpl, err)
	return tpl
}

// PrintMultiSelectMessage prints multi selected options.
func PrintMultiSelectMessage(messages []string) error {
	var message, coloredMessage string
	var tpl *template.Template
	if len(messages) > 0 {
		var templateText string
		if len(messages) == 1 {
			templateText = fmt.Sprintf("%s %d option selected: {{ . }}\n",
				constants.IconSelect,
				len(messages))
		} else {
			templateText = fmt.Sprintf("%s %d options selected: {{ . }}\n",
				constants.IconSelect,
				len(messages))
		}
		for _, option := range messages {
			message = fmt.Sprintf("%s%s ", message, option)
		}
		coloredMessage = color.GreenString("%s", message)
		tpl = CreateMessageTemplate("message", templateText)
	} else {
		message = "No options selected"
		coloredMessage = color.YellowString("%s", message)
		tpl = CreateMessageTemplate("responseMessage", fmt.Sprintf("%s {{ . }}\n", constants.IconWarn))
	}
	err := tpl.Execute(os.Stdout, coloredMessage)
	return err
}

// PrintWarningMessage prints given message in yellow color as warning message in terminal.
func PrintWarningMessage(message string) error {
	tpl := CreateMessageTemplate("warningMessage", fmt.Sprintf("\n%s {{ . }}\n", constants.IconWarn))
	message = color.YellowString("%s", message)
	err := tpl.Execute(os.Stdout, message)
	return err
}

// PrintInfoMessage prints given message in cyan color as info message in terminal.
func PrintInfoMessage(message string) error {
	tpl := CreateMessageTemplate("InfoMessage", fmt.Sprintf("\n%s {{ . }}\n", constants.IconChoose))
	message = color.CyanString("%s", message)
	err := tpl.Execute(os.Stdout, message)
	return err
}

// GetSuffixOfStack returns suffix name for the given stack and database.
func GetSuffixOfStack(stack, database string) string {
	var suffix string
	switch stack {
	case constants.ReactJS:
		suffix = constants.ReactTemplate
	case constants.NextJS:
		suffix = constants.NextTemplate
	case constants.NodeHapiTemplate:
		if database == constants.PostgreSQL {
			suffix = constants.NodeHapiPgTemplate
		} else if database == constants.MySQL {
			suffix = constants.NodeHapiMySqlTemplate
		}
	case constants.NodeExpressGraphqlTemplate:
		if database == constants.PostgreSQL {
			suffix = constants.NodeGraphqlPgTemplate
		} else if database == constants.MySQL {
			suffix = constants.NodeGraphqlMySqlTemplate
		}
	case constants.NodeExpressTemplate:
		if database == constants.MongoDB {
			suffix = constants.NodeExpressMongoTemplate
		}
	case constants.GolangEchoTemplate:
		if database == constants.PostgreSQL {
			suffix = constants.GolangPgTemplate
		} else if database == constants.MySQL {
			suffix = constants.GolangMySqlTemplate
		}
	case constants.ReactNative:
		suffix = constants.ReactNativeTemplate
	case constants.Android:
		suffix = constants.AndroidTemplate
	case constants.IOS:
		suffix = constants.IOSTemplate
	case constants.Flutter:
		suffix = constants.FlutterTemplate
	}
	return suffix
}

type StackDetails struct {
	Name      string
	Language  string
	Framework string
	Type      string
	Databases string
}

// GetStackDetails returns an array of StackDetails for showing details
// when user selects stacks prompt.
func GetStackDetails(service string) []StackDetails {
	var stacksDetails []StackDetails
	switch service {
	case constants.Backend:
		stacksDetails = []StackDetails{
			{
				Name:      constants.NodeHapiTemplate,
				Language:  "JavaScript",
				Framework: "Node JS & Hapi",
				Type:      "REST API",
				Databases: fmt.Sprintf("%s %s", constants.PostgreSQL, constants.MySQL),
			},
			{
				Name:      constants.NodeExpressGraphqlTemplate,
				Language:  "JavaScript",
				Framework: "Node JS & Express",
				Type:      "GraphQL API",
				Databases: fmt.Sprintf("%s %s", constants.PostgreSQL, constants.MySQL),
			},
			{
				Name:      constants.NodeExpressTemplate,
				Language:  "JavaScript",
				Framework: "Node JS & Express",
				Type:      "REST API",
				Databases: constants.MongoDB,
			},
			{
				Name:      constants.GolangEchoTemplate,
				Language:  "Golang",
				Framework: "Echo",
				Type:      "GraphQL API",
				Databases: fmt.Sprintf("%s %s", constants.PostgreSQL, constants.MySQL),
			},
		}
	case constants.Web:
		stacksDetails = []StackDetails{
			{
				Name:      constants.ReactJS,
				Language:  "JavaScript",
				Framework: "React",
			},
			{
				Name:      constants.NextJS,
				Language:  "JavaScript",
				Framework: "Next.js",
			},
		}
	case constants.Mobile:
		stacksDetails = []StackDetails{
			{
				Name:      constants.ReactNative,
				Language:  "JavaScript",
				Framework: "React Native",
			},
			{
				Name:      constants.Android,
				Language:  "Kotlin",
				Framework: "-",
			},
			{
				Name:      constants.IOS,
				Language:  "Swift",
				Framework: "-",
			},
			{
				Name:      constants.Flutter,
				Language:  "Dart",
				Framework: "Flutter",
			},
		}
	}
	return stacksDetails
}

// FindConfigStacks will return an array of existing stack functions in sst.config.js
// Eg: [ApiNodeHapiMysql, FeReactWeb]
func FindConfigStacks(configLine string) []string {
	var stack string
	var stacks []string
	stackFound := false

	for _, char := range configLine {
		if char == '(' {
			stackFound = true
			continue
		} else if char == ')' {
			stackFound = false
			stacks = append(stacks, stack)
			stack = ""
		}
		if stackFound {
			stack = fmt.Sprintf("%s%s", stack, string(char))
		}
	}
	return stacks
}

// FindStackDirectoriesByConfigStacks will return an array of stack directories present
// in the root directory. Eg: [api-node-hapi-mysql, fe-react-web]
func FindStackDirectoriesByConfigStacks(configStacks []string) []string {
	var stacks []string

	_, _, directories := ExistingStacksDatabasesAndDirectories()
	camelCaseDirectories := ConvertToCamelCase(directories)

	for _, configStack := range configStacks {
		for idx, camelCaseDirName := range camelCaseDirectories {
			if configStack == camelCaseDirName {
				stacks = append(stacks, directories[idx])
			}
		}
	}
	return stacks
}

// RunCommandWithLogs runs the given command with logs.
func RunCommandWithLogs(path string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if path != "" {
		cmd.Dir = path
	}
	err := cmd.Run()
	fmt.Printf("\n")
	return err
}

// RunCommandWithLogs runs the given command without logs.
func RunCommandWithoutLogs(path string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if path != "" {
		cmd.Dir = path
	}
	err := cmd.Run()
	return err
}

// IsYarnOrNpmInstalled checks whether yarn or npm is installed in the user's machine.
// If both are not installed, then the system will throw error.
func IsYarnOrNpmInstalled() string {
	var pkgManager string
	err := RunCommandWithoutLogs("", constants.Yarn, "-v")
	if err != nil {
		err = RunCommandWithoutLogs("", constants.Npm, "-v")
		if err != nil {
			// Throw error either yarn or npm not installed
			errorhandler.CheckNilErr(fmt.Errorf("Please install 'yarn' or 'npm' in your machine.\n"))
		} else {
			// install 'yarn' if 'npm' installed already.
			err = RunCommandWithoutLogs("", constants.Npm, "install", "--global", constants.Yarn)
			errorhandler.CheckNilErr(err)
			pkgManager = constants.Yarn
		}
	} else {
		pkgManager = constants.Yarn
	}
	return pkgManager
}

// GetShortEnvironment return short environment name for the given environment.
func GetShortEnvironment(environment string) string {
	switch environment {
	case constants.Development:
		return constants.Dev
	case constants.QA:
		return constants.QA
	case constants.Production:
		return constants.Prod
	default:
		return environment
	}
}
