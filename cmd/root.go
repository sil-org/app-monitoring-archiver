package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os"
)

const NodepingTokenKey = "NODEPING_TOKEN"
var nodepingToken string

var rootCmd = &cobra.Command{
	Use: "app-monitoring-archiver",
	Short: "Write Nodeping uptime results to Google Sheets",
	Long:  `Script for getting uptime results from Nodeping for a certain contact group for the previous month and saving them to Google Sheets.`,
}


func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}

func init() {

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.terraform-enterprise-migrator.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")


	// Get Nodeping Token from env vars
	nodepingToken = os.Getenv(NodepingTokenKey)

	if nodepingToken == "" {
		log.Fatal("Error: Environment variable for NODEPING_TOKEN is required to execute plan and migration \n")
	}
}