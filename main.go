package main

import (
	"github.com/silinternational/app-monitoring-archiver/lib/googlesheets"
	"os"
	"log"
)

func main() {
	spreadsheetID := os.Getenv("SPREADSHEET_ID")

	if spreadsheetID == "" {
		log.Fatal("Error: Environment variable for SPREADSHEET_ID is required to execute plan and migration")
	}

	nodepingToken := os.Getenv("NODEPING_TOKEN")

	if spreadsheetID == "" {
		log.Fatal("Error: Environment variable for NODEPING_TOKEN is required to execute plan and migration")
	}

	googlesheets.ArchiveResultsForMonth("AppsDev Alerts", "November", "2018", spreadsheetID, nodepingToken)
}
