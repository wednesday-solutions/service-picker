package helpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func CreateInfrastructure(stack, dirName, database string) error {

	path := fileutils.CurrentDirectory()
	var infraFiles map[string]string

	switch stack {
	case constants.REACT:

		packageDotJsonSource := fmt.Sprintf(`{
	"name": "react-app",
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
		"typescript": "^4.9.5",
		"@tsconfig/node16": "^1.0.3"
	},
	"workspaces": [
		"%s/*"
	]
}`, dirName)

		sstConfigSource := `export default {
	config(_input) {
		return {
			name: "react-app",
			region: "us-east-1",
		};
	},
};
`

		frontendStackSource := `import { StaticSite } from "sst/constructs";

export function FrontendStack({ stack }) {
	// Deploy our React app
	const site = new StaticSite(stack, "ReactSite", {
		path: "frontend",
		buildCommand: "yarn run build",
		buildOutput: "build",
	});

	// Show the URLs in the output
	stack.addOutputs({
		SiteUrl: site.url || "http://localhost:3000/",
	});
}
`

		infraFiles = map[string]string{
			"package.json":     packageDotJsonSource,
			"sst.config.js":    sstConfigSource,
			"FrontendStack.js": frontendStackSource,
		}
	default:
		return fmt.Errorf("Only react template is integrated now")
	}

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
	return nil
}
