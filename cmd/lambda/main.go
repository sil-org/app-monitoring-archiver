package main

import (
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/getsentry/sentry-go"

	"github.com/sil-org/app-monitoring-archiver/cmd"
	"github.com/sil-org/app-monitoring-archiver/lib/googlesheets"
)

type ArchiveToGoogleSheetsConfig struct {
	ContactGroupName string
	Period           string
	SpreadSheetID    string
	CountLimit       string
}

func main() {
	sentryInit()
	lambda.Start(handler)
}

func handler(config ArchiveToGoogleSheetsConfig) error {
	if config.Period == "" {
		config.Period = "LastMonth"
	}

	nodePingToken := os.Getenv(cmd.NodePingTokenKey)

	if nodePingToken == "" {
		log.Fatal("Error: Environment variable for NODEPING_TOKEN is required to execute plan and migration")
	}

	intCountLimit, err := strconv.Atoi(config.CountLimit)
	if err != nil {
		log.Fatalf("Error converting CountLimit of %s to integer. %v", config.CountLimit, err)
	}

	googlesheets.ArchiveResultsForMonth(
		config.ContactGroupName,
		config.Period,
		config.SpreadSheetID,
		nodePingToken,
		intCountLimit,
	)
	return nil
}

func sentryInit() {
	dsn := os.Getenv("SENTRY_DSN")
	if dsn == "" {
		return
	}

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:         dsn,
		EnableLogs:  true,
		Environment: os.Getenv("APP_ENV"),
	}); err != nil {
		log.Printf("Sentry initialization failed: %v\n", err)
	}
}
