package helpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func CreateInfrastructure(stack, service string) error {

	var infraFiles map[string]string
	path := fileutils.CurrentDirectory()

	files := []string{
		constants.PackageDotJsonFile,
		constants.EnvFile,
		constants.SstConfigJsFile,
		constants.FrontendStackJsFile,
	}

	for _, file := range files {
		status, _ := fileutils.IsExists(path + "/" + file)
		if status {
			errorhandler.CheckNilErr(fmt.Errorf("%s file already exist", file))
		}
	}

	switch stack {
	case constants.ReactJS:

		packageDotJsonSource := fmt.Sprintf(`{
	"name": "app",
	"version": "0.0.0",
	"private": true,
	"type": "module",
	"scripts": {
		"dev": "sst dev",
		"build": "sst build",
		"deploy": "sst deploy",
		"remove": "sst remove",
		"console": "sst console",
		"typecheck": "tsc --noEmit"
	},
	"devDependencies": {
		"sst": "^2.0.18",
		"aws-cdk-lib": "2.62.2",
		"constructs": "10.1.156",
		"dotenv": "^10.0.0",
		"typescript": "^4.9.5",
		"@tsconfig/node16": "^1.0.3"
	},
	"workspaces": [
		"%s/*"
	]
}`, service)

		envSource := `APP_NAME=app
WEB_AWS_REGION=us-east-1`

		sstConfigSource := `const dotenv = require('dotenv');
		
dotenv.config({ path: ".env" });

export default {
	config(_input) {
		return {
			name: process.env.APP_NAME || "frontend-app",
			region: process.env.WEB_AWS_REGION || "us-east-1",
		};
	},
};
`

		frontendStackSource := fmt.Sprintf(`import { StaticSite } from "sst/constructs";

export function FrontendStack({ stack }) {
	// Deploy our React app
	const site = new StaticSite(stack, "ReactSite", {
		path: "%s",
		buildCommand: "yarn run build",
		buildOutput: "build",
	});

	// Show the URLs in the output
	stack.addOutputs({
		SiteUrl: site.url || "http://localhost:3000/",
	});
}
`, service)

		infraFiles = map[string]string{
			constants.PackageDotJsonFile:  packageDotJsonSource,
			constants.EnvFile:             envSource,
			constants.SstConfigJsFile:     sstConfigSource,
			constants.FrontendStackJsFile: frontendStackSource,
		}
	default:
		return fmt.Errorf("Only react template is integrated now")
	}

	done := make(chan bool)
	go ProgressBar(30, "Generating", done)

	for fileName, fileSource := range infraFiles {

		if fileName == "FrontendStack.js" {
			err := fileutils.MakeDirectory(path, "stacks")
			errorhandler.CheckNilErr(err)
			path = fileutils.CurrentDirectory() + "/stacks"

		} else {
			path = fileutils.CurrentDirectory()
		}
		err := fileutils.TruncateAndWriteToFile(path, fileName, fileSource)
		errorhandler.CheckNilErr(err)
	}
	<-done

	return nil
}
