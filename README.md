<img align="left" src="https://github.com/wednesday-solutions/service-picker/assets/114065489/fb48bfef-29d6-4568-b2da-acb3ceff924d"  width="380" height="500" />

<div>
  <a href="https://www.wednesday.is/?utm_source=github&utm_medium=service-picker" align="left" style="margin-left: 0;">
    <img src="https://uploads-ssl.webflow.com/5ee36ce1473112550f1e1739/5f5879492fafecdb3e5b0e75_wednesday_logo.svg">
  </a>
  <p>
    <h1 align="left">Service Picker
    </h1>
  </p>

  <p>
A CLI tool that helps to setup full-stack javascript applications without having to touch any code. You'll be able to pick templates and databases of your choice, integrate it, set up automation pipelines and create infrastructure in AWS with ease.
It contains a number of <a href="https://github.com/wednesday-solutions">Wednesday Solutions</a>'s open source projects, connected and working together. Pick whatever you need and build your own ecosystem.
  </p>

---

  <p>
    <h4>
      Expert teams of digital product strategists, developers, and designers.
    </h4>
  </p>

  <div>
    <a href="https://www.wednesday.is/contact-us/?utm_source=github&utm_medium=service-picker" target="_blank">
      <img src="https://uploads-ssl.webflow.com/5ee36ce1473112550f1e1739/5f6ae88b9005f9ed382fb2a5_button_get_in_touch.svg" width="121" height="34">
    </a>
    <a href="https://github.com/wednesday-solutions/" target="_blank">
      <img src="https://uploads-ssl.webflow.com/5ee36ce1473112550f1e1739/5f6ae88bb1958c3253756c39_button_follow_on_github.svg" width="168" height="34">
    </a>
  </div>

---

<span>We're always looking for people who value their work, so come and join us. <a href="https://www.wednesday.is/hiring/?utm_source=github&utm_medium=service-picker">We are hiring!</a></span>

</div>

<br>

## Table of contents

- [Table of contents](#table-of-contents)
- [Overview](#overview)
- [Tech Stacks](#tech-stacks)
- [Setup and Configuration.](#setup-and-configuration)
  - [Pre-requisites](#pre-requisites)
  - [Installation](#installation)
- [Creating a Project](#creating-a-project)
- [User Guide](#user-guide)
- [Project Structure](#project-structure)
- [Feedback](#feedback)
- [License](#license)
- [Future Plans](#future-plans)

## Overview

Once business gives the sign-off and it's time for execution, the question thats most frequently asked is "What's the tech stack that we should use for this?"

Fast-forward past the long debates where engineers are pitching their favourite languages and frameworks, we kick "git init" a bunch of repos for the frontend, infrastructure & backend. We then spend some time creating some boilerplate, or use some templates, setting up the CI/CD pipelines, IaC etc. The next thing you know, it's 6 weeks later and you're still configuring the connection between your database and your ec2 instance. The amount of business logic that you've written is minimal cause all of your team's time was spent configuring repos, environments, security groups, and other nitty-grittys.

Thats where the service-picker comes in. We're working on building a cli tool that allows you to scaffold batteries included code, IaC & CI/CD pipelines, across environments and stacks, that is completely integrated with each other and ready to go.

This means that setting up the infra and codebase for your next project which needs a React web app, with a node.js backend and a postgreSQL db is as simple as a hitting a few arrow buttons, and enter a couple of times.

Service picker works on macOS, Windows and Linux.<br>
If something doesn't work, please file an [issue](https://github.com/wednesday-solutions/service-picker/issues).<br>
If you have questions, suggestions or need help, please ask in [GitHub Discussions](https://github.com/wednesday-solutions/service-picker/discussions)

## Tech Stacks

This tool will have support for production applications using the following tech stacks.

**Web:**

- [React JS](https://github.com/wednesday-solutions/react-template)
- [Next JS](https://github.com/wednesday-solutions/nextjs-template)

**Backend:**

- [Node (Hapi - REST API)](https://github.com/wednesday-solutions/nodejs-hapi-template)
- [Node (Express - GraphQL API)](https://github.com/wednesday-solutions/node-express-graphql-template)

**Databases:**

- MySQL
- PostgreSQL

**Cache:**

- Redis

**Infrastructure Provider:**

- [AWS](https://aws.amazon.com/)

## Setup and Configuration.

### Pre-requisites

- [Golang](https://go.dev/doc/install)
- [Node JS](https://nodejs.org/en/download)
- Package Manager([npm](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm) or [yarn](https://classic.yarnpkg.com/lang/en/docs/install/#mac-stable))
- [Docker](https://docs.docker.com/engine/install/) - Install and have it running in your local to docker compose applications and setup infrastructures in AWS.
- [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html) - Configure to your AWS account.
```bash
 $ aws configure
AWS Access Key ID: MYACCESSKEYID
AWS Secret Access Key: MYSECRETKEY
Default region name [us-west-2]: MYAWSREGION
Default output format [None]:
```
- Create a repository in your [AWS ECR](https://aws.amazon.com/ecr/).

```bash
 $ aws ecr create-repository --repository-name cdk-hnb659fds-container-assets-MYAWSACCOUNTID-MYAWSREGION
```


### Installation

Using Picky is easy. First use `go install` to install the latest version of the library (`go` should be installed in your system).

```bash
go install github.com/wednesday-solutions/picky@latest
```

Please make sure the installation is successful by running the following command.

```bash
picky -v
```

## Creating a Project

To create a new project, you need to pick stacks which are mentioned in [tech stacks](#tech-stacks)<br>
To start using `picky`

```bash
mkdir my-project
cd my-project
```

```bash
picky service
```

![pick_a_service](https://github.com/wednesday-solutions/service-picker/assets/114065489/a3440e2a-9486-4419-85a6-eae2029ba45a)

Use the arrow keys to navigate and pick a service you want.<br>

The complete stack initialization tutorial is given below.

![stack_initialisation_demo](https://github.com/wednesday-solutions/service-picker/assets/114065489/0c2f0e0b-bc05-4a3d-9d27-a69c6654419b)

You can see `picky`'s home page if you initialized atleast one stack. You can choose any option in the following.

<img width="551" alt="Picky Home Preview Image" src="https://github.com/wednesday-solutions/service-picker/assets/114065489/322fc49b-43b3-4d0e-9753-6db6d3dd96a5">

***Tips:***
- If you want to go back from the prompt, click `Ctrl + D`
- If you want to exit from the prompt, click `Ctrl + C`

## User Guide

| Option           | Use                                                                                       |
| ---------------- | ----------------------------------------------------------------------------------------- |
| `Init Service`   | Initialize a stack.                                                                       |
| `CI/CD`          | Create CI/CD Pipeline in GitHub.                                                          |
| `Docker Compose` | Create Docker Compose file for the mono-repo. It consist of all the selected stacks.      |
| `Setup Infra`    | Setup infrastructure for initialized stacks.                                              |
| `Deploy`         | Deploy the infrastructure in AWS. It can deploy Frontend, Backend or Full stack projects. |
| `Remove Deploy`  | Remove the deployed infrastructure.                                                       |
| `Git Init`       | Initialize empty git repository in the current directory.                                 |
| `Exit`           | Exit from the tool.                                                                       |

## Project Structure

It will be like the following in the current directory.

```
my-project
├── .github
│   └── workflows
│       ├── cd-backend-node-hapi-pg.yml
│       ├── cd-frontend-next-js-web.yml
│       ├── ci-backend-node-hapi-pg.yml
│       └── ci-frontend-next-js-web.yml
├── .sst
│   ├── artifacts
│   ├── dist
│   ├── types
│   ├── debug.log
│   └── outputs.json
├── node_modules
├── stacks
│   ├── BackendNodeHapiPg.js
│   └── FrontendNextJsWeb.js
├── backend-node-hapi-pg
│   └── ...
├── frontend-next-js-web
│   └── ...
├── .env
├── .git
├── .gitignore
├── cdk.context.json
├── docker-compose.yml
├── backend-node-hapi-outputs.json
├── frontend-next-js-web-outputs.json
├── package.json
├── parseSstOutputs.js
├── sst.config.js
└── yarn.lock
```

## Feedback

If you have any feedback, please reach out to us at [GitHub Discussions](https://github.com/wednesday-solutions/service-picker/discussions)

## License

This project is under the [MIT License](https://github.com/wednesday-solutions/service-picker/blob/main/LICENSE).

## Future Plans

Currently the service-picker is capable of setting up full-stack javascript applications. In it's end state the service picker will allow you to choose right from your cloud infra structure provider (GCP, AWS, AZURE) to different backends and databases that you'd like to use, to your caching strategy, message broker, mobile app release tooling and any other tooling choice that you must make along your product development journey.
