package prompt

import "github.com/wednesday-solutions/picky/utils"

func PromptHome() {
	var label string
	var initService bool
	stacks, databases, _ := utils.ExistingStacksDatabasesAndDirectories()
	if len(stacks) > 0 {
		label = "Pick an option"
		items := []string{"Init Service"}
		showCreateDC := ShowCreateDockerCompose(databases)
		if showCreateDC {
			items = append(items, "Create docker-compose")
		}
		items = append(items, "Create CD", "Setup Infra", "Deploy", "Exit")
		response := PromptSelect(label, items)
		switch response {
		case "Init Service":
			initService = true
		case "Create docker-compose":
			PromptDockerCompose()
		case "Create CD":
			PromptCreateCD()
		case "Setup Infra":
			PromptSetupInfra()
		case "Deploy":
			PromptDeploy()
		case "Exit":
			PromptExit()
		}
	}
	if len(stacks) == 0 || initService {
		PromptSelectService()
	}
}
