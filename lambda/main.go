package main

import (
	"github.com/silinternational/app-monitoring-archiver/lib/googlesheets"
	"github.com/aws/aws-lambda-go/lambda"
	"os"
	"log"
	"github.com/silinternational/app-monitoring-archiver/cmd"
	"strconv"
)


type ArchiveToGoogleSheetsConfig struct {
	ContactGroupName string
	Period 			 string
	SpreadSheetID 	 string
	CountLimit 		 string
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

	intCountLimit, err := strconv.Atoi(config.CountLimit)
	if err != nil {
		log.Fatalf("Error converting CountLimit of %s to integer. %v", config.CountLimit, err)
	}

	googlesheets.ArchiveResultsForMonth(
		config.ContactGroupName,
		config.Period,
		config.SpreadSheetID,
		nodepingToken,
		intCountLimit,
	)
	return nil
}