package sources

import (
	"fmt"

	"github.com/wednesday-solutions/picky/internal/utils"
)

func CISource(stack, stackDir, environment string) string {
	envName := utils.GetShortEnvName(environment)
	source := fmt.Sprintf(`name: CI %s
on:
  push:
    branches:
      - develop
      - qa
      - master
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
        run: yarn run test
`,
		stackDir, stackDir, stackDir, stackDir,
		stackDir, envName,
	)
	return source
}
