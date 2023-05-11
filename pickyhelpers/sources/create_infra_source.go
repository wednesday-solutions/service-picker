package sources

import (
	"fmt"
	"strconv"

	"github.com/iancoleman/strcase"
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/utils"
)

func PackageDotJsonSource() string {

	source := `{
	"name": "app",
	"version": "0.0.0",
	"private": true,
	"type": "module",
	"scripts": {
		"dev": "sst dev",
		"build": "sst build",
		"deploy:dev": "sst deploy --stage dev --outputs-file outputs.json",
		"deploy:qa": "sst deploy --stage qa --outputs-file outputs.json",
		"deploy:prod": "sst deploy --stage prod --outputs-file outputs.json",
		"remove:dev": "sst remove --stage dev",
		"remove:qa": "sst remove --stage qa",
		"remove:prod": "sst remove --stage prod",
		"console": "sst console",
		"typecheck": "tsc --noEmit"
	},
	"devDependencies": {
		"sst": "^2.0.18",
		"aws-cdk-lib": "2.62.2",
		"constructs": "10.1.156",
		"dotenv": "^10.0.0"
	},
	"workspaces": [
		"/*"
	]
}`
	return source
}

func EnvFileSource() string {

	source := `APP_NAME=web-app
AWS_REGION=ap-south-1`

	return source
}

func SstConfigSource() string {

	source := `import dotenv from "dotenv";
{{{sstImportStacks existingDirectories}}}
dotenv.config({ path: ".env" });

export default {
	config(_input) {
		return {
			name: process.env.APP_NAME || "web-app",
			region: process.env.AWS_REGION || "ap-south-1",
		};
	},
	stacks(app) {
		// deploy stacks
		{{deployStacks existingDirectories}}
	},
};
`

	return source
}

func WebStackSource(dirName, camelCaseDirName, environment string) string {
	var shortEnvironment string
	switch environment {
	case constants.Development:
		environment = constants.Dev
		shortEnvironment = constants.Dev
	case constants.QA:
		shortEnvironment = constants.QA
	case constants.Production:
		environment = constants.Prod
		shortEnvironment = constants.Prod
	}
	buildOutput, singleQuote := "", "`"
	stack, _ := utils.FindStackAndDatabase(dirName)
	if stack == constants.ReactJS {
		buildOutput = "build"
	} else if stack == constants.NextJS {
		buildOutput = "out"
	}

	source := fmt.Sprintf(`import { StaticSite } from "sst/constructs";

export function %s({ stack }) {
	const bucketprefix = "%s";
	const environment = "%s";
	const bucketName = %s${bucketprefix}-${environment}%s;

	// Deploy our web app
	const site = new StaticSite(stack, "%sSite", {
		path: "%s",
		buildCommand: "yarn run build:%s",
		buildOutput: "%s",
		cdk: {
			bucket: {
				bucketName,
			},
			distribution: {
				comment: %sDistribution for ${bucketName}%s,
			},
		},
	});

	// Show the URLs in the output
	stack.addOutputs({
		siteUrl: site.url || "http://localhost:3000/",
		distributionId: site.cdk?.distribution?.distributionId,
		bucketName: site.cdk?.bucket?.bucketName,
	});
}
`, camelCaseDirName, dirName, environment, singleQuote, singleQuote, camelCaseDirName,
		dirName, shortEnvironment, buildOutput, singleQuote, singleQuote)

	return source
}

func BackendStackSource(database, dirName, environment string) string {
	var envFileEnvironment, shortEnvironment string
	switch environment {
	case constants.Development:
		envFileEnvironment = fmt.Sprintf(".%s", constants.Development)
		shortEnvironment = constants.Dev
	case constants.QA:
		envFileEnvironment = fmt.Sprintf(".%s", constants.QA)
		shortEnvironment = constants.QA
	case constants.Production:
		envFileEnvironment = ""
		shortEnvironment = constants.Prod
	}
	userInputStackName := utils.FindUserInputStackName(dirName)
	singleQuote := "`"
	dbName := fmt.Sprintf("%s%s_%s_%s%s",
		singleQuote,
		strcase.ToSnake(userInputStackName),
		constants.Database,
		"${environment}",
		singleQuote,
	)
	dbUsername := constants.DBUsername
	camelCaseDirName := strcase.ToCamel(dirName)
	var (
		dbEngineVersion string
		dbPortNumber    string
		dbEngine        string
		dbUri           string
		dbHost          string
	)
	backendPortNumber := utils.FetchExistingPortNumber(dirName, constants.BackendPort)
	redisPortNumber := strconv.Itoa(constants.RedisPortNumber)

	if database == constants.PostgreSQL {
		dbEngineVersion = "PostgresEngineVersion"
		dbPortNumber = utils.FetchExistingPortNumber(dirName, constants.PostgresPort)
		dbEngine = "DatabaseInstanceEngine.postgres({\n\t\t\t\tversion: PostgresEngineVersion.VER_14_2,\n\t\t\t})"
		dbUri = "`postgres://${username}:${password}@${database.dbInstanceEndpointAddress}/${dbName}`"
		dbHost = "POSTGRES_HOST: database.dbInstanceEndpointAddress"
	} else if database == constants.MySQL {
		dbEngineVersion = "MysqlEngineVersion"
		dbPortNumber = utils.FetchExistingPortNumber(dirName, constants.MysqlPort)
		dbEngine = "DatabaseInstanceEngine.mysql({\n\t\t\t\tversion: MysqlEngineVersion.VER_8_0_31,\n\t\t\t})"
		dbUri = "`mysql://${username}:${password}@${database.dbInstanceEndpointAddress}/${dbName}`"
		dbHost = "MYSQL_HOST: database.dbInstanceEndpointAddress"
	}

	source := fmt.Sprintf(`import * as cdk from "aws-cdk-lib";
import * as ecs from "aws-cdk-lib/aws-ecs";
import * as ec2 from "aws-cdk-lib/aws-ec2";
import * as iam from "aws-cdk-lib/aws-iam";
import * as elasticloadbalancing from "aws-cdk-lib/aws-elasticloadbalancingv2";
import * as secretsManager from "aws-cdk-lib/aws-secretsmanager";
import {
	DatabaseInstance,
	DatabaseInstanceEngine,
	%s,
	Credentials,
} from "aws-cdk-lib/aws-rds";
import { CfnOutput } from "aws-cdk-lib";
import * as elasticcache from "aws-cdk-lib/aws-elasticache";
import { Platform } from "aws-cdk-lib/aws-ecr-assets";

export function %s({ stack }) {
	const clientName = "%s";
	const environment = "%s";
	const dbName = %s;
	const dbUsername = "%s";

	const vpc = new ec2.Vpc(stack, %s${clientName}-vpc-${environment}%s, {
		maxAzs: 3,
		natGateways: 1,
		subnetConfiguration: [
			{
				name: "public-subnet",
				subnetType: ec2.SubnetType.PUBLIC,
				cidrMask: 24,
			},
			{
				cidrMask: 24,
				name: "private-subnet",
				subnetType: ec2.SubnetType.PRIVATE_WITH_EGRESS,
			},
		],
	});

	// Load Balancer Security groups
	const elbSG = new ec2.SecurityGroup(
		stack,
		%s${clientName}-elb-security-group-${environment}%s,
		{
			vpc,
			allowAllOutbound: true,
		}
	);

	elbSG.addIngressRule(
		ec2.Peer.anyIpv4(),
		ec2.Port.tcp(80),
		"Allow http traffic"
	);

	// ECS Security groups
	const ecsSG = new ec2.SecurityGroup(
		stack,
		%s${clientName}-ecs-security-group-${environment}%s,
		{
			vpc,
			allowAllOutbound: true,
		}
	);

	ecsSG.connections.allowFrom(
		elbSG,
		ec2.Port.allTcp(),
		"Application load balancer"
	);

	// Database security group
	const databaseSecurityGroup = new ec2.SecurityGroup(
		stack,
		%s${clientName}-database-security-group-${environment}%s,
		{
			vpc,
			allowAllOutbound: false,
		}
	);

	databaseSecurityGroup.addIngressRule(
		ecsSG,
		ec2.Port.tcp(%s),
		"Permit the database to accept requests from the fargate service"
	);

	const databaseCredentialsSecret = new secretsManager.Secret(
		stack,
		%s${clientName}-database-credentials-secret-${environment}%s,
		{
			secretName: %s${clientName}-database-credentials-${environment}%s,
			description: %sDatabase credentials for ${clientName}-${environment}%s,
			generateSecretString: {
				excludeCharacters: "'\\;@$\"%s!/ ",
				generateStringKey: "password",
				passwordLength: 30,
				secretStringTemplate: JSON.stringify({ username: dbUsername }),
				excludePunctuation: true,
			},
		}
	);

	const databaseCredentials = Credentials.fromSecret(
		databaseCredentialsSecret,
		dbUsername
	);

	const database = new DatabaseInstance(
		stack,
		%s${clientName}-database-instance-${environment}%s,
		{
			vpc,
			securityGroups: [databaseSecurityGroup],
			credentials: databaseCredentials,
			engine: %s,
			removalPolicy: cdk.RemovalPolicy.DESTROY, // CHANGE TO .SNAPSHOT FOR PRODUCTION
			instanceType: ec2.InstanceType.of(
				ec2.InstanceClass.BURSTABLE3,
				ec2.InstanceSize.MICRO
			),
			vpcSubnets: {
				subnetType: ec2.SubnetType.PRIVATE_WITH_EGRESS,
			},
			backupRetention: cdk.Duration.days(1),
			allocatedStorage: 10,
			maxAllocatedStorage: 30,
			databaseName: dbName,
		}
	);

	// Elasticache
	const redisSubnetGroup = new elasticcache.CfnSubnetGroup(
		stack,
		%s${clientName}-redis-subnet-group-${environment}%s,
		{
			description: "Subnet group for the redis cluster",
			subnetIds: vpc.privateSubnets.map((subnet) => subnet.subnetId),
			cacheSubnetGroupName: %s${clientName}-redis-subnet-group-${environment}%s,
		}
	);

	const redisSecurityGroup = new ec2.SecurityGroup(
		stack,
		%s${clientName}-redis-security-group-${environment}%s,
		{
			vpc,
			allowAllOutbound: true,
			description: "Security group for the redis cluster",
		}
	);

	redisSecurityGroup.addIngressRule(
		ecsSG,
		ec2.Port.tcp(%s),
		"Permit the redis cluster to accept requests from the fargate service"
	);

	const redisCache = new elasticcache.CfnCacheCluster(
		stack,
		%s${clientName}-redis-cache-${environment}%s,
		{
			engine: "redis",
			cacheNodeType: "cache.t3.micro",
			numCacheNodes: 1,
			clusterName: %s${clientName}-redis-cluster-${environment}%s,
			vpcSecurityGroupIds: [redisSecurityGroup.securityGroupId],
			cacheSubnetGroupName: redisSubnetGroup.ref,
			engineVersion: "6.2",
		}
	);

	redisCache.addDependency(redisSubnetGroup);

	// Creating your ECS
	const cluster = new ecs.Cluster(stack, %s${clientName}-cluster%s, {
		clusterName: %s${clientName}-cluster-${environment}%s,
		vpc,
	});

	// Creating your Load Balancer
	const elb = new elasticloadbalancing.ApplicationLoadBalancer(
		stack,
		%s${clientName}-elb-${environment}%s,
		{
			vpc,
			vpcSubnets: { subnets: vpc.publicSubnets },
			internetFacing: true,
			loadBalancerName: %s${clientName}-elb-${environment}%s,
		}
	);

	elb.addSecurityGroup(elbSG);

	// Creating your target group
	const targetGroupHttp = new elasticloadbalancing.ApplicationTargetGroup(
		stack,
		%s${clientName}-target-${environment}%s,
		{
			port: 80,
			vpc,
			protocol: elasticloadbalancing.ApplicationProtocol.HTTP,
			targetType: elasticloadbalancing.TargetType.IP,
		}
	);

	targetGroupHttp.configureHealthCheck({
		path: "/",
		protocol: elasticloadbalancing.Protocol.HTTP,
	});

	// Adding your listeners
	const listener = elb.addListener("Listener", {
		open: true,
		port: 80,
	});

	listener.addTargetGroups(%s${clientName}-target-group-${environment}%s, {
		targetGroups: [targetGroupHttp],
	});

	const taskRole = new iam.Role(
		stack,
		%s${clientName}-ecs-task-role-${environment}%s,
		{
			assumedBy: new iam.ServicePrincipal("ecs-tasks.amazonaws.com"),
			roleName: %s${clientName}-task-role-${environment}%s,
		}
	);

	const executionRole = new iam.Role(
    stack,
    %s${clientName}-ecs-execution-role-${environment}%s,
    {
      assumedBy: new iam.ServicePrincipal("ecs-tasks.amazonaws.com"),
      roleName: %s${clientName}-ecs-execution-role-${environment}%s,
    }
  );

  databaseCredentialsSecret.grantRead(taskRole);
  databaseCredentialsSecret.grantRead(executionRole);

	const taskDefinition = new ecs.TaskDefinition(
		stack,
		%s${clientName}-task-${environment}%s, 
		{
			family: %s${clientName}-task-definition-${environment}%s,
			compatibility: ecs.Compatibility.EC2_AND_FARGATE,
			cpu: "512",
			memoryMiB: "1024",
			networkMode: ecs.NetworkMode.AWS_VPC,
			taskRole: taskRole,
			executionRole: executionRole
		}
	);

	const username = databaseCredentialsSecret
		.secretValueFromJson("username")
		.toString();
	const password = databaseCredentialsSecret
		.secretValueFromJson("password")
		.toString();

	const dbURI = %s;

	const image = ecs.ContainerImage.fromAsset("%s/", {
		exclude: ["node_modules", ".git"],
		platform: Platform.LINUX_AMD64,
		buildArgs: {
			ENVIRONMENT_NAME: "%s",
			BUILD_NAME: environment,
		},
	});

	const container = taskDefinition.addContainer(
		%s${clientName}-container-${environment}%s, 
		{
			image,
			cpu: 512,
			memoryLimitMiB: 1024,
			environment: {
				BUILD_NAME: environment,
				// ENVIRONMENT_NAME: "%s",
				DB_URI: dbURI,
				%s,
				REDIS_HOST: redisCache.attrRedisEndpointAddress,
				REDIS_PORT: "%s",
			},
			logging: ecs.LogDriver.awsLogs({
				streamPrefix: %s${clientName}-log-group-${environment}%s,
			}),
		}
	);

	container.addPortMappings({ containerPort: %s });

	const service = new ecs.FargateService(
		stack,
		%s${clientName}-service-${environment}%s,
		{
			cluster,
			desiredCount: 1,
			taskDefinition,
			securityGroups: [ecsSG],
			assignPublicIp: true,
			serviceName: %s${clientName}-service-${environment}%s,
		}
	);

	service.attachToApplicationTargetGroup(targetGroupHttp);

	new CfnOutput(stack, "databaseHost", {
		value: database.dbInstanceEndpointAddress,
	});

	new CfnOutput(stack, "databaseName", {
		value: dbName,
	});

	new CfnOutput(stack, "redisHost", {
		value: redisCache.attrRedisEndpointAddress,
	});

	new CfnOutput(stack, "loadBalancerDns", {
		value: elb.loadBalancerDnsName,
	});

	new CfnOutput(stack, "awsRegion", {
		value: stack.region,
	});

  new CfnOutput(stack, "elasticContainerRegistryRepo", {
    value: stack.synthesizer.repositoryName,
  });

  new CfnOutput(stack, "image", {
    value: container.imageName,
  });

  new CfnOutput(stack, "taskDefinition", {
    value: taskDefinition.taskDefinitionArn,
  });

  new CfnOutput(stack, "taskRole", {
    value: taskRole.roleArn,
  });

  new CfnOutput(stack, "executionRole", {
    value: executionRole.roleArn,
  });

  new CfnOutput(stack, "family", {
    value: taskDefinition.family,
  });

  new CfnOutput(stack, "containerName", {
    value: container.containerName,
  });

  new CfnOutput(stack, "containerPort", {
    value: container.containerPort.toString(),
  });

  new CfnOutput(stack, "logDriver", {
    value: container.logDriverConfig.logDriver,
  });

  new CfnOutput(stack, "logDriverOptions", {
    value: JSON.stringify(container.logDriverConfig.options),
  });

	new CfnOutput(stack, "serviceName", {
		value: service.serviceName,
	});

	new CfnOutput(stack, "clusterName", {
    value: cluster.clusterName,
  });

	new CfnOutput(stack, "secretName", {
    value: databaseCredentialsSecret.secretName,
  });

	new CfnOutput(stack, "secretArn", {
    value: databaseCredentialsSecret.secretArn,
  });
}
`,
		dbEngineVersion, camelCaseDirName, userInputStackName, shortEnvironment,
		dbName, dbUsername, singleQuote, singleQuote, singleQuote, singleQuote,
		singleQuote, singleQuote, singleQuote, singleQuote, dbPortNumber, singleQuote,
		singleQuote, singleQuote, singleQuote, singleQuote, singleQuote, singleQuote,
		singleQuote, singleQuote, dbEngine, singleQuote, singleQuote, singleQuote,
		singleQuote, singleQuote, singleQuote, redisPortNumber, singleQuote, singleQuote,
		singleQuote, singleQuote, singleQuote, singleQuote, singleQuote, singleQuote,
		singleQuote, singleQuote, singleQuote, singleQuote, singleQuote, singleQuote,
		singleQuote, singleQuote, singleQuote, singleQuote, singleQuote, singleQuote,
		singleQuote, singleQuote, singleQuote, singleQuote, singleQuote, singleQuote,
		singleQuote, singleQuote, dbUri, dirName, envFileEnvironment, singleQuote,
		singleQuote, environment, dbHost, redisPortNumber, singleQuote, singleQuote,
		backendPortNumber, singleQuote, singleQuote, singleQuote, singleQuote,
	)
	return source
}

// EnvSource return the source string with respect to the given environment.
func EnvSource(dirName, environment, database string, backendObj utils.BackendOutputKeys) string {

	var dbHost, dbUser, dbName string
	if database == constants.PostgreSQL {
		dbUser = constants.PostgresUser
		dbHost = constants.PostgresHost
		dbName = constants.PostgresDB
	} else if database == constants.MySQL {
		dbUser = constants.MysqlUser
		dbHost = constants.MysqlHost
		dbName = constants.MysqlDatabase
	}
	redisPortNumber := strconv.Itoa(constants.RedisPortNumber)
	backendPortNumber := utils.FetchExistingPortNumber(dirName, constants.BackendPort)

	source := fmt.Sprintf(`NAME=Node Template
NODE_ENV=%s
ENVIRONMENT_NAME=%s
PORT=%s
`,
		environment,
		environment,
		backendPortNumber,
	)
	if backendObj.DatabaseHost != "" && backendObj.DatabaseName != "" && backendObj.RedisHost != "" {
		source = fmt.Sprintf(`%s
%s=%s
%s=%s
%s=%s
%s=%s
%s=%s
`,
			source,
			dbUser,
			constants.DBUsername,
			dbHost,
			backendObj.DatabaseHost,
			dbName,
			backendObj.DatabaseName,
			constants.RedisHost,
			backendObj.RedisHost,
			constants.RedisPort,
			redisPortNumber,
		)
	}
	return source
}

func ParseSstOutputsSource() string {
	source := `import * as fs from "fs";

function parseOutputs() {
	const outputFile = "./.sst/outputs.json";
	let fileContent = JSON.parse(fs.readFileSync(outputFile, "utf-8"));

	Object.keys(fileContent).some((k) => {
		if (k.endsWith("Pg") || k.endsWith("Mysql")) {
			if (fileContent[k]?.logDriverOptions) {
				fileContent[k].logDriverOptions = JSON.parse(
					fileContent[k].logDriverOptions
				);
			}
		}
	});
	fileContent = JSON.stringify(fileContent, null, 2);
	fs.writeFileSync(outputFile, fileContent);
}
parseOutputs();	
`
	return source
}
