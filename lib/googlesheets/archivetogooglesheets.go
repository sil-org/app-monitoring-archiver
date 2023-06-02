package googlesheets

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/silinternational/nodeping-cli/lib"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/sheets/v4"
)

const MonthHeaderRow = 2

type SheetsData struct {
	SpreadsheetID string // The ID of the whole Google Sheets file
	SheetID       int64  // The index of the individual sheet
	Service       *sheets.Service
}

func EnsureSheetExists(sheetName string, sheetsData SheetsData) (int64, error) {
	doesSheetExist, sheetID, err := GetSheetIDFromTitle(sheetName, sheetsData)
	if err != nil {
		return 0, err
	}

	if !doesSheetExist {
		request := sheets.Request{
			AddSheet: &sheets.AddSheetRequest{
				Properties: &sheets.SheetProperties{
					Title: sheetName,
				},
			},
		}

		rbb := &sheets.BatchUpdateSpreadsheetRequest{
			Requests: []*sheets.Request{&request},
		}

		spreadsheetID := sheetsData.SpreadsheetID
		srv := sheetsData.Service
		_, err := srv.Spreadsheets.BatchUpdate(spreadsheetID, rbb).Context(context.Background()).Do()
		if err != nil {
			return 0, fmt.Errorf("Unable to create new sheet %s. %s", sheetName, err)
		}

		_ = WriteToCellWithColumnLetter(1, "B", "Uptime Percent", sheetName, spreadsheetID, srv)
		_ = WriteToCellWithColumnLetter(2, "A", "Checks", sheetName, spreadsheetID, srv)
	}

	doesSheetExist, sheetID, err = GetSheetIDFromTitle(sheetName, sheetsData)
	if err != nil {
		return 0, fmt.Errorf("Error finding newly created sheet %s. %v", sheetName, err)
	}

	if !doesSheetExist {
		return 0, fmt.Errorf("Unable to find newly created sheet %s.", sheetName)
	}

	return sheetID, nil
}

func EnsureMonthColumnExists(month, year string, sheetsData SheetsData) (int, error) {
	desiredMonthPosition, err := GetMonthPosition(month)
	monthHeader := fmt.Sprintf("%s %s", month, year)

	monthsRange := fmt.Sprintf("%s!B2:Z2", year)
	srv := sheetsData.Service
	spreadsheetID := sheetsData.SpreadsheetID
	sheetID := sheetsData.SheetID

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, monthsRange).Do()
	if err != nil {
		return 0, fmt.Errorf("Error getting month headings for %s: %s", monthsRange, err)
	}

	indexOfFirstMonth := 1

	// No Month Heading in first results column, so just use that column
	if len(resp.Values) < 1 {
		err := WriteToCellWithColumnIndex(MonthHeaderRow, int64(indexOfFirstMonth), monthHeader, year, spreadsheetID, srv)
		return indexOfFirstMonth, err
	}

	chosenColumn := 0
	lastIndex := 0
	for index, value := range resp.Values[0] {
		lastIndex = index
		columnHeader := fmt.Sprintf("%v", value)

		if columnHeader == "" {
			chosenColumn = index + indexOfFirstMonth
			break
		}

		colMonthPosition, err := GetMonthPosition(columnHeader)
		if err != nil {
			continue
		}
		if desiredMonthPosition < colMonthPosition {
			chosenColumn = index + indexOfFirstMonth
			InsertColumn(int64(chosenColumn), sheetID, spreadsheetID, srv)
			err := WriteToCellWithColumnIndex(MonthHeaderRow, int64(chosenColumn), monthHeader, year, spreadsheetID, srv)
			return chosenColumn, err
		}
	}

	if chosenColumn == 0 {
		AddColumn(sheetID, spreadsheetID, srv)
		chosenColumn = lastIndex + indexOfFirstMonth + 1
	}

	err = WriteToCellWithColumnIndex(MonthHeaderRow, int64(chosenColumn), monthHeader, year, spreadsheetID, srv)
	return chosenColumn, err
}

// EnsureCheckRowExists looks for a match for the check name in the Sheet's A column (starting at row 3)
// If it finds a match or a blank cell, it returns that row number.  Otherwise, it looks down the column
// until it finds an existing check name that comes after it in terms of alphabetical order.
// Once it finds such an existing check name, it inserts a row above the existing row and then
// inserts the new check name into the first cell of the inserted row.
func EnsureCheckRowExists(nodepingCheck, year string, sheetsData SheetsData) (int, error) {
	checksRange := fmt.Sprintf("%s!A3:A100", year)
	srv := sheetsData.Service
	spreadsheetID := sheetsData.SpreadsheetID
	sheetID := sheetsData.SheetID

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, checksRange).Do()
	if err != nil {
		return 0, fmt.Errorf("Error getting Nodeping Check names for %s: %s", nodepingCheck, err)
	}

	indexOfFirstCheck := 3

	npCheckLower := strings.ToLower(nodepingCheck)

	chosenRow := 0
	for index, cells := range resp.Values {

		// There should always be a result for a cell, even if it's empty,
		// But, just to be careful, if it does happen, then use that row for this Check name
		if len(cells) < 1 {
			chosenRow = index + indexOfFirstCheck
			break
		}

		rowCheckName := fmt.Sprintf("%v", cells[0])

		// If the cell is blank, then use that row for this Check name
		if rowCheckName == "" {
			chosenRow = index + indexOfFirstCheck
			break
		}

		// If the new Check name comes before (alphabetically, case-insensitive) than
		// the existing Check name in the current row, then
		// insert a row and use the newly inserted row
		if npCheckLower < strings.ToLower(rowCheckName) {
			chosenRow = index + indexOfFirstCheck - 1 // It must be doing an "insert below"
			InsertRow(int64(chosenRow), sheetID, spreadsheetID, srv)
			chosenRow += 1 // Reverse the -1 from two code lines up
			break

			// If the new Check name matches the existing check name in the current row (case-insensitive)
			// then use the current row
		} else if npCheckLower == strings.ToLower(rowCheckName) {
			chosenRow = index + indexOfFirstCheck
			return chosenRow, nil
		}
	}

	// Just to be careful, if somehow there weren't any results returned for the range, then
	// use the first row of the range.
	if chosenRow == 0 {
		chosenRow = indexOfFirstCheck
	}

	err = WriteToCellWithColumnLetter(int64(chosenRow), "A", nodepingCheck, year, spreadsheetID, srv)
	return chosenRow, err
}

func GetAuthConfig() *jwt.Config {
	privateKey := GetRequiredEnvVar("GOOGLE_AUTH_PRIVATE_KEY")
	privateKey = strings.Replace(privateKey, "\\n", "\n", -1)

	config := &jwt.Config{
		Email:        GetRequiredEnvVar("GOOGLE_AUTH_CLIENT_EMAIL"),
		PrivateKeyID: GetRequiredEnvVar("GOOGLE_AUTH_PRIVATE_KEY_ID"),
		PrivateKey:   []byte(privateKey),
		TokenURL:     GetRequiredEnvVar("GOOGLE_AUTH_TOKEN_URI"),
		Scopes:       []string{"https://www.googleapis.com/auth/spreadsheets"},
	}

	return config
}

func ArchiveResultsForMonth(contactGroupName, period, spreadsheetID, nodePingToken string, countLimit int) {
	if countLimit < 1 {
		countLimit = 1000
	}

	config := GetAuthConfig()
	client := config.Client(context.Background())

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	uptimeResults, err := lib.GetUptimesForContactGroup(nodePingToken, contactGroupName, period)
	if err != nil {
		log.Fatalf("Error getting Nodeping results.  %v", err)
	}

	// Get the human readable form of the month and year
	monthTime := uptimeResults.StartTime + 86400 // Add seconds per day to ensure time zone issues don't point to previous month
	month := time.Unix(monthTime, 0).Format("January")
	year := time.Unix(monthTime, 0).Format("2006")

	sheetsData := SheetsData{
		SpreadsheetID: spreadsheetID,
		Service:       srv,
	}

	sheetID, err := EnsureSheetExists(year, sheetsData)
	if err != nil {
		log.Fatal(err.Error())
	}

	sheetsData.SheetID = sheetID

	monthColumn, err := EnsureMonthColumnExists(month, year, sheetsData)
	if err != nil {
		log.Fatalf("Error choosing column for %s.  %v", month, err)
	}

	index := 1
	delaySeconds := time.Duration(22)

	for nodepingCheck, percentage := range uptimeResults.Uptimes {
		if index > countLimit {
			break
		}

		// The quota is 100 writes per 100 seconds per user
		if index%20 == 0 {
			fmt.Printf("Waiting %v seconds at index %d to avoid Google Api rate limiting.\n", delaySeconds.Seconds(), index)
			time.Sleep(time.Second * delaySeconds)
		}

		checkRow, err := EnsureCheckRowExists(nodepingCheck, year, sheetsData)
		if err != nil {
			log.Fatalf("Error adding row for %s", nodepingCheck)
		}

		err = WriteToCellWithColumnIndex(
			int64(checkRow), int64(monthColumn),
			fmt.Sprintf("%.3f", percentage), year,
			spreadsheetID, srv,
		)

		index += 1
	}
}
