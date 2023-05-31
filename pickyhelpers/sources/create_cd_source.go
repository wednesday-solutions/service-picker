package sources

import (
	"fmt"

	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/utils"
)

func CDBackendSource(stack, stackDir, environment string) string {

	userInput := utils.FindUserInputStackName(stackDir)
	source := fmt.Sprintf(`# CD pipeline for %s for %s branch

name: CD %s - %s

on:
  push:
    branches:
      - dev
      - qa
      - master
    # paths: "%s/**"
  workflow_dispatch:

jobs:
  docker-build-and-push:
    name: Docker build image and push
    runs-on: ubuntu-latest
    defaults:
      run: 
        working-directory: ./%s
    strategy:
      matrix:
        node-version: [16.14.x]

    steps:
      # Checkout
      - name: Checkout to branch
        uses: actions/checkout@v2

      - name: Use Node.js ${{ matrix.node-version }}
        uses: actions/setup-node@v2
        with:
          node-version: ${{ matrix.node-version }}

      - name: Get branch name
        id: vars
        run: echo ::set-output name=short_ref::${GITHUB_REF_NAME}

      - name: Set env.ENV_NAME and env.BUILD_NAME
        run: |
          if [[ ${{ steps.vars.outputs.short_ref }} == master ]]; then
              echo "BUILD_NAME=prod" >> "$GITHUB_ENV"
          else
              echo "ENV_NAME=.development" >> "$GITHUB_ENV"
              echo "BUILD_NAME=dev" >> "$GITHUB_ENV"
          fi

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
          DOCKER_BUILDKIT: 1
        run: |
          docker build --no-cache -t $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG . --build-arg BUILD_NAME=${{ env.BUILD_NAME }} --build-arg ENVIRONMENT_NAME=${{ env.ENV_NAME }}
          docker push $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG

      # Create and configure Amazon ECS task definition
      - name: Render Amazon ECS task definition
        id: %s-container
        uses: aws-actions/amazon-ecs-render-task-definition@v1
        with:
          task-definition: %s/task-definition-${{ env.BUILD_NAME }}.json
          container-name: %s-container-${{ env.BUILD_NAME }}
          image: ${{ steps.login-ecr.outputs.registry }}/${{ secrets.AWS_ECR_REPOSITORY }}:${{ github.sha }}

      # Deploy to Amazon ECS
      - name: Deploy to Amazon ECS
        uses: aws-actions/amazon-ecs-deploy-task-definition@v1
        with:
          task-definition: ${{ steps.%s-container.outputs.task-definition }}
          service: %s-service-${{ env.BUILD_NAME }}
          cluster: %s-cluster-${{ env.BUILD_NAME }}

      # Logout of Amazon
      - name: Logout of Amazon ECR
        if: always()
        run: docker logout ${{ steps.login-ecr.outputs.registry }}
`,
		stackDir, environment, stackDir, environment, stackDir, stackDir,
		userInput, stackDir, userInput, userInput, userInput, userInput,
	)
	return source
}

func TaskDefinitionSource(environment, stackDir string) string {

	taskDefinition := utils.GetOutputsBackendObject(environment, stackDir)

	taskDefinitionSource := fmt.Sprintf(`{
  "taskRoleArn": "%s",
  "executionRoleArn": "%s",
  "taskDefinitionArn": "%s",
  "family": "%s",
  "containerDefinitions": [
    {
      "name": "%s",
      "image": "%s",
      "logConfiguration": {
        "logDriver": "%s",
        "secretOptions": null,
        "options": {
          "awslogs-group": "%s",
          "awslogs-stream-prefix": "%s",
          "awslogs-region": "%s"
        }
      },
      "portMappings": [
        {
          "hostPort": "9000",
          "protocol": "tcp",
          "containerPort": "%s"
        }
      ],
      "environment": [
        {
          "name": "BUILD_NAME",
          "value": "%s"
        },
        {
          "name": "ENVIRONMENT_NAME",
          "value": "%s"
        }
      ],
      "secrets": [
        {
          "name": "%s",
          "valueFrom": "%s:password::"
        }
      ],
      "cpu": 0,
      "memory": null,
      "command": null,
      "entryPoint": null,
      "dnsSearchDomains": null,
      "linuxParameters": null,
      "resourceRequirements": null,
      "ulimits": null,
      "dnsServers": null,
      "mountPoints": [],
      "workingDirectory": null,
      "dockerSecurityOptions": null,
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
      "privileged": null
    }
  ],
  "placementConstraints": [],
  "memory": "2048",
  "compatibilities": ["EC2", "FARGATE"],
  "requiresAttributes": [
    {
      "name": "com.amazonaws.ecs.capability.logging-driver.awslogs"
    },
    {
      "name": "ecs.capability.execution-role-awslogs"
    },
    {
      "name": "com.amazonaws.ecs.capability.ecr-auth"
    },
    {
      "name": "com.amazonaws.ecs.capability.docker-remote-api.1.19"
    },
    {
      "name": "ecs.capability.execution-role-ecr-pull"
    },
    {
      "name": "com.amazonaws.ecs.capability.docker-remote-api.1.18"
    },
    {
      "name": "ecs.capability.task-eni",
      "targetId": null,
      "targetType": null,
      "value": null
    }
  ],
  "ipcMode": null,
  "pidMode": null,
  "requiresCompatibilities": ["FARGATE"],
  "networkMode": "awsvpc",
  "cpu": "1024",
  "revision": 17,
  "status": "ACTIVE",
  "inferenceAccelerators": null,
  "proxyConfiguration": null,
  "volumes": [],
  "tags": [
    {
      "key": "sst:app",
      "value": "web-app"
    },
    {
      "key": "sst:stage",
      "value": "dev"
    }
  ]
}
`,
		taskDefinition.BackendObj.TaskRole,
		taskDefinition.BackendObj.ExecutionRole,
		taskDefinition.BackendObj.TaskDefinition,
		taskDefinition.BackendObj.Family,
		taskDefinition.BackendObj.ContainerName,
		taskDefinition.BackendObj.Image,
		taskDefinition.BackendObj.LogDriver,
		taskDefinition.BackendObj.LogDriverOptions.AwsLogsGroup,
		taskDefinition.BackendObj.LogDriverOptions.AwsLogsStreamPrefix,
		taskDefinition.BackendObj.LogDriverOptions.AwsLogsRegion,
		taskDefinition.BackendObj.ContainerPort,
		taskDefinition.Environment,
		taskDefinition.EnvName,
		taskDefinition.SecretName,
		taskDefinition.BackendObj.SecretArn,
	)
	return taskDefinitionSource
}

func CDWebSource(stack, dirName string) string {
	var sourceDir string
	if stack == constants.ReactJS {
		sourceDir = "build"
	} else if stack == constants.NextJS {
		sourceDir = "out"
	}
	source := fmt.Sprintf(`name: CD %s
on:
  push:
    branches:
      - dev
      - qa
      - master
    # paths: "%s/**"
  workflow_dispatch:

jobs:
  deploy:
    name: Deploy
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/master' || github.ref == 'refs/heads/develop' || github.ref == 'refs/heads/qa'
    defaults:
      run:
        working-directory: ./%s
    strategy:
      matrix:
        node-version: [16.13.0]
    env:
      SOURCE_DIR: "./%s/%s/"
      AWS_REGION: ${{ secrets.AWS_REGION }}
      AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
      AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
      PATHS: "/*"

    steps:
      - name: Checkout to branch
        uses: actions/checkout@v2

      - name: Use Node.js ${{ matrix.node-version }}
        uses: actions/setup-node@v2
        with:
          node-version: ${{ matrix.node-version }}

      - name: Get branch name
        id: vars
        run: echo ::set-output name=short_ref::${GITHUB_REF_NAME}

      - name: Set short environment name
        id: environment
        run: |
          if [[ ${{ steps.vars.outputs.short_ref }} == master ]]; then
               echo ::set-output name=short_env::prod
          elif [[ ${{ steps.vars.outputs.short_ref }} == qa ]]; then
               echo ::set-output name=short_env::qa
          else
               echo ::set-output name=short_env::dev
          fi

      - name: Install dependencies
        run: yarn

      - name: Build
        run: yarn build:${{ steps.environment.outputs.short_env }}

      - name: AWS Deploy to S3
        uses: jakejarvis/s3-sync-action@v0.5.1
        with:
          args: --follow-symlinks --delete
        env:
          AWS_S3_BUCKET: %s-${{ steps.environment.outputs.short_env }}

      - name: Invalidate CloudFront
        uses: chetan/invalidate-cloudfront-action@v2.4
        env:
          DISTRIBUTION: ${{ secrets.DISTRIBUTION_ID }}
`,
		dirName, dirName, dirName, dirName, sourceDir, dirName)
	return source
}
