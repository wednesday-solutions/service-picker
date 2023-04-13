package pickyhelpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
)

func CreateCI(stackDirs []string) error {
	currentDir := utils.CurrentDirectory()
	workflowsPath := fmt.Sprintf("%s/%s", currentDir,
		constants.GithubWorkflowsDir)

	status, _ := utils.IsExists(workflowsPath)
	if !status {
		githubFolderPath := fmt.Sprintf("%s/%s", currentDir, ".github")
		githubStatus, _ := utils.IsExists(githubFolderPath)
		if !githubStatus {
			err := utils.CreateDirectory(githubFolderPath)
			errorhandler.CheckNilErr(err)
		}
		err := utils.CreateDirectory(workflowsPath)
		errorhandler.CheckNilErr(err)
	}
	var stackCIPath string
	for _, dir := range stackDirs {
		stackCIPath = fmt.Sprintf("%s/ci-%s.yaml",
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
        run: yarn build:dev

      - name: Test
        run: yarn test
`, stackDir, stackDir, stackDir, stackDir, stackDir)

	err := utils.WriteToFile(path, source)
	return err
}
