package sources

import "fmt"

func PackageDotJsonSource() string {

	source := `{
	"name": "app",
	"version": "0.0.0",
	"private": true,
	"type": "module",
	"scripts": {
		"dev": "sst dev",
		"build": "sst build",
		"deploy": "sst deploy",
		"remove": "sst remove",
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

func SstConfigJsSource() string {

	source := `import dotenv from "dotenv";
{{#if webStatus}}
import { WebStack } from "./stacks/WebStack";
{{/if}}
{{#if backendStatus}}
import { BackendStack } from "./stack/BackendStack";
{{/if}}

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
		{{deployStacks webStatus backendStatus}}
	},
};
`

	return source
}

func WebStackJsSource() string {

	source := `import { StaticSite } from "sst/constructs";
		
export function WebStack({ stack }) {
	// Deploy our web app
	const site = new StaticSite(stack, "WebSite", {
		path: "/",
		buildCommand: "yarn run build:prod",
		buildOutput: "out",
	});

	// Show the URLs in the output
	stack.addOutputs({
		SiteUrl: site.url || "http://localhost:3000/",
	});
}
`
	return source
}

func BackendStackJsSource() string {

	source := fmt.Sprintf(`import * as cdk from "aws-cdk-lib";
import * as ecs from "aws-cdk-lib/aws-ecs";
import * as ec2 from "aws-cdk-lib/aws-ec2";
import * as iam from "aws-cdk-lib/aws-iam";
import * as elasticloadbalancing from "aws-cdk-lib/aws-elasticloadbalancingv2";
import * as secretsManager from "aws-cdk-lib/aws-secretsmanager";
import {
	DatabaseInstance,
	DatabaseInstanceEngine,
	MysqlEngineVersion,
	Credentials,
} from "aws-cdk-lib/aws-rds";
import { CfnOutput } from "aws-cdk-lib";
import * as elasticcache from "aws-cdk-lib/aws-elasticache";
import { Platform } from "aws-cdk-lib/aws-ecr-assets";

export function BackendStack({ stack }) {
	const clientName = "test-sst-ecs";
	const environment = "develop";
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
		ec2.Port.tcp(3306),
		"Permit the database to accept requests from the fargate service"
	);

	// database
	const mysqlUsername = "username";

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
				secretStringTemplate: JSON.stringify({ username: mysqlUsername }),
			},
		}
	);

	const mysqlCredentials = Credentials.fromSecret(
		databaseCredentialsSecret,
		mysqlUsername
	);

	const database = new DatabaseInstance(
		stack,
		%s${clientPrefix}-database-instance%s,
		{
			vpc,
			securityGroups: [databaseSecurityGroup],
			credentials: mysqlCredentials,
			engine: DatabaseInstanceEngine.mysql({
				version: MysqlEngineVersion.VER_8_0_23,
			}),
			removalPolicy: cdk.RemovalPolicy.DESTROY, // CHANGE TO .SNAPSHOT FOR PRODUCTION
			instanceType: ec2.InstanceType.of(
				ec2.InstanceClass.BURSTABLE3,
				ec2.InstanceSize.MICRO
			),
			vpcSubnets: {
				subnetType: ec2.SubnetType.PRIVATE_WITH_EGRESS,
			},
			backupRetention: cdk.Duration.days(7),
			allocatedStorage: 10,
			maxAllocatedStorage: 30,
			databaseName: "sst_test_database",
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

	const DB_URI = %smysql://${username}:${password}@${database.dbInstanceEndpointAddress}/${db}%s;

	const image = ecs.ContainerImage.fromAsset("backend/", {
		exclude: ["node_modules", ".git"],
		platform: Platform.LINUX_AMD64,
		buildArgs: {
			ENVIRONMENT_NAME: "development",
		},
	});

	const container = taskDefinition.addContainer(%s${clientPrefix}-container%s, {
		image,
		memoryLimitMiB: 512,
		environment: {
			BUILD_NAME: "develop",
			ENVIRONMENT_NAME: "development",
			DB_URI,
			MYSQL_HOST: database.dbInstanceEndpointAddress,
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
`, "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`",
		"`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`",
		"`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`", "`",
		"`", "`", "`", "`", "`", "`", "`",
	)

	return source
}

func EnvDevSource() string {

	source := `NAME=Node Template (DEV)
NODE_ENV=development
ENVIRONMENT_NAME=development
PORT=9000`

	return source
}
