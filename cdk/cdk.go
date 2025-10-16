package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"

	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CdkStackProps struct {
	awscdk.StackProps
}

func NewCdkStack(scope constructs.Construct, id string, props *CdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	const appName = "app-monitoring-archiver"
	const customer = "gtis"
	envName := os.Getenv("STAGE")

	functionName := "lambda_function-" + appName + "-" + customer + "-" + envName

	nodepingToken := os.Getenv("NODEPING_TOKEN")
	contactGroupName := os.Getenv("CONTACT_GROUP_NAME")
	countLimit := os.Getenv("COUNT_LIMIT")
	period := os.Getenv("PERIOD")
	spreadsheetID := os.Getenv("SPREADSHEET_ID")

	googleAuthClientEmail := os.Getenv("GOOGLE_AUTH_CLIENT_EMAIL")
	googleAuthPrivateKeyID := os.Getenv("GOOGLE_AUTH_PRIVATE_KEY_ID")
	googleAuthPrivateKey := os.Getenv("GOOGLE_AUTH_PRIVATE_KEY")
	googleAuthTokenURI := os.Getenv("GOOGLE_AUTH_TOKEN_URI")


	stack := awscdk.NewStack(scope, &id, &sprops)

	logGroup := awslogs.NewLogGroup(stack, jsii.String("LambdaLogGroup"), &awslogs.LogGroupProps{
		LogGroupName:  jsii.String("/aws/lambda/" + functionName + "-cdk"),
		Retention:     awslogs.RetentionDays_TWO_MONTHS,
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})


	function := awslambda.NewFunction(stack,  jsii.String(functionName), &awslambda.FunctionProps{
		Code: awslambda.Code_FromAsset(jsii.String("../bin"), nil),
		Environment: &map[string]*string{
			"NODEPING_TOKEN":             &nodepingToken,
			"CONTACT_GROUP_NAME":         &contactGroupName,
			"COUNT_LIMIT":                &countLimit,
			"PERIOD":                     &period,
			"SPREADSHEET_ID":             &spreadsheetID,
			"GOOGLE_AUTH_CLIENT_EMAIL":   &googleAuthClientEmail,
			"GOOGLE_AUTH_PRIVATE_KEY_ID": &googleAuthPrivateKeyID,
			"GOOGLE_AUTH_PRIVATE_KEY":    &googleAuthPrivateKey,
			"GOOGLE_AUTH_TOKEN_URI":      &googleAuthTokenURI,
		},
		FunctionName:  jsii.String(appName + "-" + customer + "-" + envName),
		Handler:       jsii.String("bootstrap"),
		LoggingFormat: awslambda.LoggingFormat_JSON,
		LogGroup:      logGroup,
		Runtime:       awslambda.Runtime_PROVIDED_AL2023(),
		Timeout:       awscdk.Duration_Seconds(jsii.Number(600)),
	})

	rule := awsevents.NewRule(stack, jsii.String("ScheduleRule"), &awsevents.RuleProps{
		Schedule: awsevents.Schedule_Cron(&awsevents.CronOptions{
			Minute:  jsii.String("30"),
			Hour:    jsii.String("3"),
			Day:     jsii.String("1"),
			Month:   jsii.String("*"),
		}),
	})

	rule.AddTarget(awseventstargets.NewLambdaFunction(function, &awseventstargets.LambdaFunctionProps{
		RetryAttempts: jsii.Number(0),
	}))

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewCdkStack(app, "AppMonitoringArchiver", &CdkStackProps{
		awscdk.StackProps{
			Env: &awscdk.Environment{
				Region: jsii.String(os.Getenv("AWS_REGION")),
			},
			Tags: &map[string]*string{
				"managed_by":        jsii.String("cdk"),
				"itse_app_name":     jsii.String("app-monitoring-archiver"),
				"itse_app_customer": jsii.String("gtis"),
				"itse_app_env":      jsii.String("production"),
			},
		},
	})

	app.Synth(nil)
}
