package helpers

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func UpdatePackageDotJson(stack, jsonFile, database, projectName string) error {

	jsonFile = fmt.Sprintf("%s%s", fileutils.CurrentDirectory(), jsonFile)

	input, err := ioutil.ReadFile(jsonFile)
	errorhandler.CheckNilErr(err)

	lines := strings.Split(string(input), "\n")
	var updatePackageDotJson bool

	switch stack {
	case constants.NODE_HAPI:

		if database == constants.POSTGRES {
			updatePackageDotJson = true
			lines[43] = fmt.Sprintf("%s\n\t\t%s", lines[43], `"cls-hooked": "4.2.2",`)
			lines[62] = fmt.Sprintf("%s\n\t\t%s", lines[62], `"pg": "^8.10.0",`)
			lines[63] = fmt.Sprintf("%s\n\t\t%s", lines[63], `"pg-native": "3.0.1",`)
			lines[66] = fmt.Sprintf("\t\t%s", `"sequelize": "^6.6.5",`)
			lines[107] = fmt.Sprintf("\t\t%s", `"sequelize-cli": "^6.6.0",`)
		}
	}
	output := strings.Join(lines, "\n")

	if updatePackageDotJson {
		err := fileutils.WriteToFile(jsonFile, output)
		errorhandler.CheckNilErr(err)
	}

	return nil
}
