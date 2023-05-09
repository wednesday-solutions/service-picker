package sources

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/wednesday-solutions/picky/internal/utils"
)

func CDSource(stack, environment string) string {
	singleQuote, projectName := "`", strcase.ToKebab(utils.GetProjectName())
	source := fmt.Sprintf(`# CD pipeline for %s for %s branch

name: %s CD -- %s

on:
  push:
    branches:
      - develop
      - qa

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
        env:
          CI: true

      # Configure AWS with credentials
      - name: Configure AWS Credentials
        uses: aws-actions/configure-aws-credentials@v1
        with:
          aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
          aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          aws-region: ${{ secrets.AWS_REGION }}

      # Populate existing .env with required envs (master only)
      - name: Append AWS ENVs in the .env file
        if: github.ref == 'refs/heads/master'
        run: |
          echo AWS_SECRET_ACCESS_KEY=${{ secrets.AWS_SECRET_ACCESS_KEY }} >> .env
          echo AWS_ACCESS_KEY_ID=${{ secrets.AWS_ACCESS_KEY_ID }} >> .env

      # Populate existing .env.qa with required envs (qa only)
      - name: Append AWS ENVs in the .env.qa file
        if: github.ref == 'refs/heads/qa'
        run: |
          echo AWS_SECRET_ACCESS_KEY=${{ secrets.AWS_SECRET_ACCESS_KEY }} >> .env.qa
          echo AWS_ACCESS_KEY_ID=${{ secrets.AWS_ACCESS_KEY_ID}} >> .env.qa

      # Populate existing .env.development with required envs (develop only)
      - name: Append AWS ENVs in the .env.development file
        if: github.ref == 'refs/heads/develop'
        run: |
          echo AWS_SECRET_ACCESS_KEY=${{ secrets.AWS_SECRET_ACCESS_KEY }} >> .env.development
          echo AWS_ACCESS_KEY_ID=${{ secrets.AWS_ACCESS_KEY_ID }} >> .env.development

      # Login to Amazon ECR
      - name: Login to Amazon ECR
        id: login-ecr
        uses: aws-actions/amazon-ecr-login@v1

      # Build, tag, and push image to Amazon ECR
      - name: Build, tag, and push image to ECR
        env:
          ECR_REGISTRY: ${{ steps.login-ecr.outputs.registry }}
          ECR_REPOSITORY: ${{ secrets.AWS_ECR_REPOSITORY }}-${{ steps.vars.outputs.short_ref }}
          AWS_REGION: ${{ secrets.AWS_REGION }}
          IMAGE_TAG: ${{ github.sha }}
        run: |
          docker compose build

      # Create and configure Amazon ECS task definition
      - name: Render Amazon ECS task definition
        id: %s-container
        uses: aws-actions/amazon-ecs-render-task-definition@v1
        with:
          task-definition: task-definition-${{ steps.vars.outputs.short_ref }}.json
          container-name: %s-${{ steps.vars.outputs.short_ref }}
          image: ${{ steps.login-ecr.outputs.registry }}/${{ secrets.AWS_ECR_REPOSITORY }}-${{ steps.vars.outputs.short_ref }}:${{ github.sha }}

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
		stack, environment, stack, environment, projectName, projectName, projectName,
		projectName, projectName, singleQuote, singleQuote, singleQuote, singleQuote,
		singleQuote, singleQuote)

	return source
}
