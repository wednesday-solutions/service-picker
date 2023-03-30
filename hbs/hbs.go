package hbs

import (
	"github.com/aymerick/raymond"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func init() {
	raymond.RegisterHelper("databaseVolume", DatabaseVolume)
	raymond.RegisterHelper("dbVersion", DBVersion)
	raymond.RegisterHelper("portConnection", PortConnection)
	raymond.RegisterHelper("dbServiceName", DBServiceName)
	raymond.RegisterHelper("globalAddDependencies", GlobalAddDependencies)
	raymond.RegisterHelper("addDependencies", AddDependencies)
	raymond.RegisterHelper("envFileBackend", EnvFileBackend)
	raymond.RegisterHelper("runBuildEnvironment", RunBuildEnvironment)
	raymond.RegisterHelper("waitForDBService", WaitForDBService)
	raymond.RegisterHelper("dependsOnFieldOfGo", DependsOnFieldOfGo)
	raymond.RegisterHelper("cmdDockerfile", CmdDockerfile)
	raymond.RegisterHelper("envEnvironmentName", EnvEnvironmentName)
	raymond.RegisterHelper("deployStacks", DeployStacks)
	raymond.RegisterHelper("sstImportStacks", SstImportStacks)
}

func ParseAndWriteToFile(source, filePath string, stackInfo map[string]interface{}) error {

	ctx := map[string]interface{}{
		constants.Frontend:                 constants.Frontend,
		constants.Web:                      constants.Web,
		constants.Mobile:                   constants.Mobile,
		constants.Backend:                  constants.Backend,
		constants.Redis:                    constants.Redis,
		constants.Postgres:                 constants.Postgres,
		constants.Mysql:                    constants.Mysql,
		constants.GolangMySQLTemplate:      constants.GolangMySQLTemplate,
		constants.GolangPostgreSQLTemplate: constants.GolangPostgreSQLTemplate,
		constants.Stack:                    stackInfo[constants.Stack].(string),
		constants.Database:                 stackInfo[constants.Database].(string),
		constants.ProjectName:              stackInfo[constants.ProjectName].(string),
		constants.WebStatus:                stackInfo[constants.WebStatus].(bool),
		constants.MobileStatus:             stackInfo[constants.MobileStatus].(bool),
		constants.BackendStatus:            stackInfo[constants.BackendStatus].(bool),
		constants.WebDirName:               stackInfo[constants.WebDirName].(string),
		constants.MobileDirName:            stackInfo[constants.MobileDirName].(string),
		constants.BackendDirName:           stackInfo[constants.BackendDirName].(string),
		constants.SstConfigStack:           stackInfo[constants.SstConfigStack].(string),
		constants.ExistingDirectories:      stackInfo[constants.ExistingDirectories].([]string),
	}
	// Parse the source string into template
	tpl, err := raymond.Parse(source)
	errorhandler.CheckNilErr(err)

	// Execute the template into string
	executedTemplate, err := tpl.Exec(ctx)
	errorhandler.CheckNilErr(err)

	err = fileutils.WriteToFile(filePath, executedTemplate)
	errorhandler.CheckNilErr(err)

	return nil
}
