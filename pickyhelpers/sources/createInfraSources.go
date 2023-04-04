package sources

import (
	"fmt"

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
		"deploy:dev": "sst deploy --stage dev",
		"deploy:qa": "sst deploy --stage qa",
		"deploy:prod": "sst deploy --stage prod",
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

	source := `APP_NAME=app
WEB_AWS_REGION=ap-south-1`

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
			region: process.env.WEB_AWS_REGION || "ap-south-1",
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
		shortEnvironment = constants.Dev
	case constants.QA:
		shortEnvironment = constants.QA
	case constants.Production:
		shortEnvironment = constants.Prod
	}
	var buildOutput string
	stack, _ := utils.FindStackAndDatabase(camelCaseDirName)
	if stack == constants.ReactJS {
		buildOutput = "build"
	} else if stack == constants.NextJS {
		buildOutput = "out"
	}

	source := fmt.Sprintf(`import { StaticSite } from "sst/constructs";
		
export function %s({ stack }) {
	// Deploy our web app
	const site = new StaticSite(stack, "%sSite", {
		path: "%s/",
		buildCommand: "yarn run build:%s",
		buildOutput: "%s",
	});

	// Show the URLs in the output
	stack.addOutputs({
		SiteUrl: site.url || "http://localhost:3000/",
	});
}
`, camelCaseDirName, camelCaseDirName, dirName, shortEnvironment, buildOutput)
	return source
}

func BackendStackSource(database, dirName, environment string) string {
	var shortEnvironment string
	switch environment {
	case constants.Development:
		shortEnvironment = "develop"
	case constants.QA:
		shortEnvironment = constants.QA
	case constants.Production:
		shortEnvironment = constants.Production
	}
	databaseName := fmt.Sprintf("%s_%s",
		utils.FindUserInputStackName(dirName),
		constants.Database,
	)
	camelCaseDirName := strcase.ToCamel(dirName)
	var (
		dbEngineVersion string
		dbPortNumber    string
		dbEngine        string
		db_uri          string
		dbHost          string
	)
	if database == constants.PostgreSQL {
		dbEngineVersion = "PostgresEngineVersion"
		dbPortNumber = "5432"
		dbEngine = "DatabaseInstanceEngine.postgres({\n\t\t\t\tversion: PostgresEngineVersion.VER_14_2,\n\t\t\t})"
		db_uri = "`postgres://${username}:${password}@${database.dbInstanceEndpointAddress}/${db}`"
		dbHost = "POSTGRES_HOST: database.dbInstanceEndpointAddress"
	} else if database == constants.MySQL {
		dbEngineVersion = "MysqlEngineVersion"
		dbPortNumber = "3306"
		dbEngine = "DatabaseInstanceEngine.mysql({\n\t\t\t\tversion: MysqlEngineVersion.VER_8_0_23,\n\t\t\t})"
		db_uri = "`mysql://${username}:${password}@${database.dbInstanceEndpointAddress}/${db}`"
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
	const clientName = "test-sst-ecs";
	const environment = "%s";
	const clientPrefix = %s${clientName}-${environment}%s;

	const vpc = new ec2.Vpc(stack, %s${clientPrefix}-vpc%s, {
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
	const elbSG = new ec2.SecurityGroup(stack, %s${clientPrefix}-elbSG%s, {
		vpc,
		allowAllOutbound: true,
	});

	elbSG.addIngressRule(
		ec2.Peer.anyIpv4(),
		ec2.Port.tcp(80),
		"Allow http traffic"
	);

	// ECS Security groups
	const ecsSG = new ec2.SecurityGroup(stack, %s${clientPrefix}-ecsSG%s, {
		vpc,
		allowAllOutbound: true,
	});

	ecsSG.connections.allowFrom(
		elbSG,
		ec2.Port.allTcp(),
		"Application load balancer"
	);

	// Database security group
	const databaseSecurityGroup = new ec2.SecurityGroup(
		stack,
		%s${clientPrefix}-database-security-group%s,
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

	// database
	const dbUsername = "username";

	const databaseCredentialsSecret = new secretsManager.Secret(
		stack,
		%s${clientPrefix}-database-credentials-secret%s,
		{
			secretName: %s${clientPrefix}-database-credentials%s,
			description: "Database Credentials",
			generateSecretString: {
				excludeCharacters: "\"@/\\ '",
				generateStringKey: "password",
				passwordLength: 30,
				secretStringTemplate: JSON.stringify({ username: dbUsername }),
			},
		}
	);

	const databaseCredentials = Credentials.fromSecret(
		databaseCredentialsSecret,
		dbUsername
	);

	const database = new DatabaseInstance(
		stack,
		%s${clientPrefix}-database-instance%s,
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
			databaseName: "%s",
		}
	);

	// Elasticache
	const redisSubnetGroup = new elasticcache.CfnSubnetGroup(
		stack,
		%s${clientPrefix}-redis-subnet-group%s,
		{
			description: "Subnet group for the redis cluster",
			subnetIds: vpc.privateSubnets.map((subnet) => subnet.subnetId),
			cacheSubnetGroupName: %s${clientPrefix}-redis-subnet-group%s,
		}
	);

	const redisSecurityGroup = new ec2.SecurityGroup(
		stack,
		%s${clientPrefix}-redis-security-group%s,
		{
			vpc,
			allowAllOutbound: true,
			description: "Security group for the redis cluster",
		}
	);

	redisSecurityGroup.addIngressRule(
		ecsSG,
		ec2.Port.tcp(6379),
		"Permit the redis cluster to accept requests from the fargate service"
	);

	const redisCache = new elasticcache.CfnCacheCluster(
		stack,
		%s${clientPrefix}-redis-cache%s,
		{
			engine: "redis",
			cacheNodeType: "cache.t3.micro",
			numCacheNodes: 1,
			clusterName: %s${clientPrefix}-redis-cluster%s,
			vpcSecurityGroupIds: [redisSecurityGroup.securityGroupId],
			cacheSubnetGroupName: redisSubnetGroup.ref,
			engineVersion: "6.2",
		}
	);

	redisCache.addDependency(redisSubnetGroup);

	// Creating your ECS
	const cluster = new ecs.Cluster(stack, %s${clientPrefix}-cluster%s, {
		clusterName: %s${clientPrefix}-cluster%s,
		vpc,
	});

	// // Creating your Load Balancer
	const elb = new elasticloadbalancing.ApplicationLoadBalancer(
		stack,
		%s${clientPrefix}-elb%s,
		{
			vpc,
			vpcSubnets: { subnets: vpc.publicSubnets },
			internetFacing: true,
		}
	);

	elb.addSecurityGroup(elbSG);

	// Creating your target group
	const targetGroupHttp = new elasticloadbalancing.ApplicationTargetGroup(
		stack,
		%s${clientPrefix}-target%s,
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

	listener.addTargetGroups(%s${clientPrefix}-tg%s, {
		targetGroups: [targetGroupHttp],
	});

	const taskRole = new iam.Role(stack, %s${clientPrefix}-task-role%s, {
		assumedBy: new iam.ServicePrincipal("ecs-tasks.amazonaws.com"),
		roleName: %s${clientPrefix}-task-role%s,
		description: "Role that the api task definitions use to run the api code",
	});

	const taskDefinition = new ecs.TaskDefinition(stack, %s${clientPrefix}-task%s, {
		family: %s${clientPrefix}-task%s,
		compatibility: ecs.Compatibility.EC2_AND_FARGATE,
		cpu: "256",
		memoryMiB: "512",
		networkMode: ecs.NetworkMode.AWS_VPC,
		taskRole: taskRole,
	});

	const username = databaseCredentialsSecret
		.secretValueFromJson("username")
		.toString();
	const password = databaseCredentialsSecret
		.secretValueFromJson("password")
		.toString();

	const db = "sst_test_database";

	const DB_URI = %s;

	const image = ecs.ContainerImage.fromAsset("%s/", {
		exclude: ["node_modules", ".git"],
		platform: Platform.LINUX_AMD64,
		buildArgs: {
			ENVIRONMENT_NAME: "%s",
		},
	});

	const container = taskDefinition.addContainer(%s${clientPrefix}-container%s, {
		image,
		memoryLimitMiB: 512,
		environment: {
			BUILD_NAME: "%s",
			ENVIRONMENT_NAME: "%s",
			DB_URI,
			%s,
			REDIS_HOST: redisCache.attrRedisEndpointAddress,
		},
		logging: ecs.LogDriver.awsLogs({
			streamPrefix: %s${clientPrefix}-log-group%s,
		}),
	});

	container.addPortMappings({ containerPort: 9000 });

	const service = new ecs.FargateService(stack, %s${clientPrefix}-service%s, {
		cluster,
		desiredCount: 1,
		taskDefinition,
		securityGroups: [ecsSG],
		assignPublicIp: true,
	});

	service.attachToApplicationTargetGroup(targetGroupHttp);

	new CfnOutput(stack, "database-host", {
		exportName: "database-host",
		value: database.dbInstanceEndpointAddress,
	});

	new CfnOutput(stack, "database-name", {
		exportName: "database-name",
		value: db,
	});

	new CfnOutput(stack, "redis-host", {
		exportName: "redis-host",
		value: redisCache.attrRedisEndpointAddress,
	});
}
`, dbEngineVersion, camelCaseDirName, shortEnvironment, "`", "`", "`", "`",
		"`", "`", "`", "`", "`", "`", dbPortNumber, "`", "`", "`", "`", "`", "`",
		dbEngine, databaseName, "`", "`", "`", "`", "`", "`", "`", "`", "`", "`",
		"`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`",
		"`", "`", "`", db_uri, dirName, environment, "`", "`", shortEnvironment,
		environment, dbHost, "`", "`", "`", "`",
	)

	return source
}

func EnvDevSource(environment string) string {

	source := fmt.Sprintf(`NAME=Node Template
NODE_ENV=%s
ENVIRONMENT_NAME=%s
PORT=9000`, environment, environment)

	return source
}
