package main

import (
	"fmt"
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
	SentryDSN        string
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

	nodePingToken, err := getRequiredEnv(cmd.NodePingTokenKey)
	if err != nil {
		sentry.CaptureException(err)
		log.Fatalln(err)
	}

	intCountLimit, err := strconv.Atoi(config.CountLimit)
	if err != nil {
		err = fmt.Errorf("error converting CountLimit '%s' to integer: %w", config.CountLimit, err)
		sentry.CaptureException(err)
		log.Fatalln(err)
	}

	err = googlesheets.ArchiveResultsForMonth(
		config.ContactGroupName,
		config.Period,
		config.SpreadSheetID,
		nodePingToken,
		intCountLimit,
	)
	if err != nil {
		sentry.CaptureException(err)
		log.Fatalln(err)
	}
	return nil
}

func initSentry(dsn string) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         dsn,
		Environment: getEnv("APP_ENV", "prod"),
	})
	if err != nil {
		log.Println("Sentry initialization failed:", err)
	}
}

func getEnv(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func getRequiredEnv(name string) (string, error) {
	if value := os.Getenv(name); value != "" {
		return value, nil
	}
	return "", fmt.Errorf("missing required env var: %s", name)
}
