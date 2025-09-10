package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/sil-org/app-monitoring-archiver/lib/googlesheets"
)

var (
	contactGroupName string
	spreadsheetID    string
	countLimit       int
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Archive NodePing results",
	Long:  "Get the uptime results from NodePing and write them to Google Sheets.",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if contactGroupName == "" {
			log.Fatal(`Error: The 'contact-group' flag is required (e.g. -g "AppsDev Alerts").`)
		}

		if spreadsheetID == "" {
			log.Fatal(`Error: The 'spreadsheetID' flag is required (found in its url).`)
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
		`Name of the NodePing Contact Group to retrieve uptime data for.`,
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
	googlesheets.ArchiveResultsForMonth(contactGroupName, "LastMonth", spreadsheetID, nodePingToken, countLimit)
}
