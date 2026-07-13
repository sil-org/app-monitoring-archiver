package cmd

import (
	"log/slog"
	"os"

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
			slog.Error("required flag is missing", "flag", "contact-group", "example", `-g "AppsDev Alerts"`)
			os.Exit(1)
		}

		if spreadsheetID == "" {
			slog.Error("required flag is missing", "flag", "spreadsheetID")
			os.Exit(1)
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
	err := googlesheets.ArchiveResultsForMonth(contactGroupName, "LastMonth", spreadsheetID, nodePingToken, countLimit)
	if err != nil {
		slog.Error("archive failed", "error", err)
		os.Exit(1)
	}
}
