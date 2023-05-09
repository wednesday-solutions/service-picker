package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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

var (
	WebPortNumber      = 3000
	BackendPortNumber  = 9000
	PostgresPortNumber = 5432
	MysqlPortNumber    = 3306
)

// ResetPortNumbers resets the port numbers to default.
func ResetPortNumbers() {
	WebPortNumber = 3000
	BackendPortNumber = 9000
	PostgresPortNumber = 5432
	MysqlPortNumber = 3306
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
	_, _, directories := ExistingStacksDatabasesAndDirectories()
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

func WebOutputsSource(webObj WebOutputKeys, env string) string {
	source := fmt.Sprintf(`{
	"%s": {
		"distributionId": "%s",
		"bucketName": "%s",
		"siteUrl": "%s",
	}
}`,
		env,
		webObj.DistributionId,
		webObj.BucketName,
		webObj.SiteUrl,
	)
	return source
}

func BackendOutputsSource(backendObj BackendOutputKeys, env string) string {
	source := fmt.Sprintf(`{
	"%s": {
		"redisHost": "%s",
    "image": "%s",
    "awsRegion": "%s",
    "secretName": "%s",
    "databaseHost": "%s",
    "databaseName": "%s",
    "secretArn": "%s",
    "logDriver": "%s",
    "loadBalancerDns": "%s",
    "serviceName": "%s",
    "containerName": "%s",
    "taskRole": "%s",
    "clusterName": "%s",
    "elasticContainerRegistryRepo": "%s",
    "containerPort": "%s",
    "executionRole": "%s",
    "taskDefinition": "%s",
    "family": "%s",
		"logDriverOptions": {
      "awslogs-group": "%s",
      "awslogs-stream-prefix": "%s",
      "awslogs-region": "%s"
    },
	}
}`,
		env,
		backendObj.RedisHost,
		backendObj.Image,
		backendObj.AwsRegion,
		backendObj.SecretName,
		backendObj.DatabaseHost,
		backendObj.DatabaseName,
		backendObj.SecretArn,
		backendObj.LogDriver,
		backendObj.LoadBalancerDns,
		backendObj.ServiceName,
		backendObj.ContainerName,
		backendObj.TaskRole,
		backendObj.ClusterName,
		backendObj.ElasticContainerRegistryRepo,
		backendObj.ContainerPort,
		backendObj.ExecutionRole,
		backendObj.TaskDefinition,
		backendObj.Family,
		backendObj.LogDriverOptions.AwsLogsGroup,
		backendObj.LogDriverOptions.AwsLogsStreamPrefix,
		backendObj.LogDriverOptions.AwsLogsRegion,
	)
	return source
}
