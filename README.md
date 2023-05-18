<img align="left" src="https://github.com/wednesday-solutions/service-picker/assets/114065489/c3366e7d-7213-4ed4-9cc5-b19c0271c03b"  width="380" height="500" />

<div>
  <a href="https://www.wednesday.is?utm_source=gthb&utm_medium=repo&utm_campaign=serverless" align="left" style="margin-left: 0;">
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
    <a href="https://www.wednesday.is/contact-us?utm_source=gthb&utm_medium=repo&utm_campaign=serverless" target="_blank">
      <img src="https://uploads-ssl.webflow.com/5ee36ce1473112550f1e1739/5f6ae88b9005f9ed382fb2a5_button_get_in_touch.svg" width="121" height="34">
    </a>
    <a href="https://github.com/wednesday-solutions/" target="_blank">
      <img src="https://uploads-ssl.webflow.com/5ee36ce1473112550f1e1739/5f6ae88bb1958c3253756c39_button_follow_on_github.svg" width="168" height="34">
    </a>
  </div>

---

<span>We're always looking for people who value their work, so come and join us. <a href="https://www.wednesday.is/hiring">We are hiring!</a></span>

</div>

<br>

## Table of contents

1. [Overview](#overview)
2. [Tech Stacks](#tech-stacks)
3. [Prerequisites](#pre-requisites)
4. [Installation](#installation)
5. [Creating a Project](#creating-a-project)
6. [User Guide](#user-guide)
7. [Project Structure](#project-structure)
8. [Feedback](#feedback)
9. [License](#license)
10. [Future Plans](#future-plans)

## Overview

Once business gives the sign-off and it's time for execution, the question thats most frequently asked is "What's the tech stack that we should use for this?"

Fast-forward past the long debates where engineers are pitching their favourite languages and frameworks, we kick "git init" a bunch of repos for the frontend, infrastructure & backend. We then spend some time creating some boilerplate, or use some templates, setting up the CI/CD pipelines, IaC etc. The next thing you know, it's 6 weeks later and you're still configuring the connection between your database and your ec2 instance. The amount of business logic that you've written is minimal cause all of your team's time was spent configuring repos, environments, security groups, and other nitty-grittys.

Thats where the service-picker comes in. We're working on building a cli tool that allows you to scaffold batteries included code, IaC & CI/CD pipelines, across environments and stacks, that is completely integrated with each other and ready to go.

This means that setting up the infra and codebase for your next project which needs a React web app, with a node.js backend and a postgreSQL db is as simple as a hitting a few arrow buttons, and enter a couple of times.

`picky` is our CLI tool. `picky` works on macOS, Windows and Linux.<br>
If something doesn't work, please file an [issue](https://github.com/wednesday-solutions/service-picker/issues).<br>
If you have questions, suggestions or need help, please ask in [GitHub Discussions](https://github.com/wednesday-solutions/service-picker/discussions)

## Tech Stacks

This tool will have support for production applications using the following tech stacks.

**Web:**

- [React JS](https://gitub.com/wednesday-solutions/react-template)
- [Next JS](https://github.com/wednesday-solutions/nextjs-template)

**Backend:**

- [Node (Hapi - REST API)](https://github.com/wednesday-solutions/nodejs-hapi-template)
- [Node (Express - GraphQL API)](https://github.com/wednesday-solutions/node-express-graphql-template)

**Databases:**

- MySQL
- PostgreSQL

**Cache:**

- Redis

## Setup and Configuration.

---

### Pre-requisites

- [Golang](https://go.dev/doc/install)
- [Node JS](https://nodejs.org/en/download)
- [Docker](https://docs.docker.com/engine/install/)
- Package Manager([npm](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm) or [yarn](https://classic.yarnpkg.com/lang/en/docs/install/#mac-stable))

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

<img width="469" alt="Pick a service" src="https://user-images.githubusercontent.com/114065489/236760233-e3dadf7a-42de-4f98-8cba-7c01161b1d3c.png">

Use the arrow keys to navigate and pick a service you want.<br>

The complete stack initialization tutorial is given below.

![stack_initialisation_demo](https://github.com/wednesday-solutions/service-picker/assets/114065489/79c66f38-f35f-451d-889b-01499c977b6e)

You can see `picky`'s home page if you initialized atleast one stack. You can choose any option in the following.

<img width="551" alt="Picky Home Preview Image" src="https://github.com/wednesday-solutions/service-picker/assets/114065489/8b2a94dd-3103-4a1c-ad4e-1eedd7e87164">

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

Currently the service-picker is capable of setting up full-stack javascipt applications. In it's end state the service picker will allow you to chose right from your cloud infra structure provider (GCP, AWS, AZURE) to different backends and databases that you'd like to use, to your caching strategy, message broker, mobile app release tooling and any other tooling choice that you must make along the product development journey.
