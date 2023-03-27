package helpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func CreateInfra(stack, service string) error {

	infraFiles := make(map[string]string)
	path := fileutils.CurrentDirectory()

	files := []string{
		constants.PackageDotJsonFile,
		constants.EnvFile,
		constants.SstConfigJsFile,
		constants.WebStackJsFile,
	}

	for _, file := range files {
		status, _ := fileutils.IsExists(path + "/" + file)
		if status {
			errorhandler.CheckNilErr(fmt.Errorf("%s file already exist", file))
		}
	}

	switch stack {
	case constants.ReactJS:

		infraFiles[constants.PackageDotJsonFile] = fmt.Sprintf(`{
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
	},
	"workspaces": [
		"%s/*"
	]
}`, service)

		infraFiles[constants.EnvFile] = `APP_NAME=app
WEB_AWS_REGION=us-east-1`

		infraFiles[constants.SstConfigJsFile] = `import dotenv from "dotenv";
import { WebStack } from "./stacks/WebStack";

dotenv.config({ path: ".env" });

export default {
	config(_input) {
		return {
			name: process.env.APP_NAME || "web-app",
			region: process.env.WEB_AWS_REGION || "us-east-1",
		};
	},
	stacks(app) {
		// deploy stacks
		app.stack(WebStack);
	},
};
`

		infraFiles[constants.WebStackJsFile] = fmt.Sprintf(`import { StaticSite } from "sst/constructs";

export function WebStack({ stack }) {
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

	default:
		return fmt.Errorf("Only react template is integrated now")
	}

	done := make(chan bool)
	go ProgressBar(30, "Generating", done)

	for fileName, fileSource := range infraFiles {

		if fileName == constants.WebStackJsFile {
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
	fmt.Printf("\nGenerating completed\n")

	return nil
}
