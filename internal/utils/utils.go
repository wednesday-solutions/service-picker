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

// GetEnvironment return short environment name for the given environment.
func GetEnvironment(environment string) string {
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

func ReadJsonDataInSstOutputs() interface{} {
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

	return data
}

type LogDriverOptionsKeys struct {
	AwsLogsGroup        string
	AwsLogsStreamPrefix string
	AwsLogsRegion       string
}

type WebOutputKeys struct {
	DistributionId string
	BucketName     string
	SiteUrl        string
}

type BackendOutputKeys struct {
	TaskRole                     string
	Image                        string
	ContainerName                string
	ContainerPort                string
	ExecutionRole                string
	TaskDefinition               string
	LogDriver                    string
	LogDriverOptions             LogDriverOptionsKeys
	Family                       string
	AwsRegion                    string
	RedisHost                    string
	SecretName                   string
	DatabaseHost                 string
	DatabaseName                 string
	SecretArn                    string
	LoadBalancerDns              string
	ServiceName                  string
	ClusterName                  string
	ElasticContainerRegistryRepo string
}

func CreateSstOutputsFile() error {
	data := ReadJsonDataInSstOutputs()
	if data == nil {
		return nil
	}
	jsonData, ok := data.(map[string]interface{})
	if !ok {
		errorhandler.CheckNilErr(fmt.Errorf("outputs.json is not valid"))
	}
	_, _, directories := GetExistingStacksDatabasesAndDirectories()
	for _, stackDir := range directories {
		service := FindService(stackDir)

		for key, value := range jsonData {
			outputsFile := fmt.Sprintf("%s/%s-%s", CurrentDirectory(), stackDir, constants.OutputsJson)
			status, _ := IsExists(outputsFile)
			if !status {
				err := CreateFile(outputsFile)
				errorhandler.CheckNilErr(err)
			}
			env, source := "", ""
			if StartsWith(key, "dev") {
				env = constants.Dev
			} else if StartsWith(key, "prod") {
				env = constants.Prod
			}
			if EndsWith(key, "Pg") || EndsWith(key, "Mysql") {
				if service == constants.Backend {
					backendObj := ParseBackendOutputsKey(key, value)
					source = BackendOutputsSource(backendObj, env)
				}
			} else if EndsWith(key, "Web") {
				if service == constants.Web {
					webObj := ParseWebOutputsKey(key, value)
					source = WebOutputsSource(webObj, env)
				}
			}
			if source != "" {
				err := WriteToFile(outputsFile, source)
				errorhandler.CheckNilErr(err)
			}
		}
	}
	return nil
}

func ParseBackendOutputsKey(key string, value interface{}) BackendOutputKeys {
	var backendObj BackendOutputKeys
	backendJson, ok := value.(map[string]interface{})
	if !ok {
		errorhandler.CheckNilErr(fmt.Errorf("outputs.json is not valid"))
	}
	backendObj.TaskRole, _ = backendJson["taskRole"].(string)
	backendObj.Image, _ = backendJson["image"].(string)
	backendObj.ContainerName, _ = backendJson["containerName"].(string)
	backendObj.ContainerPort, _ = backendJson["containerPort"].(string)
	backendObj.ExecutionRole, _ = backendJson["executionRole"].(string)
	backendObj.TaskDefinition, _ = backendJson["taskDefinition"].(string)
	backendObj.LogDriver, _ = backendJson["logDriver"].(string)
	backendObj.Family, _ = backendJson["family"].(string)
	backendObj.AwsRegion, _ = backendJson["awsRegion"].(string)
	backendObj.RedisHost, _ = backendJson["redisHost"].(string)
	backendObj.SecretArn, _ = backendJson["secretArn"].(string)
	backendObj.DatabaseHost, _ = backendJson["databaseHost"].(string)
	backendObj.DatabaseName, _ = backendJson["databaseName"].(string)
	backendObj.SecretName, _ = backendJson["secretName"].(string)
	backendObj.LoadBalancerDns, _ = backendJson["loadBalancerDns"].(string)
	backendObj.ServiceName, _ = backendJson["serviceName"].(string)
	backendObj.ClusterName, _ = backendJson["clusterName"].(string)
	backendObj.ElasticContainerRegistryRepo, _ = backendJson["elasticContainerRegistryRepo"].(string)

	logdriveroptions, _ := backendJson["logDriverOptions"].(map[string]interface{})
	backendObj.LogDriverOptions.AwsLogsGroup, _ = logdriveroptions["awslogs-group"].(string)
	backendObj.LogDriverOptions.AwsLogsStreamPrefix, _ = logdriveroptions["awslogs-stream-prefix"].(string)
	backendObj.LogDriverOptions.AwsLogsRegion, _ = logdriveroptions["awslogs-region"].(string)

	return backendObj
}

func ParseWebOutputsKey(key string, value interface{}) WebOutputKeys {
	var webObj WebOutputKeys
	webJson, ok := value.(map[string]interface{})
	if !ok {
		errorhandler.CheckNilErr(fmt.Errorf("outputs.json is not valid"))
	}
	webObj.DistributionId, _ = webJson["distributionId"].(string)
	webObj.BucketName, _ = webJson["bucketName"].(string)
	webObj.SiteUrl, _ = webJson["siteUrl"].(string)

	return webObj
}

func WebOutputsSource(webObj WebOutputKeys, key string) string {
	source := fmt.Sprintf(`{
	"%s": {
		"siteUrl": "%s",
		"bucketName": "%s",
		"distributionId": "%s"
	}
}
`,
		key,
		webObj.SiteUrl,
		webObj.BucketName,
		webObj.DistributionId,
	)
	return source
}

func BackendOutputsSource(backendObj BackendOutputKeys, key string) string {
	source := fmt.Sprintf(`{
	"%s": {
    "image": "%s",
    "family": "%s",
    "taskRole": "%s",
    "executionRole": "%s",
    "databaseName": "%s",
    "databaseHost": "%s",
		"redisHost": "%s",
    "awsRegion": "%s",
    "secretName": "%s",
    "secretArn": "%s",
    "loadBalancerDns": "%s",
    "serviceName": "%s",
    "containerPort": "%s",
    "containerName": "%s",
    "clusterName": "%s",
    "taskDefinition": "%s",
    "elasticContainerRegistryRepo": "%s",
    "logDriver": "%s",
		"logDriverOptions": {
      "awslogs-group": "%s",
      "awslogs-stream-prefix": "%s",
      "awslogs-region": "%s"
    }
	}
}
`,
		key,
		backendObj.Image,
		backendObj.Family,
		backendObj.TaskRole,
		backendObj.ExecutionRole,
		backendObj.DatabaseName,
		backendObj.DatabaseHost,
		backendObj.RedisHost,
		backendObj.AwsRegion,
		backendObj.SecretName,
		backendObj.SecretArn,
		backendObj.LoadBalancerDns,
		backendObj.ServiceName,
		backendObj.ContainerPort,
		backendObj.ContainerName,
		backendObj.ClusterName,
		backendObj.TaskDefinition,
		backendObj.ElasticContainerRegistryRepo,
		backendObj.LogDriver,
		backendObj.LogDriverOptions.AwsLogsGroup,
		backendObj.LogDriverOptions.AwsLogsStreamPrefix,
		backendObj.LogDriverOptions.AwsLogsRegion,
	)
	return source
}

type TaskDefinitionDetails struct {
	BackendObj  BackendOutputKeys
	Environment string
	EnvName     string
	SecretName  string
}

func GetOutputsBackendObject(environment, stackDir string) TaskDefinitionDetails {
	stackKey := strcase.ToCamel(stackDir)
	data := ReadJsonDataInSstOutputs()
	if data == nil {
		return TaskDefinitionDetails{}
	}
	jsonData, ok := data.(map[string]interface{})
	if !ok {
		errorhandler.CheckNilErr(fmt.Errorf("outputs.json is not valid."))
	}
	var taskDefinition TaskDefinitionDetails
	for key, value := range jsonData {
		if EndsWith(key, stackKey) {
			if EndsWith(key, "Pg") {
				taskDefinition.SecretName = "POSTGRES_PASSWORD"
			} else if EndsWith(key, "Mysql") {
				taskDefinition.SecretName = "MYSQL_PASSWORD"
			}
			if environment == constants.Develop {
				taskDefinition.EnvName = constants.Dev
			} else if environment == constants.Production {
				taskDefinition.EnvName = constants.Prod
			}
			taskDefinition.BackendObj = ParseBackendOutputsKey(key, value)
			taskDefinition.Environment = environment

			return taskDefinition
		}
	}
	return TaskDefinitionDetails{}
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

func CreateJsonFile() error {
	data := ReadJsonDataInSstOutputs()
	if data == nil {
		return nil
	}
	jsonData, ok := data.(map[string]interface{})
	if !ok {
		errorhandler.CheckNilErr(fmt.Errorf("outputs.json is not valid"))
	}
	_, _, directories := GetExistingStacksDatabasesAndDirectories()
	for _, stackDir := range directories {
		service := FindService(stackDir)

		for key, value := range jsonData {
			outputsFile := fmt.Sprintf("%s/%s-%s", CurrentDirectory(), key, constants.OutputsJson)
			status, _ := IsExists(outputsFile)
			if !status {
				err := CreateFile(outputsFile)
				errorhandler.CheckNilErr(err)
			}
			var source string
			if EndsWith(key, "Pg") || EndsWith(key, "Mysql") {
				if service == constants.Backend {
					backendObj := ParseBackendOutputsKey(key, value)
					// source = BackendOutputsSource(backendObj, key)

					for key, value := range backendObj {
						source = fmt.Sprintf(`\t"%s": "%s"`)
					}

					err := WriteToFile(outputsFile, fmt.Sprintf("%+v\n", backendObj))
					errorhandler.CheckNilErr(err)
				}
			} else if EndsWith(key, "Web") {
				if service == constants.Web {
					webObj := ParseWebOutputsKey(key, value)
					// source = WebOutputsSource(webObj, key)
					err := WriteToFile(outputsFile, fmt.Sprintf("%+v\n", webObj))
					errorhandler.CheckNilErr(err)
				}
			}
		}
	}
	return nil
}
