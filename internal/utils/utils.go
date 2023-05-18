package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
)

// SplitStackDirectoryName returns user-input, stack-suffix and last-suffix of the given stack directory name.
func SplitStackDirectoryName(dirName string) (string, string, string, string) {
	var userInput, langSuffix, stackSuffix, lastSuffix string
	var splitDirName []string
	var isBackendStack, isWebStack, isMobileStack bool
	splitDirName = strings.Split(dirName, "-")
	if len(splitDirName) > constants.Two {
		lastSuffix = splitDirName[len(splitDirName)-constants.One]
		stackSuffix = splitDirName[len(splitDirName)-constants.Two]
		langSuffix = splitDirName[len(splitDirName)-constants.Three]
		if lastSuffix == constants.Pg || lastSuffix == constants.Mysql || lastSuffix == constants.Mongo {
			isBackendStack = true
		} else if lastSuffix == constants.Web {
			isWebStack = true
		} else if lastSuffix == constants.Mobile {
			isMobileStack = true
		}
		var suffixSize int
		if isBackendStack {
			suffixSize = constants.BackendSuffixSize
		} else if isWebStack {
			suffixSize = constants.WebSuffixSize
		} else if isMobileStack {
			suffixSize = constants.MobileSuffixSize
		}
		userInput = splitDirName[constants.Zero]
		for _, split := range splitDirName[constants.One : len(splitDirName)-suffixSize] {
			userInput = fmt.Sprintf("%s-%s", userInput, split)
		}
	}
	return userInput, langSuffix, stackSuffix, lastSuffix
}

// CovertToCamelCase return camel cased array of string of the given array of string.
func ConvertToCamelCase(slice []string) []string {
	camelSlice := []string{}
	for _, str := range slice {
		camelSlice = append(camelSlice, strcase.ToCamel(str))
	}
	return camelSlice
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

// GetPackageManagerOfUser checks whether yarn or npm is installed in the user's machine.
// If both are not installed, then the system will throw error.
func GetPackageManagerOfUser() string {
	var pkgManager string
	err := RunCommandWithoutLogs("", constants.Yarn, "-v")
	if err != nil {
		err = RunCommandWithoutLogs("", constants.Npm, "-v")
		if err != nil {
			// Throw error either yarn or npm not installed
			errorhandler.CheckNilErr(fmt.Errorf("Please install 'yarn' or 'npm' in your machine.\n"))
		} else {
			pkgManager = constants.Npm
		}
	} else {
		pkgManager = constants.Yarn
	}
	return pkgManager
}

// GetProjectName returns projectName
func GetProjectName() string {
	currentDir := CurrentDirectory()
	splitDirs := strings.Split(currentDir, "/")
	projectName := splitDirs[len(splitDirs)-1]
	return projectName
}

func CreateGithubWorkflowDir() {
	currentDir := CurrentDirectory()
	workflowsPath := fmt.Sprintf("%s/%s", currentDir,
		constants.GithubWorkflowsDir,
	)
	workflowStatus, _ := IsExists(workflowsPath)
	if !workflowStatus {
		githubFolderPath := fmt.Sprintf("%s/%s", currentDir, ".github")
		githubStatus, _ := IsExists(githubFolderPath)
		if !githubStatus {
			err := CreateDirectory(githubFolderPath)
			errorhandler.CheckNilErr(err)
		}
		err := CreateDirectory(workflowsPath)
		errorhandler.CheckNilErr(err)
	}
}

func EndsWith(inputString, endString string) bool {
	if len(inputString) < len(endString) {
		return false
	} else if len(inputString) == len(endString) {
		return inputString == endString
	} else {
		for i, j := len(endString)-1, len(inputString)-1; i >= 0; i, j = i-1, j-1 {
			if endString[i] != inputString[j] {
				return false
			}
		}
		return true
	}
}

func StartsWith(inputString, startString string) bool {
	if len(inputString) < len(startString) {
		return false
	} else if len(inputString) == len(startString) {
		return inputString == startString
	} else {
		for i, j := 0, 0; i < len(startString); i, j = i+1, j+1 {
			if startString[i] != inputString[j] {
				return false
			}
		}
		return true
	}
}

type LogDriverOptionsKeys struct {
	AwsLogsGroup        string `json:"awsLogsGroup"`
	AwsLogsStreamPrefix string `json:"awsLogsStreamPrefix"`
	AwsLogsRegion       string `json:"awsLogsRegion"`
}

type WebOutputKeys struct {
	DistributionId string `json:"distributionId"`
	BucketName     string `json:"bucketName"`
	SiteUrl        string `json:"siteUrl"`
}

type BackendOutputKeys struct {
	TaskRole                     string `json:"taskRole"`
	Image                        string `json:"image"`
	ContainerName                string `json:"containerName"`
	ContainerPort                string `json:"containerPort"`
	ExecutionRole                string `json:"executionRole"`
	TaskDefinition               string `json:"taskDefinition"`
	LogDriver                    string `json:"logDriver"`
	LogDriverOptions             LogDriverOptionsKeys
	Family                       string `json:"family"`
	AwsRegion                    string `json:"awsRegion"`
	RedisHost                    string `json:"redisHost"`
	SecretName                   string `json:"secretName"`
	DatabaseHost                 string `json:"databaseHost"`
	DatabaseName                 string `json:"databaseName"`
	SecretArn                    string `json:"secretArn"`
	LoadBalancerDns              string `json:"loadBalancerDns"`
	ServiceName                  string `json:"serviceName"`
	ClusterName                  string `json:"clusterName"`
	ElasticContainerRegistryRepo string `json:"elasticContainerRegistryRepo"`
}

type TaskDefinitionDetails struct {
	BackendObj  BackendOutputKeys
	Environment string
	EnvName     string
	SecretName  string
}

func ReadJsonDataInSstOutputs() map[string]interface{} {
	file := fmt.Sprintf(
		"%s/%s/%s", CurrentDirectory(), constants.DotSst, constants.OutputsJson,
	)
	status, _ := IsExists(file)
	if !status {
		return nil
	}
	sstOutputFile, err := os.Open(file)
	errorhandler.CheckNilErr(err)

	fileContent, err := ioutil.ReadAll(sstOutputFile)
	errorhandler.CheckNilErr(err)

	var data interface{}
	err = json.Unmarshal(fileContent, &data)
	errorhandler.CheckNilErr(err)

	if jsonData, ok := data.(map[string]interface{}); ok {
		return jsonData
	}
	return nil
}

func GetOutputsBackendObject(environment, stackDir string) TaskDefinitionDetails {
	camelCaseDir := strcase.ToCamel(stackDir)
	jsonData := ReadJsonDataInSstOutputs()
	if jsonData == nil {
		errorhandler.CheckNilErr(fmt.Errorf("outputs.json is not valid."))
	}
	var td TaskDefinitionDetails
	td.Environment = environment
	td.EnvName = GetShortEnvName(environment)
	key := fmt.Sprintf("%s-web-app-%s", td.EnvName, camelCaseDir)

	if value, ok := jsonData[key]; ok {

		jsonOutput, err := json.Marshal(value)
		errorhandler.CheckNilErr(err)

		var backendObj BackendOutputKeys
		err = json.Unmarshal(jsonOutput, &backendObj)
		errorhandler.CheckNilErr(err)

		td.BackendObj = backendObj
		return td
	}
	return TaskDefinitionDetails{}
}

func CreateInfraOutputsJson(environment string) error {

	jsonData := ReadJsonDataInSstOutputs()
	if jsonData == nil {
		errorhandler.CheckNilErr(fmt.Errorf("outputs.json is not valid"))
	}
	envName := GetShortEnvName(environment)
	_, _, directories := GetExistingStacksDatabasesAndDirectories()
	for _, stackDir := range directories {
		camelCaseDir := strcase.ToCamel(stackDir)

		key := fmt.Sprintf("%s-web-app-%s", envName, camelCaseDir)

		if value, ok := jsonData[key]; ok {
			outputFile := fmt.Sprintf("%s/%s-%s", CurrentDirectory(), stackDir, constants.OutputsJson)

			status, _ := IsExists(outputFile)
			if !status {
				err := CreateFile(outputFile)
				errorhandler.CheckNilErr(err)
			}

			jsonOutput, err := json.MarshalIndent(value, "", "\t")
			errorhandler.CheckNilErr(err)

			err = WriteToFile(outputFile, string(jsonOutput))
			errorhandler.CheckNilErr(err)
		}
	}
	return nil
}

// GetShortEnvName return short environment name for the given environment.
func GetShortEnvName(environment string) string {
	var shortEnv string
	if environment == constants.Dev || environment == constants.Develop || environment == constants.Development {
		shortEnv = constants.Dev
	} else if environment == constants.QA {
		shortEnv = constants.QA
	} else if environment == constants.Prod || environment == constants.Production {
		shortEnv = constants.Prod
	}
	return shortEnv
}

func GetPortNumber(defaultPN int) int {
	_, _, stackDirs := GetExistingStacksDatabasesAndDirectories()
	var backendServices []string
	for _, dir := range stackDirs {
		service := FindService(dir)
		if service == constants.Backend {
			backendServices = append(backendServices, service)
		}
	}
	currentPN := defaultPN + len(backendServices) - 1
	return currentPN
}

func GetDatabasePortNumber(driver string) int {
	_, _, stackDirs := GetExistingStacksDatabasesAndDirectories()
	var postgresStacks, mysqlStacks []string
	for _, dir := range stackDirs {
		service := FindService(dir)
		if service == constants.Backend {
			_, database := FindStackAndDatabase(dir)
			if database == constants.PostgreSQL {
				postgresStacks = append(postgresStacks, dir)
			} else if database == constants.MySQL {
				mysqlStacks = append(mysqlStacks, dir)
			}
		}
	}
	var currentPortNumber int
	if driver == constants.PostgreSQL {
		currentPortNumber = constants.PostgresPortNumber + len(postgresStacks) - 1
	} else if driver == constants.MySQL {
		currentPortNumber = constants.MysqlPortNumber + len(mysqlStacks) - 1
	}
	return currentPortNumber
}

func FetchExistingPortNumber(stackDir, portName string) string {
	envFile := fmt.Sprintf("%s/%s/%s",
		CurrentDirectory(),
		stackDir,
		constants.DockerEnvFile,
	)
	var portNumber string
	content, err := os.ReadFile(envFile)
	errorhandler.CheckNilErr(err)
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if StartsWith(line, portName) {
			portLine := strings.Split(line, "=")
			portNumber = portLine[len(portLine)-1]
			return portNumber
		}
	}
	if portName == constants.BackendPort {
		portNumber = strconv.Itoa(constants.BackendPortNumber)
	} else if portName == constants.PostgresPort {
		portNumber = strconv.Itoa(constants.PostgresPortNumber)
	} else if portName == constants.MysqlPort {
		portNumber = strconv.Itoa(constants.MysqlPortNumber)
	} else if portName == constants.RedisPort {
		portNumber = strconv.Itoa(constants.RedisPortNumber)
	}
	return portNumber
}
