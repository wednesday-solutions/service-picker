package flagcommand

import (
	"fmt"
	"strings"

	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
	"github.com/wednesday-solutions/picky/prompt"
)

type InitInfo struct {
	Service   string
	Stack     string
	Database  string
	Directory string
}

func (i *InitInfo) FlagInit() error {

	i.Database = utils.GetDatabase(i.Database)

	i.Stack = utils.GetStackConstantNameFromLower(i.Stack)
	if i.Stack == "" {
		return fmt.Errorf("Entered stack is invalid")
	}
	if i.Stack == constants.GolangEchoTemplate {
		i.Stack = fmt.Sprintf("%s-%s", strings.Split(i.Stack, " ")[0], i.Database)
	}
	i.Directory = utils.CreateStackDirectory(i.Directory, i.Stack, i.Database)

	status, _ := utils.IsExists(fmt.Sprintf("%s/%s", utils.CurrentDirectory(), i.Directory))
	if status {
		return fmt.Errorf("Entered directory %s already exists\n", i.Directory)
	}
	err := utils.MakeDirectory(utils.CurrentDirectory(), i.Directory)
	errorhandler.CheckNilErr(err)

	fmt.Printf("\nService: %s\nStack: %s\nDatabase: %s\nDirectory: %s\n\n",
		i.Service, i.Stack, i.Database, i.Directory,
	)

	var ii prompt.InitInfo
	ii.Service = i.Service
	ii.Stack = i.Stack
	ii.Database = i.Database
	ii.DirName = i.Directory
	err = ii.StackInitialize()
	return err
}
