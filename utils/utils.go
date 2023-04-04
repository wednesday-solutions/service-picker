package utils

import (
	"fmt"
	"strings"

	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func DirectoryName(dirName, service string) string {
	switch service {
	case constants.Web:
		return fmt.Sprintf("%s-%s", dirName, constants.Web)
	case constants.Mobile:
		return fmt.Sprintf("%s-%s", dirName, constants.Mobile)
	case constants.Backend:
		return fmt.Sprintf("%s-%s", dirName, constants.Backend)
	default:
		return ""
	}
}

func ServiceExist(service string) (string, bool) {
	var splitDirName []string
	var existingService string
	dirNames, err := fileutils.ReadAllContents(fileutils.CurrentDirectory())
	errorhandler.CheckNilErr(err)
	for _, dirName := range dirNames {
		splitDirName = strings.Split(dirName, "-")
		if len(splitDirName) > 0 {
			existingService = splitDirName[len(splitDirName)-1]
			if existingService == service {
				return dirName, true
			}
		}
	}
	return "", false
}

func ServicesExist() (map[string]bool, map[string]string) {
	var splitDirName []string
	var service string
	serviceStatuses := make(map[string]bool)
	serviceDirectories := make(map[string]string)
	dirNames, err := fileutils.ReadAllContents(fileutils.CurrentDirectory())
	errorhandler.CheckNilErr(err)
	for _, dirName := range dirNames {
		splitDirName = strings.Split(dirName, "-")
		if len(splitDirName) > 0 {
			service = splitDirName[len(splitDirName)-1]
			switch service {
			case constants.Web:
				serviceStatuses[constants.WebStatus] = true
				serviceDirectories[constants.Web] = dirName
			case constants.Mobile:
				serviceStatuses[constants.MobileStatus] = true
				serviceDirectories[constants.Mobile] = dirName
			case constants.Backend:
				serviceStatuses[constants.BackendStatus] = true
				serviceDirectories[constants.Backend] = dirName
			}
		}
	}
	return serviceStatuses, serviceDirectories
}
