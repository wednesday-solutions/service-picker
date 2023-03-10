package prompt

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
	cp "github.com/otiai10/copy"
	"github.com/stoewer/go-strcase"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
	"github.com/wednesday-solutions/picky/utils/helpers"
)

func PromptSelect(label string, items []string) string {

	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	_, result, err := prompt.Run()

	errorhandler.CheckNilErr(err)

	return result
}

func PromptSelectCloudProviderConfig(service string, stack string) {
	cloudProviderConfigLabel := "Choose a cloud provider config"
	cloudProviderConfigItems := []string{constants.CREATE_CD, constants.CREATE_INFRA}

	selectedCloudConfig := PromptSelect(cloudProviderConfigLabel, cloudProviderConfigItems)
	dirName := service
	if dirName != constants.BACKEND {
		dirName = constants.FRONTEND
	}
	if selectedCloudConfig == constants.CREATE_CD {

		githubUrl := "https://raw.githubusercontent.com/wednesday-solutions/"
		var cdFileUrl string
		cdFile := "/.github/workflows/cd.yml"

		switch stack {
		case constants.NODE_HAPI:
			cdFileUrl = githubUrl + "nodejs-hapi-template/main" + cdFile
		case constants.NODE_EXPRESS:
			cdFileUrl = githubUrl + "node-express-graphql-template/develop" + cdFile
		case constants.REACT:
			cdFileUrl = githubUrl + "react-graphql-ts-template/master" + cdFile
		default:
			cdFileUrl = ""
		}

		cdDestination := fileutils.CurrentDirectory() + "/" + dirName + cdFile
		status, _ := fileutils.IsExists(cdDestination)
		if !status {
			resp, err := http.Get(cdFileUrl)
			errorhandler.CheckNilErr(err)
			defer resp.Body.Close()

			cdFileData, err := io.ReadAll(resp.Body)
			errorhandler.CheckNilErr(err)

			// Create CD File
			err = fileutils.CreateFile(cdDestination)
			errorhandler.CheckNilErr(err)

			// Write CDFileData to CD File
			err = fileutils.WriteToFile(cdDestination, string(cdFileData))
			errorhandler.CheckNilErr(err)

		} else {
			fmt.Println("The", dirName, stack, "CD you are looking to create already exists")
		}
	} else if selectedCloudConfig == constants.CREATE_INFRA {
		infraSource := "infrastructure/" + dirName
		infraDestination := fileutils.CurrentDirectory() + "/"
		status, _ := fileutils.IsExists(infraDestination + "/stacks")
		if !status {
			err := cp.Copy(infraSource, infraDestination)
			errorhandler.CheckNilErr(err)
		} else {
			fmt.Println("The", dirName, stack, "infrastructure you are looking to create already exists")
		}
	}
}

func PromptSelectCloudProvider(service string, stack string) {
	cloudProviderLabel := "Choose a cloud provider"
	cloudProviderItems := []string{constants.AWS}

	selectedCloudProvider := PromptSelect(cloudProviderLabel, cloudProviderItems)
	if selectedCloudProvider == constants.AWS {
		PromptSelectCloudProviderConfig(service, stack)
	}
}

func PromptSelectInit(service, stack, database string) {

	currentDir := fileutils.CurrentDirectory()
	splitDirs := strings.Split(currentDir, "/")
	projectName := splitDirs[len(splitDirs)-1]
	projectName = strcase.SnakeCase(projectName)

	if stack == constants.GOLANG {
		stack = fmt.Sprintf("%s-%s", strings.Split(stack, " ")[0], database)
	}

	var createDockerFile bool
	dirName := service
	if service != constants.BACKEND {
		dirName = constants.FRONTEND
	}
	destination := currentDir + "/" + dirName

	status, _ := fileutils.IsExists(destination)
	if !status {
		makeDirErr := fileutils.MakeDirectory(currentDir+"/", dirName)
		errorhandler.CheckNilErr(makeDirErr)
		cmd := exec.Command("git", "clone", constants.Repos()[stack], dirName)
		err := cmd.Run()
		errorhandler.CheckNilErr(err)

		// Delete cd.yml file from the cloned repo.
		cdFilePatch := currentDir + "/" + dirName + "/.github/workflows/cd.yml"
		status, _ := fileutils.IsExists(cdFilePatch)
		if status {
			err = fileutils.RemoveFile(cdFilePatch)
			errorhandler.CheckNilErr(err)
		}

		// Database selection
		if stack == constants.NODE_HAPI && database == constants.POSTGRES {
			// Convert DB Connection of Hapi template into Postgres.
			file := "/backend/config/db.js"
			err = helpers.ConvertDBConnection(stack, file, database, projectName)
			errorhandler.CheckNilErr(err)

			err = helpers.UpdateEnvFiles(stack, database, projectName)
			errorhandler.CheckNilErr(err)

			err = helpers.ConvertMysqlToPostgres(stack, projectName)
			errorhandler.CheckNilErr(err)
		}

		// Docker-compose file
		if dirName == constants.FRONTEND {
			destination = currentDir + "/" + constants.BACKEND
			status, _ := fileutils.IsExists(destination)
			if status {
				createDockerFile = true
			}
		} else if dirName == constants.BACKEND {
			destination = currentDir + "/" + constants.FRONTEND
			status, _ := fileutils.IsExists(destination)
			if status {
				createDockerFile = true
			}
		}
		if createDockerFile {
			// create Docker File
			dockerComposeFile := "docker-compose.yml"
			err = fileutils.MakeFile(currentDir, dockerComposeFile)
			errorhandler.CheckNilErr(err)

			// write Docker File
			err = helpers.WriteDockerFile(dockerComposeFile, database, projectName)
			errorhandler.CheckNilErr(err)
		}
	} else {
		fmt.Println("The", service, "service already exists. You can initialize only one stack in a service")
	}
}

func PromptSelectStackConfig(service, stack, database string) {
	configLabel := "Choose the config to setup"
	configItems := []string{constants.INIT, constants.CLOUD_NATIVE}

	selectedConfig := PromptSelect(configLabel, configItems)

	if selectedConfig == constants.INIT {
		PromptSelectInit(service, stack, database)
	} else {
		PromptSelectCloudProvider(service, stack)
	}
}

func PromptSelectStackDatabase(service, stack string) {
	var database string
	label := "Choose a database"
	if service == constants.BACKEND {
		switch stack {
		case constants.NODE_HAPI:
			database = PromptSelect(label, []string{constants.POSTGRES, constants.MYSQL})
		case constants.NODE_EXPRESS:
			database = PromptSelect(label, []string{constants.POSTGRES, constants.MYSQL})
		case constants.NODE_EXPRESS_TS:
			database = PromptSelect(label, []string{})
		case constants.GOLANG:
			database = PromptSelect(label, []string{constants.POSTGRES, constants.MYSQL})
		default:
			log.Fatalln("Something went wrong")
		}
	} else {
		switch stack {
		case constants.REACT, constants.NEXT:
			database = PromptSelect(label, []string{constants.POSTGRES, constants.MYSQL, constants.MONGODB})
		default:
			log.Fatalln("Something went wrong")
		}
	}
	PromptSelectStackConfig(service, stack, database)
}

func PromptSelectStack(service string, items []string) {
	stack := PromptSelect("Pick a stack", items)

	var status bool
	var err error
	if service != constants.BACKEND {
		status, err = fileutils.IsExists(fileutils.CurrentDirectory() + "/" + constants.BACKEND)
		errorhandler.CheckNilErr(err)
	}

	// Choose database
	if status || service == constants.BACKEND {
		PromptSelectStackDatabase(service, stack)
	} else {
		PromptSelectStackConfig(service, stack, "")
	}
}
