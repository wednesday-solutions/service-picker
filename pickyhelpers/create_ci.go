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
	var stackCIPath, stack string
	for _, dir := range stackDirs {
		stackCIPath = fmt.Sprintf("%s/ci-%s.yml",
			workflowsPath,
			dir,
		)
		status, _ = utils.IsExists(stackCIPath)
		if !status {
			stack, _ = utils.FindStackAndDatabase(dir)
			err := CreateStackCI(stackCIPath, dir, stack)
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
        run: echo SERVICE PICKER`, projectName)

	err := utils.WriteToFile(path, source)
	return err
}

// CreateStackCI creates and writes CI for the given stack.
func CreateStackCI(path, stackDir, stack string) error {
	var environment, source string
	if stack == constants.NodeExpressGraphqlTemplate {
		environment = constants.Development
	} else {
		environment = constants.Dev
	}
	if stack != constants.GolangEchoTemplate {
		source = fmt.Sprintf(`name: CI %s
on:
  push:
    branches: ["master", "develop", "qa"]
    paths: "%s/**"
  pull_request:
    paths: "%s/**"

  # Allows to run this workflow manually from the Actions tab
  workflow_dispatch:

jobs:
  build-and-test:
    name: Build & Test
    runs-on: ubuntu-latest
    strategy:
      matrix:
        node-version: [16.14.x]

    steps:
      - uses: actions/checkout@v3
      - name: Use Node.js ${{ matrix.node-version }}
        uses: actions/setup-node@v2
        with:
          node-version: ${{ matrix.node-version }}
          cache: 'yarn'
          cache-dependency-path: ./%s/package.json

      - name: Install dependencies
        working-directory: ./%s
        run: yarn

      - name: Lint
        working-directory: ./%s
        run: yarn lint

      - name: Build
        working-directory: ./%s
        run: yarn build:%s

      - name: Test
        working-directory: ./%s
        run: yarn test`,
			stackDir, stackDir, stackDir, stackDir, stackDir,
			stackDir, stackDir, environment, stackDir,
		)
	}

	err := utils.WriteToFile(path, source)
	return err
}
