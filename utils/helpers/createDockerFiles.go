package helpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
	"github.com/wednesday-solutions/picky/utils/hbs"
)

func CreateDockerFiles(stackData map[string]interface{}) error {

	var (
		path      string
		source    string
		fileFound bool
		err       error
	)
	dockerfile := "Dockerfile"
	dockerEnvFile := ".env.docker"
	dockerIgnoreFile := ".dockerignore"

	if stackData["webStatus"].(bool) {

		path = fmt.Sprintf("%s/%s/%s", fileutils.CurrentDirectory(),
			constants.Web,
			dockerfile,
		)
		fileFound, _ = fileutils.IsExists(path)
		if !fileFound {
			source = `FROM node:14-alpine as baseimage
RUN mkdir app/
ADD . app/
WORKDIR /app

RUN npm install

FROM baseimage
CMD ["yarn", "start"]
EXPOSE 3000`

			err = fileutils.WriteToFile(path, source)
			errorhandler.CheckNilErr(err)
		}
		path = fmt.Sprintf("%s/%s/%s", fileutils.CurrentDirectory(),
			constants.Web,
			dockerEnvFile,
		)
		fileFound, _ = fileutils.IsExists(path)
		if !fileFound {
			source = `GITHUB_URL=https://api.github.com/`

			err = fileutils.WriteToFile(path, source)
			errorhandler.CheckNilErr(err)
		}
		path = fmt.Sprintf("%s/%s/%s", fileutils.CurrentDirectory(),
			constants.Web,
			dockerIgnoreFile,
		)
		fileFound, _ = fileutils.IsExists(path)
		if !fileFound {
			source = "node_modules\n.git\nbadges"
			err = fileutils.WriteToFile(path, source)
			errorhandler.CheckNilErr(err)
		}
	}

	// Add mobile related files.
	// if stackData["mobileStatus"].(bool) {}

	if stackData["backendStatus"].(bool) {

		path = fmt.Sprintf("%s/%s/%s", fileutils.CurrentDirectory(),
			constants.Web,
			dockerIgnoreFile,
		)
		fileFound, _ = fileutils.IsExists(path)
		if !fileFound {
			source = "node_modules\n.git\nbadges"
			err = fileutils.WriteToFile(path, source)
			errorhandler.CheckNilErr(err)
		}

		path = fmt.Sprintf("%s/%s/%s", fileutils.CurrentDirectory(),
			constants.Backend,
			dockerfile,
		)
		fileFound, _ = fileutils.IsExists(path)
		if fileFound {
			source = `FROM node:14
ARG ENVIRONMENT_NAME
RUN mkdir -p /app-build
ADD . /app-build
WORKDIR /app-build
RUN --mount=type=cache,target=/root/.yarn YARN_CACHE_FOLDER=/root/.yarn yarn --frozen-lockfile
RUN yarn
RUN yarn {{runBuildEnvironment stack}}

FROM node:14-alpine
ARG ENVIRONMENT_NAME
ENV ENVIRONMENT_NAME $ENVIRONMENT_NAME
RUN mkdir -p /dist
RUN apk add yarn
RUN yarn global add {{globalAddDependencies database}}
RUN yarn add {{addDependencies database}}
ADD scripts/migrate-and-run.sh /
ADD package.json /
ADD . /
COPY --from=0 /app-build/dist ./dist

CMD ["sh", "./migrate-and-run.sh"]
EXPOSE 9000`

			err = hbs.ParseAndWriteToFile(source, path, stackData)
			errorhandler.CheckNilErr(err)
		}
	}
	return nil
}
