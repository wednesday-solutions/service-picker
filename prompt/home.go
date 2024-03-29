package prompt

import (
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/utils"
)

func PromptHome() {
	var p PromptInput
	p.Label = "Pick an option"
	p.GoBack = PromptAlertMessage
	var initService bool
	stacks, _, _ := utils.GetExistingStacksDatabasesAndDirectories()
	if len(stacks) > 0 {
		p.Items = []string{constants.InitService}
		if ShowPromptGitInit() {
			p.Items = append(p.Items, constants.GitInit)
		}
		p.Items = append(p.Items,
			constants.CICD,
			constants.DockerCompose,
			constants.SetupInfra,
			constants.Deploy,
		)
		if ShowRemoveDeploy() {
			p.Items = append(p.Items, constants.RemoveDeploy)
		}
		p.Items = append(p.Items, constants.Exit)
		response, _ := p.PromptSelect()
		switch response {
		case constants.InitService:
			initService = true
		case constants.DockerCompose:
			PromptDockerCompose()
		case constants.CICD:
			PromptCICD()
		case constants.SetupInfra:
			PromptSetupInfra()
		case constants.Deploy:
			PromptDeploy()
		case constants.RemoveDeploy:
			PromptRemoveDeploy()
		case constants.GitInit:
			PromptGitInit()
		case constants.Exit:
			PromptExit()
		}
	}
	if len(stacks) == 0 || initService {
		PromptSelectService()
	}
}
