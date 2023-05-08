<img align="left" src="" width="480" height="620" />

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

&nbsp; [![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)

<br><br>

## Table of contents

1. [Overview](#overview)
2. [Tech Stacks](#tech-stacks)
3. [Installing](#installing)
4. [Creating an App](#creating-an-app)
5. [Feedback](#feedback)
6. [License](#license)

## Overview

Service-Picker is a simple library for creating a new full stack project and deploy it in aws by just doing some cli stuffs. We are using `picky` as our CLI tool.

Keep the repos you're interested in and you're good to go.

`picky` works on macOS, Windows, Linux.<br>
If something doesn't work, please file an [issue](https://github.com/wednesday-solutions/service-picker/issues).<br>
If you have questions, suggestions or need help, please ask in [GitHub Discussions](https://github.com/wednesday-solutions/service-picker/discussions)

## Tech Stacks

This tool will have support for production applications using the following tech stacks.

**Web:**

- [React JS](https://gitub.com/wednesday-solutions/react-template)
- [Next JS](https://github.com/wednesday-solutions/nextjs-template)

**Mobile:**

- [Android App](https://github.com/wednesday-solutions/android-template)
- [iOS App](https://github.com/wednesday-solutions/ios-template)
- [React Native App](https://github.com/wednesday-solutions/react-native-template)
- [Flutter App](https://github.com/wednesday-solutions/flutter_template)

**Backend:**

- [Node (Hapi - REST API)](https://github.com/wednesday-solutions/nodejs-hapi-template)
- [Node (Express - GraphQL API)](https://github.com/wednesday-solutions/node-express-graphql-template)
- [Node (Express - REST API)](https://github.com/wednesday-solutions/node-mongo-express)
- [Golang (Echo - GraphQL API)](https://github.com/wednesday-solutions/go-template)

**Databases:**

- MySQL
- PostgreSQL
- MongoDB
- DynamoDB
- Neo4j

**Infrastructure:**

- Redis
- Kafka

## Installing

Using Picky is easy. First use `go get` to install the latest version of the library (`go` should be installed in your system).

```bash
go install github.com/wednesday-solutions/picky@latest
```

Please make sure the installation is successful by running the following command.

```bash
picky -v
```

## Creating an App

To create a new app, you need to pick stacks which are mentioned in [tech stacks](#tect-stacks)
To start using `picky`, first, create a new directory and run the following command.

```bash
picky service
```

<img width="469" alt="Pick a service" src="https://user-images.githubusercontent.com/114065489/236760233-e3dadf7a-42de-4f98-8cba-7c01161b1d3c.png">

Use the arrow keys to navigate and pick a service you want.<br>
If you select `web`, the following will come up,

<img width="504" alt="pick-stack" src="https://user-images.githubusercontent.com/114065489/236762803-3e8d6b67-bcf8-4ff3-a70c-a43424ca1457.png">

The complete tutorial is given below.

https://user-images.githubusercontent.com/114065489/236762965-ff6b9dab-e357-4c17-b57f-1ec09cdfb440.mp4

## Feedback

If you have any feedback, please reach out to us at [GitHub Discussions](https://github.com/wednesday-solutions/service-picker/discussions)

## License

This project is under the [MIT License](https://github.com/wednesday-solutions/service-picker).
