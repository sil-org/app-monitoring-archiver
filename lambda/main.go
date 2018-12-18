package main

import (
	"github.com/silinternational/app-monitoring-archiver/lib/googlesheets"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
	"log"
	"github.com/silinternational/app-monitoring-archiver/cmd"
)


type ArchiveToGoogleSheetsConfig struct {
	ContactGroupName string
	Period 			 string
	SpreadSheetID 	 string
	CountLimit 		 int
}

func main() {
	lambda.Start(handler)
}

func handler(config ArchiveToGoogleSheetsConfig) error {
	if config.Period == "" {
		config.Period = "LastMonth"
	}

	nodepingToken := os.Getenv(cmd.NodepingTokenKey)

	if nodepingToken == "" {
		log.Fatal("Error: Environment variable for NODEPING_TOKEN is required to execute plan and migration \n")
	}

	googlesheets.ArchiveResultsForMonth(
		config.ContactGroupName,
		config.Period,
		config.SpreadSheetID,
		nodepingToken,
		config.CountLimit,
	)
	return nil
}