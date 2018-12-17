package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"github.com/silinternational/app-monitoring-archiver/lib/googlesheets"
)


var contactGroupName string
var spreadsheetID string
var countLimit int

var runCmd = &cobra.Command{
	Use: "run",
	Short: "Archive Nodeping results",
	Long: "Get the uptime results from Nodeping and write them to Google Sheets.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if contactGroupName == "" {
			log.Fatal(`Error: The 'contact-group' flag is required (e.g. -g "AppsDev Alerts"). \n`)
		}

		if spreadsheetID == "" {
			log.Fatal(`Error: The 'spreadsheetID' flag is required (found in its url). \n`)
		}

		runArchive()
	},
}


func init() {

	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(
		&contactGroupName,
		"contact-group",
		"g",
		"",
		`Name of the Nodeping Contact Group to retrieve uptime data for.`,
	)
	runCmd.Flags().StringVarP(
		&spreadsheetID,
		"spreadsheetID",
		"s",
		"",
		`The ID of the spreadsheet as found in its url.`,
	)
	runCmd.Flags().IntVarP(
		&countLimit,
		"count-limit",
		"l",
		0,
		`(Optional) The maximum number of results to write to Google Sheets`,
	)
}

func runArchive() {
	googlesheets.ArchiveResultsForMonth(contactGroupName, "LastMonth", spreadsheetID, nodepingToken, countLimit)
}
