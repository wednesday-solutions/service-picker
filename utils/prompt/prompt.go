package prompt

import (
	"fmt"
	"log"
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
	var cloudProviderConfigLabel string = "Choose a cloud provider config"
	var cloudProviderConfigItems = []string{constants.CREATE_CD, constants.CREATE_INFRA}

	selectedCloudConfig := PromptSelect(cloudProviderConfigLabel, cloudProviderConfigItems)
	dirName := service
	if dirName != constants.BACKEND {
		dirName = constants.FRONTEND
	}
	if selectedCloudConfig == constants.CREATE_CD {
		cdSource := "workflows/" + dirName + "/cd/" + stack + ".yml"
		cdDestination := fileutils.CurrentDirectory() + "/" + dirName + "/.github/workflows/cd.yml"

		status, _ := fileutils.IsExists(cdDestination)
		if !status {
			err := cp.Copy(cdSource, cdDestination)
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
	var cloudProviderLabel string = "Choose a cloud provider"
	var cloudProviderItems = []string{constants.AWS}

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

	if service == constants.BACKEND {
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

	var configLabel string = "Choose the config to setup"
	var configItems = []string{constants.INIT, constants.CLOUD_NATIVE}

	selectedConfig := PromptSelect(configLabel, configItems)

	if selectedConfig == constants.INIT {
		PromptSelectInit(service, stack, database)
	} else {
		PromptSelectCloudProvider(service, stack)
	}
}

func PromptSelectStackDatabase(service, stack string) {
	var database string
	var label string = "Choose a database"
	if service == constants.BACKEND {
		switch stack {
		case constants.NODE_HAPI:
			database = PromptSelect(label, []string{constants.MYSQL})
		case constants.NODE_EXPRESS:
			database = PromptSelect(label, []string{constants.POSTGRES})
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
