package pickyhelpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
)

func CreateCI(stackDirs []string) error {
	workflowsPath := fmt.Sprintf("%s/%s", utils.CurrentDirectory(),
		constants.GithubWorkflowsDir)

	utils.CreateGithubWorkflowDir()
	var stackCIPath, stack string
	var status bool
	for _, dir := range stackDirs {
		stackCIPath = fmt.Sprintf("%s/ci-%s.yml", workflowsPath, dir)
		status, _ = utils.IsExists(stackCIPath)
		if !status {
			stack, _ = utils.FindStackAndDatabase(dir)
			err := CreateStackCI(stackCIPath, dir, stack)
			errorhandler.CheckNilErr(err)
		}
	}
	return nil
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
    defaults: 
      run:
        working-directory: ./%s
    strategy:
      matrix:
        node-version: [16.14.x]

    steps:
      - uses: actions/checkout@v3
      - name: Use Node.js ${{ matrix.node-version }}
        uses: actions/setup-node@v2
        with:
          node-version: ${{ matrix.node-version }}
          cache: "yarn"
          cache-dependency-path: ./%s/package.json

      - name: Install dependencies
        run: yarn

      - name: Lint
        run: yarn lint

      - name: Build
        run: yarn build:%s

      - name: Test
        run: yarn test
`, stackDir, stackDir, stackDir, stackDir, stackDir, environment)

		err := utils.WriteToFile(path, source)
		errorhandler.CheckNilErr(err)
	} else {
		err := utils.PrintWarningMessage(fmt.Sprintf(
			"CI of '%s' is in work in progress..!", stack,
		))
		errorhandler.CheckNilErr(err)
	}
	return nil
}
