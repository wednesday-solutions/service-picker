package utils

import (
	"fmt"

	cp "github.com/otiai10/copy"

	"github.com/manifoldco/promptui"
)

var (
		INIT = "init"
		CLOUD_NATIVE = "cloud native"
		AWS = "AWS"
		CREATE_CD = "Create CD pipeline"
		CREATE_INFRA = "Create Infra"
	)


func PromptSelect(label string, items []string ) string  {
	
	prompt := promptui.Select{
	Label: label,
	Items: items,
}

_, result, err:= prompt.Run()

checkNilErr(err)

return result
}

func PromptSelectCloudProviderConfig(service string, stack string) {
	var cloudProviderConfigLabel string = "Choose a cloud provider config"
	var cloudProviderConfigItems = []string{CREATE_CD, CREATE_INFRA}

	selectedCloudConfig := PromptSelect(cloudProviderConfigLabel, cloudProviderConfigItems)
		if selectedCloudConfig == CREATE_CD {
			cdSource:= "workflows/" + service +  "/cd/" + stack + ".yml"
			cdDestination := CurrentDirectory() + "/" + service + "/.github/workflows/cd.yml"

			status, _ := IsExists(cdDestination)
			if !status {
				err := cp.Copy(cdSource, cdDestination)
				checkNilErr(err)
			} else {
				fmt.Println("Error:", "The", service, stack, "CD you are looking to create already exists")
			}
			

			} else if selectedCloudConfig == CREATE_INFRA {
		      infraSource:= "infrastructure/" + service
			  infraDestination:= CurrentDirectory() + "/"
			  status,_ := IsExists(infraDestination + "/stacks")
			  if !status {
				  err := cp.Copy(infraSource, infraDestination)
				  checkNilErr(err)
			  } else {
				  fmt.Println("Error:", "The", service, stack, "infrastructure you are looking to create already exists")
			  }
			  
			  
			}
			
		}

func PromptSelectCloudProvider(service string, stack string){
	var cloudProviderLabel string = "Choose a cloud provider"
	var cloudProviderItems = []string{AWS}

    selectedCloudProvider := PromptSelect(cloudProviderLabel, cloudProviderItems)
	if selectedCloudProvider == AWS {
		PromptSelectCloudProviderConfig(service, stack)
	}
}

func PromptSelectStackConfig(service string, stack string) {
	
	var configLabel string = "Choose the config to setup"
	var configItems = []string{INIT, CLOUD_NATIVE}

	selectedConfig := PromptSelect(configLabel, configItems)

	if selectedConfig == INIT {
		destination := CurrentDirectory() + "/" + service
		status, _ := IsExists(destination)
		if !status {
			source:= "services/" + service +  "/" + stack
		    err := cp.Copy(source, destination)
		    checkNilErr(err)
		} else {
			fmt.Println("Error:", "The", service, "service already exists. You can initialise only one stack in a service")
		}
		
		
	} else {
		PromptSelectCloudProvider(service, stack)
	}
}


func PromptSelectStack(service string, items []string){
	 stack := PromptSelect("Pick a stack", items)
	 PromptSelectStackConfig(service, stack)
}
