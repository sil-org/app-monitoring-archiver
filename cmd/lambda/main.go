package main

import (
	"log"
	"os"
	"strconv"
	"time"

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
	if dsn := os.Getenv("SENTRY_DSN"); dsn != "" {
		initSentry(dsn)
		defer sentry.Flush(2 * time.Second)
	}

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

func initSentry(dsn string) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         dsn,
		EnableLogs:  true,
		Environment: getEnv("APP_ENV", "production"),
	})
	if err != nil {
		log.Printf("Sentry initialization failed: %v\n", err)
	}
}

func getEnv(name, fallback string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return fallback
}
