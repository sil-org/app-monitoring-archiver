package googlesheets

import (
	"fmt"
	"io/ioutil"
	"log"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"strings"
	"github.com/silinternational/nodeping-cli/lib"
)


const CredentialsForGoogle = `./lib/googlesheets/auth.json`
const MonthHeaderRow = 2



func EnsureSheetExists(spreadsheetID, title string, srv *sheets.Service) (int64, error) {
	doesSheetExist, sheetID, err := GetSheetIDFromTitle(spreadsheetID, title, srv)
	if err != nil {
		return 0, err
	}

	if !doesSheetExist {
		request := sheets.Request{
			AddSheet: &sheets.AddSheetRequest{
				Properties: &sheets.SheetProperties{
					Title: title,
				},
			},
		}

		rbb := &sheets.BatchUpdateSpreadsheetRequest{
			Requests: []*sheets.Request{&request},
		}
		_, err := srv.Spreadsheets.BatchUpdate(spreadsheetID, rbb).Context(context.Background()).Do()
		if err != nil {
			return 0, fmt.Errorf("Unable to create new sheet %s. %s", title, err)
		}
	}

	doesSheetExist, sheetID, err = GetSheetIDFromTitle(spreadsheetID, title, srv)
	if err != nil {
		return 0, fmt.Errorf("Error finding newly created sheet %s. %v", title, err)
	}

	if ! doesSheetExist {
		return 0, fmt.Errorf("Unable to find newly created sheet %s.", title)
	}

	return sheetID, nil
}


func EnsureMonthColumnExists(sheetID int64, month, year, spreadsheetID string, srv *sheets.Service) (int, error) {

	desiredMonthPosition, err := GetMonthPosition(month)
	monthHeader := fmt.Sprintf("%s %s", month, year)

	monthsRange := fmt.Sprintf("%s!B2:Z2", year)
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


func EnsureCheckRowExists(sheetID int64, nodepingCheck, year, spreadsheetID string, srv *sheets.Service) (int, error) {
	checksRange := fmt.Sprintf("%s!A3:A100", year)
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, checksRange).Do()
	if err != nil {
		return 0, fmt.Errorf("Error getting Nodeping Check names for %s: %s", nodepingCheck, err)
	}


	indexOfFirstCheck := 3

	// No Check Heading in first row, so just use that row
	if len(resp.Values) < 1 {
		err = WriteToCellWithColumnLetter(int64(indexOfFirstCheck), "A", nodepingCheck, year, spreadsheetID, srv)
		return indexOfFirstCheck, err
	}

	npCheckLower := strings.ToLower(nodepingCheck)

	chosenRow := 0
	lastIndex := 0

	for index, value := range resp.Values {
		lastIndex = index

		if len(value) < 1 {
			chosenRow = index + indexOfFirstCheck
			break
		}

		rowCheckName := fmt.Sprintf("%v", value[0])

		if npCheckLower < strings.ToLower(rowCheckName) {
			chosenRow = index + indexOfFirstCheck - 1  // It must be doing an "insert below"
			InsertRow(int64(chosenRow), sheetID, spreadsheetID, srv)
			chosenRow += 1
			err := WriteToCellWithColumnLetter(int64(chosenRow), "A", nodepingCheck, year, spreadsheetID, srv)
			return chosenRow, err
		} else if npCheckLower == strings.ToLower(rowCheckName) {
			chosenRow = index +indexOfFirstCheck
			return chosenRow, nil
		}
	}

	if chosenRow == 0 {
		AddColumn(sheetID, spreadsheetID, srv)
		chosenRow = lastIndex + indexOfFirstCheck + 1
	}

	err = WriteToCellWithColumnLetter(int64(chosenRow), "A", nodepingCheck, year, spreadsheetID, srv)
	return chosenRow, err
}

func ArchiveResultsForMonth(contactGroupName, month, year, spreadsheetID, nodePingToken string) {
	credBytes, err := ioutil.ReadFile(CredentialsForGoogle)
	if err != nil {
		log.Fatalf("Unable to read google credentials file: %v", err)
	}

	config, err := google.JWTConfigFromJSON(credBytes, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse google credentials file to config: %v", err)
	}
	client := config.Client(context.Background())

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	sheetID, err := EnsureSheetExists(spreadsheetID, year, srv)
	if err != nil {
		log.Fatal(err.Error())
	}


	uptimeResults, err := lib.GetUptimesForContactGroup(nodePingToken, contactGroupName, "LastMonth")
	if err != nil {
		log.Fatalf("Error getting Nodeping results.  %v", err)
	}

	monthColumn, err := EnsureMonthColumnExists(sheetID, month, year, spreadsheetID, srv)
	if err != nil {
		log.Fatalf("Error choosing column for %s.  %v", month, err)
	}

	index := 0
	for nodepingCheck, percentage := range uptimeResults.Uptimes {
		checkRow, err := EnsureCheckRowExists(sheetID, nodepingCheck, year, spreadsheetID, srv)
		if err != nil {
			log.Fatalf("Error adding row for %s", nodepingCheck)
		}

		err = WriteToCellWithColumnIndex(int64(checkRow), int64(monthColumn), fmt.Sprintf("%.3f", percentage), year, spreadsheetID, srv)

		if index > 2 {
			break
		}
		index += 1
	}
}
