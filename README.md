<img align="left" src="https://github.com/wednesday-solutions/service-picker/assets/114065489/c3366e7d-7213-4ed4-9cc5-b19c0271c03b"  width="450" height="585" />

<div>
  <a href="https://www.wednesday.is?utm_source=gthb&utm_medium=repo&utm_campaign=serverless" align="left" style="margin-left: 0;">
    <img src="https://uploads-ssl.webflow.com/5ee36ce1473112550f1e1739/5f5879492fafecdb3e5b0e75_wednesday_logo.svg">
  </a>
  <p>
    <h1 align="left">Service Picker
    </h1>
  </p>

  <p>
A CLI tool for creating real projects without touching any kind of code.
It contain a number of <a href="https://github.com/wednesday-solutions">Wednesday Solutions</a>'s open source projects, connected and working together. Pick whatever you need and build your own ecosystem.
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

<span>We’re always looking for people who value their work, so come and join us. <a href="https://www.wednesday.is/hiring">We are hiring!</a></span>

</div>

<br>

## Table of contents

1. [Overview](#overview)
2. [Tech Stacks](#tech-stacks)
3. [Installing](#installing)
4. [Creating a Project](#creating-a-project)
5. [User Guide](#user-guide)
6. [Project Structure](#project-structure)
7. [Feedback](#feedback)
8. [License](#license)
9. [Future Plans](#future-plans)

## Overview

Service-Picker is a simple library for creating a new full stack project and deploy it in aws by just doing some cli stuffs. We are using `picky` as our CLI tool.

Keep the repos you're interested in and you're good to go.

`picky` works on macOS, Windows and Linux.<br>
If something doesn't work, please file an [issue](https://github.com/wednesday-solutions/service-picker/issues).<br>
If you have questions, suggestions or need help, please ask in [GitHub Discussions](https://github.com/wednesday-solutions/service-picker/discussions)

## Tech Stacks

This tool will have support for production applications using the following tech stacks.

**Web:**

- [React JS](https://gitub.com/wednesday-solutions/react-template)
- [Next JS](https://github.com/wednesday-solutions/nextjs-template)

<!-- **Mobile:**

- [Android App](https://github.com/wednesday-solutions/android-template)
- [iOS App](https://github.com/wednesday-solutions/ios-template)
- [React Native App](https://github.com/wednesday-solutions/react-native-template)
- [Flutter App](https://github.com/wednesday-solutions/flutter_template) -->

**Backend:**

- [Node (Hapi - REST API)](https://github.com/wednesday-solutions/nodejs-hapi-template)
- [Node (Express - GraphQL API)](https://github.com/wednesday-solutions/node-express-graphql-template)
<!-- - [Node (Express - REST API)](https://github.com/wednesday-solutions/node-mongo-express)
- [Golang (Echo - GraphQL API)](https://github.com/wednesday-solutions/go-template) -->

**Databases:**

- MySQL
- PostgreSQL
<!-- - MongoDB
- DynamoDB
- Neo4j -->

**Infrastructure:**

- Redis
<!-- - Kafka -->

## Installing

Using Picky is easy. First use `go get` to install the latest version of the library (`go` should be installed in your system).

```bash
go install github.com/wednesday-solutions/picky@latest
```

Please make sure the installation is successful by running the following command.

```bash
picky -v
```

## Creating a Project

To create a new project, you need to pick stacks which are mentioned in [tech stacks](#tech-stacks)
To start using `picky`,

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

[<img src="https://uploads-ssl.webflow.com/5ee36ce1473112550f1e1739/5f5879492fafecdb3e5b0e75_wednesday_logo.svg">](https://user-images.githubusercontent.com/114065489/236762965-ff6b9dab-e357-4c17-b57f-1ec09cdfb440.mp4 "Stacks initialization tutorial")

You can see `picky`'s home page if you initialized atleast one stack. You can choose any option in the following.

<img width="461" alt="Picky Home" src="https://user-images.githubusercontent.com/114065489/236789527-d108e2ef-dc46-4115-b946-19511d00304e.png">

## User Guide

| Option          | Use                                                                                       |
| --------------- | ----------------------------------------------------------------------------------------- |
| `Init`          | Initialize a stack.                                                                       |
| `CI/CD`         | Create CI/CD Pipeline in GitHub.                                                          |
| `Setup Infra`   | Setup infrastructure for initialized stacks.                                              |
| `Deploy`        | Deploy the infrastructure in AWS. It can deploy Frontend, Backend or Full stack projects. |
| `Remove Deploy` | Remove the deployed infrastructure.                                                       |
| `Git Init`      | Initialize empty git repository in the current directory.                                 |
| `Exit`          | Exit from the tool.                                                                       |

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
├── package.json
├── parseSstOutputs.js
├── sst.config.js
└── yarn.lock
```

## Feedback

If you have any feedback, please reach out to us at [GitHub Discussions](https://github.com/wednesday-solutions/service-picker/discussions)

## License

This project is under the [MIT License](https://github.com/wednesday-solutions/service-picker).

## Future Plans

As of now, we can build full stack project which consist of the backend `Node JS` and the frontend `React JS` or `Next JS`. We don't have mobile application support. So We are planned to add support for mobile applications, currently we are working on it. And we have only `AWS` support now. Also we will add support of different cloud providers such as `GCP` and `Azure` and other stacks.
