package pickyhelpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/hbs"
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
)

func (s StackDetails) CreateDockerFiles() error {

	var (
		path      string
		source    string
		fileFound bool
		err       error
	)

	if s.StackInfo[constants.WebStatus].(bool) {

		path = fmt.Sprintf("%s/%s/%s", utils.CurrentDirectory(),
			s.DirName,
			constants.DockerFile,
		)
		fileFound, _ = utils.IsExists(path)
		if !fileFound {
			source = `FROM node:14-alpine as baseimage
RUN mkdir app/
ADD . app/
WORKDIR /app

RUN npm install

FROM baseimage
CMD {{{cmdDockerfile stack}}}
EXPOSE 3000`

			err = hbs.ParseAndWriteToFile(source, path, s.StackInfo)
			errorhandler.CheckNilErr(err)
		}
		path = fmt.Sprintf("%s/%s/%s", utils.CurrentDirectory(),
			s.DirName,
			constants.DockerEnvFile,
		)
		fileFound, _ = utils.IsExists(path)
		if !fileFound {
			source = `GITHUB_URL=https://api.github.com/`

			err = utils.WriteToFile(path, source)
			errorhandler.CheckNilErr(err)
		}
		path = fmt.Sprintf("%s/%s/%s", utils.CurrentDirectory(),
			s.DirName,
			constants.DockerIgnoreFile,
		)
		fileFound, _ = utils.IsExists(path)
		if !fileFound {
			source = "node_modules\n.git\nbadges"
			err = utils.WriteToFile(path, source)
			errorhandler.CheckNilErr(err)
		}
	}

	// Add mobile related files.
	// if stackInfo[constants.MobileStatus].(bool) {}

	if s.StackInfo[constants.BackendStatus].(bool) {

		switch s.StackInfo[constants.Stack] {
		case constants.NodeExpressGraphqlTemplate, constants.NodeHapiTemplate:

			path = fmt.Sprintf("%s/%s/%s", utils.CurrentDirectory(),
				s.DirName,
				constants.DockerIgnoreFile,
			)
			fileFound, _ = utils.IsExists(path)
			if !fileFound {
				source = "node_modules\n.git\nbadges"
				err = utils.WriteToFile(path, source)
				errorhandler.CheckNilErr(err)
			}

			path = fmt.Sprintf("%s/%s/%s", utils.CurrentDirectory(),
				s.DirName,
				constants.DockerFile,
			)
			fileFound, _ = utils.IsExists(path)
			backendPortNumber := utils.FetchExistingPortNumber(s.DirName, constants.BackendPort)
			if fileFound {
				source = fmt.Sprintf(`FROM node:16
ARG ENVIRONMENT_NAME
ARG BUILD_NAME
RUN mkdir -p /app-build
ADD . /app-build
WORKDIR /app-build
RUN --mount=type=cache,target=/root/.yarn YARN_CACHE_FOLDER=/root/.yarn yarn --frozen-lockfile
RUN yarn
RUN yarn {{runBuildEnvironment stack}}

FROM node:16-alpine
ARG ENVIRONMENT_NAME
ARG BUILD_NAME
RUN mkdir -p /dist
RUN apk add yarn
RUN yarn global add {{globalAddDependencies database}}
RUN yarn add {{addDependencies database}}
ADD scripts/migrate-and-run.sh /
ADD package.json /
ADD . /
COPY --from=0 /app-build/dist ./dist

CMD ["sh", "./migrate-and-run.sh"]
EXPOSE %s`, backendPortNumber)

				err = hbs.ParseAndWriteToFile(source, path, s.StackInfo)
				errorhandler.CheckNilErr(err)
			}
		}
	}
	return nil
}
