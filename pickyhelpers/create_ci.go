package pickyhelpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
)

func CreateCI(stackDirs []string) error {
	projectName := utils.GetProjectName()
	currentDir := utils.CurrentDirectory()
	workflowsPath := fmt.Sprintf("%s/%s", currentDir,
		constants.GithubWorkflowsDir)
	rootCIPath := fmt.Sprintf("%s/ci-%s.yml",
		workflowsPath,
		projectName,
	)
	status, _ := utils.IsExists(rootCIPath)
	if !status {
		utils.CreateGithubWorkflowDir()
		err := CreateRootCI(rootCIPath, projectName)
		errorhandler.CheckNilErr(err)
	}
	var stackCIPath string
	for _, dir := range stackDirs {
		stackCIPath = fmt.Sprintf("%s/ci-%s.yml",
			workflowsPath,
			dir,
		)
		status, _ = utils.IsExists(stackCIPath)
		if !status {
			err := CreateStackCI(stackCIPath, dir)
			errorhandler.CheckNilErr(err)
		}
	}
	return nil
}

// CreateRootCI will create and write CI for root.
func CreateRootCI(path, projectName string) error {
	source := fmt.Sprintf(`name: CI %s
on:
  push:
    branches: ["master", "develop", "qa"]
  pull_request:

  # Allows to run this workflow manually from the Actions tab
  workflow_dispatch:

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v3

      - name: Run the logs
        run: SERVICE PICKER, %s
`, projectName, projectName)

	err := utils.WriteToFile(path, source)
	return err
}

// CreateStackCI creates and writes CI for the given stack.
func CreateStackCI(path, stackDir string) error {
	source := fmt.Sprintf(`name: CI %s
on:
  push:
    branches: ["master", "develop", "qa"]
    paths: "%s/**"
  pull_request:
    paths: "%s/**"

  # Allows to run this workflow manually from the Actions tab
  workflow_dispatch:

jobs:
  build_and_test:
    name: Build & Test
    runs-on: ubuntu-latest
    strategy:
      matrix:
        node-version: [14.x]

    steps:
      - uses: actions/checkout@v3
      - name: Install dependencies
        working-directory: ./%s
        run: yarn

      - name: Lint
        working-directory: ./%s
        run: yarn lint

      - name: Build
        working-directory: ./%s
        run: yarn build:dev

      - name: Test
        working-directory: ./%s
        run: yarn test
`, stackDir, stackDir, stackDir, stackDir, stackDir, stackDir, stackDir)

	err := utils.WriteToFile(path, source)
	return err
}
