package googlesheets

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/sil-org/nodeping-cli"
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
			return 0, fmt.Errorf("unable to create new sheet %s. %s", sheetName, err)
		}

		_ = WriteToCellWithColumnLetter(1, "B", "Uptime Percent", sheetName, spreadsheetID, srv)
		_ = WriteToCellWithColumnLetter(2, "A", "Checks", sheetName, spreadsheetID, srv)
	}

	doesSheetExist, sheetID, err = GetSheetIDFromTitle(sheetName, sheetsData)
	if err != nil {
		return 0, fmt.Errorf("error finding newly created sheet %s. %w", sheetName, err)
	}

	if !doesSheetExist {
		return 0, fmt.Errorf("unable to find newly created sheet %s.", sheetName)
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
		return 0, fmt.Errorf("error getting month headings for %s: %w", monthsRange, err)
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
			if err := InsertColumn(int64(chosenColumn), sheetID, spreadsheetID, srv); err != nil {
				return 0, fmt.Errorf("error inserting column in Google Sheets. %w", err)
			}
			err := WriteToCellWithColumnIndex(MonthHeaderRow, int64(chosenColumn), monthHeader, year, spreadsheetID, srv)
			return chosenColumn, err
		}
	}

	if chosenColumn == 0 {
		if err := AddColumn(sheetID, spreadsheetID, srv); err != nil {
			return 0, err
		}
		chosenColumn = lastIndex + indexOfFirstMonth + 1
	}

	err = WriteToCellWithColumnIndex(MonthHeaderRow, int64(chosenColumn), monthHeader, year, spreadsheetID, srv)
	return chosenColumn, err
}

// This returns a row number (0-indexed) and a boolean as to whether a row needs to be inserted.
//
//	It begins by converting the first column of the rows to string values.
//	If there are no cells in that column, it just returns (0, false).
//	If it comes to a cell that is empty or matches the checkName value, it returns the
//	  corresponding row number and false.
//	If it comes to a cell that has a value greater (alphabetically) than the checkName, it
//	  returns the corresponding row number and true (i.e. a row needs to be inserted)
func findRowPositionAndWhetherToInsertARow(checkName string, rows [][]any) (int, bool) {
	if len(rows) == 0 || len(rows[0]) == 0 {
		return 0, false
	}
	checkName = strings.ToLower(checkName)

	rowCount := 0

	for i, cells := range rows {
		// This should never happen, but let's just be extra careful
		if len(cells) == 0 {
			return i, false
		}

		value := strings.ToLower(fmt.Sprintf("%v", cells[0]))
		if value == "" || value == checkName {
			return i, false
		}
		if value > checkName {
			return i, true
		}
		rowCount++
	}

	return rowCount, false
}

// EnsureCheckRowExists looks for a match for the check name in the Sheet's A column (starting at row 3)
// If it finds a match or a blank cell, it returns that row number.  Otherwise, it looks down the column
// until it finds an existing check name that comes after it in terms of alphabetical order.
// Once it finds such an existing check name, it inserts a row above the existing row and then
// inserts the new check name into the first cell of the inserted row.
func EnsureCheckRowExists(nodePingCheck, year string, sheetsData SheetsData) (int, error) {
	const indexOfFirstCheck = 3
	checksRange := fmt.Sprintf("%s!A%d:A100", year, indexOfFirstCheck)
	srv := sheetsData.Service
	spreadsheetID := sheetsData.SpreadsheetID
	sheetID := sheetsData.SheetID

	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, checksRange).Do()
	if err != nil {
		return 0, fmt.Errorf("error getting NodePing Check names for %s: %w", nodePingCheck, err)
	}

	rowInRange, insertRow := findRowPositionAndWhetherToInsertARow(nodePingCheck, resp.Values)
	chosenRow := rowInRange + indexOfFirstCheck

	if insertRow {
		row := chosenRow - 1 // It must be doing an "insert below"
		log.Printf("Inserting row above row %d for NodePing check %s", chosenRow, nodePingCheck)
		if err := InsertRow(int64(row), sheetID, spreadsheetID, srv); err != nil {
			return 0, fmt.Errorf("error inserting a row in Google sheets: %w", err)
		}
	}

	err = WriteToCellWithColumnLetter(int64(chosenRow), "A", nodePingCheck, year, spreadsheetID, srv)
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

func ArchiveResultsForMonth(contactGroupName, period, spreadsheetID, nodePingToken string, countLimit int) error {
	if countLimit < 1 {
		countLimit = 1000
	}

	config := GetAuthConfig()
	client := config.Client(context.Background())

	srv, err := sheets.New(client)
	if err != nil {
		return fmt.Errorf("unable to retrieve Sheets client: %w", err)
	}

	p, err := nodeping.GetPeriod(period)
	if err != nil {
		return fmt.Errorf("error getting NodePing period: %w", err)
	}

	uptimeResults, err := nodeping.GetUptimesForContactGroup(nodePingToken, contactGroupName, *p)
	if err != nil {
		return fmt.Errorf("error getting NodePing results: %w", err)
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
		return err
	}

	sheetsData.SheetID = sheetID

	monthColumn, err := EnsureMonthColumnExists(month, year, sheetsData)
	if err != nil {
		return fmt.Errorf("error choosing column for '%s': %w", month, err)
	}

	index := 1
	const delaySeconds = time.Second * 22

	for nodePingCheck, percentage := range uptimeResults.Uptimes {
		if index > countLimit {
			break
		}

		// The quota is 100 writes per 100 seconds per user
		if index%20 == 0 {
			fmt.Printf("Waiting %v seconds at index %d to avoid Google Api rate limiting.\n", delaySeconds.Seconds(), index)
			time.Sleep(delaySeconds)
		}

		checkRow, err := EnsureCheckRowExists(nodePingCheck, year, sheetsData)
		if err != nil {
			return fmt.Errorf("error adding row for '%s'", nodePingCheck)
		}

		err = WriteToCellWithColumnIndex(
			int64(checkRow), int64(monthColumn),
			fmt.Sprintf("%.3f", percentage), year,
			spreadsheetID, srv,
		)

		index += 1
	}
	return nil
}
