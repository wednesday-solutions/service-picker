package sources

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
)

func CDBackendSource(stack, stackDir, environment string) string {
	userInput, singleQuote := utils.FindUserInputStackName(stackDir), "`"
	source := fmt.Sprintf(`# CD pipeline for %s for %s branch

name: %s CD -- %s

on:
  push:
    branches:
      - develop

jobs:
  docker-build-and-push:
    name: Docker build image and push
    runs-on: ubuntu-latest
    steps:
      # Checkout
      - name: Checkout to branch
        uses: actions/checkout@v3

      - name: Get branch name
        id: vars
        run: echo ::set-output name=short_ref::${GITHUB_REF_NAME}

      # Configure AWS with credentials
      - name: Configure AWS Credentials
        uses: aws-actions/configure-aws-credentials@v1
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: ${{ secrets.AWS_REGION }}

      # Login to Amazon ECR
      - name: Login to Amazon ECR
        id: login-ecr
        uses: aws-actions/amazon-ecr-login@v1

      # Build, tag, and push image to Amazon ECR
      - name: Build, tag, and push image to ECR
        env:
          ECR_REGISTRY: ${{ steps.login-ecr.outputs.registry }}
          ECR_REPOSITORY: ${{ secrets.AWS_ECR_REPOSITORY }}
          AWS_REGION: ${{ secrets.AWS_REGION }}
          IMAGE_TAG: ${{ github.sha }}
        working-directory: ./%s
        run: |
          docker compose build

      # Create and configure Amazon ECS task definition
      - name: Render Amazon ECS task definition
        id: %s-container
        uses: aws-actions/amazon-ecs-render-task-definition@v1
        with:
          task-definition: %s/task-definition-${{ steps.vars.outputs.short_ref }}.json
          container-name: %s-container-${{ steps.vars.outputs.short_ref }}
          image: ${{ steps.login-ecr.outputs.registry }}/${{ secrets.AWS_ECR_REPOSITORY }}:${{ github.sha }}

      # Deploy to Amazon ECS
      - name: Deploy to Amazon ECS
        uses: aws-actions/amazon-ecs-deploy-task-definition@v1
        with:
          task-definition: ${{ steps.%s-container.outputs.task-definition }}
          service: %s-service-${{ steps.vars.outputs.short_ref }}
          cluster: %s-cluster-${{ steps.vars.outputs.short_ref }}

      # Logout of Amazon
      - name: Logout of Amazon ECR
        if: always()
        run: docker logout ${{ steps.login-ecr.outputs.registry }}

      # Set %sBRANCH%s variable
      - name: Set env BRANCH
        run: echo "BRANCH=$(echo $GITHUB_REF | cur -d'/' -f 3)" >> $GITHUB_ENV

      # Get the current %senvironment%s
      - name: Get %senvironment_name%s
        id: env_vars
        run: |
          if [[ $BRANCH == 'master' ]]; then
            echo ::set-output name=environment_name::production
          elif [[ $BRANCH == 'qa' ]]; then
            echo ::set-output name=environment_name::qa
          else
            echo ::set-output name=environment_name::development
          fi
`,
		stackDir, environment, stackDir, environment, stackDir, userInput, stackDir,
		userInput, userInput, userInput, userInput, singleQuote, singleQuote, singleQuote,
		singleQuote, singleQuote, singleQuote,
	)
	return source
}

func TaskDefinitionSource(environment string) string {

	type logDriverOptionsType struct {
		awsLogsGroup        string
		awsLogsStreamPrefix string
		awsLogsRegion       string
	}

	type output struct {
		taskRole          string
		image             string
		containerName     string
		containerPort     string
		executionRole     string
		taskDefinitionArn string
		logDriver         string
		logDriverOptions  logDriverOptionsType
		family            string
		awsRegion         string
	}

	file := fmt.Sprintf(
		"%s/%s/%s", utils.CurrentDirectory(), ".sst", "outputs.json",
	)
	status, _ := utils.IsExists(file)
	if !status {
		return ""
	}
	sstOutputFile, err := os.Open(file)
	errorhandler.CheckNilErr(err)

	fileContent, err := ioutil.ReadAll(sstOutputFile)
	errorhandler.CheckNilErr(err)

	var data interface{}
	err = json.Unmarshal(fileContent, &data)
	errorhandler.CheckNilErr(err)

	jsonData, ok := data.(map[string]interface{})
	if !ok {
		errorhandler.CheckNilErr(fmt.Errorf("Something error happened when converting map"))
	}

	for key, value := range jsonData {

		if utils.EndsWith(key, "Pg") || utils.EndsWith(key, "Mysql") {

			var backendObj output
			var envName string
			if environment == constants.Development {
				envName = constants.Develop
			}
			backendJson, ok := value.(map[string]interface{})
			if !ok {
				errorhandler.CheckNilErr(fmt.Errorf("Something error happened when converting map"))
			}
			backendObj.taskRole, _ = backendJson["taskRole"].(string)
			backendObj.image, _ = backendJson["image"].(string)
			backendObj.containerName, _ = backendJson["containerName"].(string)
			backendObj.containerPort, _ = backendJson["containerPort"].(string)
			backendObj.executionRole, _ = backendJson["executionRole"].(string)
			backendObj.taskDefinitionArn, _ = backendJson["taskDefinitionArn"].(string)
			backendObj.logDriver, _ = backendJson["logDriver"].(string)
			backendObj.family, _ = backendJson["family"].(string)
			backendObj.awsRegion, _ = backendJson["awsRegion"].(string)

			logdriveroptions, _ := backendJson["logDriverOptions"].(map[string]interface{})
			backendObj.logDriverOptions.awsLogsGroup, _ = logdriveroptions["awslogs-group"].(string)
			backendObj.logDriverOptions.awsLogsStreamPrefix, _ = logdriveroptions["awslogs-stream-prefix"].(string)
			backendObj.logDriverOptions.awsLogsRegion, _ = logdriveroptions["awslogs-region"].(string)

			source := fmt.Sprintf(`{
  "ipcMode": null,
  "executionRoleArn": "%s",
  "containerDefinitions": [
    {
      "dnsSearchDomains": null,
      "logConfiguration": {
        "logDriver": "%s",
        "secretOptions": null,
        "options": {
          "awslogs-group": "%s",
          "awslogs-stream-prefix": "%s",
          "awslogs-region": "%s"
        }
      },
      "entryPoint": null,
      "portMappings": [
        {
          "hostPort": "9000",
          "protocol": "tcp",
          "containerPort": "%s"
        }
      ],
      "command": null,
      "linuxParameters": null,
      "cpu": 0,
      "environment": [
        {
          "name": "BUILD_NAME",
          "value": "%s"
        },
        {
          "name": "ENVIRONMENT_NAME",
          "value": ".%s"
        }
      ],
      "resourceRequirements": null,
      "ulimits": null,
      "dnsServers": null,
      "mountPoints": [],
      "workingDirectory": null,
      "secrets": null,
      "dockerSecurityOptions": null,
      "memory": null,
      "memoryReservation": null,
      "volumesFrom": [],
      "stopTimeout": null,
      "startTimeout": null,
      "firelensConfiguration": null,
      "dependsOn": null,
      "disableNetworking": null,
      "interactive": null,
      "healthCheck": null,
      "essential": true,
      "links": null,
      "hostname": null,
      "extraHosts": null,
      "pseudoTerminal": null,
      "user": null,
      "readonlyRootFilesystem": null,
      "dockerLabels": null,
      "systemControls": null,
      "privileged": null,
      "image": "%s",
      "name": "%s"
    }
  ],
  "placementConstraints": [],
  "memory": "2048",
  "taskRoleArn": "%s",
  "compatibilities": ["EC2", "FARGATE"],
  "taskDefinitionArn": "%s",
  "family": "%s",
  "requiresAttributes": [
    {
      "targetId": null,
      "targetType": null,
      "value": null,
      "name": "com.amazonaws.ecs.capability.logging-driver.awslogs"
    },
    {
      "targetId": null,
      "targetType": null,
      "value": null,
      "name": "ecs.capability.execution-role-awslogs"
    },
    {
      "targetId": null,
      "targetType": null,
      "value": null,
      "name": "com.amazonaws.ecs.capability.ecr-auth"
    },
    {
      "targetId": null,
      "targetType": null,
      "value": null,
      "name": "com.amazonaws.ecs.capability.docker-remote-api.1.19"
    },
    {
      "targetId": null,
      "targetType": null,
      "value": null,
      "name": "ecs.capability.execution-role-ecr-pull"
    },
    {
      "targetId": null,
      "targetType": null,
      "value": null,
      "name": "com.amazonaws.ecs.capability.docker-remote-api.1.18"
    },
    {
      "targetId": null,
      "targetType": null,
      "value": null,
      "name": "ecs.capability.task-eni"
    }
  ],
  "pidMode": null,
  "requiresCompatibilities": ["FARGATE"],
  "networkMode": "awsvpc",
  "cpu": "1024",
  "revision": 31,
  "status": "ACTIVE",
  "inferenceAccelerators": null,
  "proxyConfiguration": null,
  "volumes": []
}
`,
				backendObj.executionRole,
				backendObj.logDriver,
				backendObj.logDriverOptions.awsLogsGroup,
				backendObj.logDriverOptions.awsLogsStreamPrefix,
				backendObj.logDriverOptions.awsLogsRegion,
				backendObj.containerPort,
				envName,
				environment,
				backendObj.image,
				backendObj.containerName,
				backendObj.taskRole,
				backendObj.taskDefinitionArn,
				backendObj.family,
			)
			return source
		}
	}
	return ""
}
