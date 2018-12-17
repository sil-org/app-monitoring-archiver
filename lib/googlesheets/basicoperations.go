package googlesheets

import (
	"strings"
	"fmt"
	"google.golang.org/api/sheets/v4"
	"golang.org/x/net/context"
	"log"
)

func GetMonthPosition(monthLabel string) (int, error) {
	monthLabel = strings.Trim(monthLabel, " ")
	monthParts := strings.Split(monthLabel, " ")
	month := strings.ToLower(monthParts[0])

	indexes := map[string]int{
		"january": 1,
		"february": 2,
		"march": 3,
		"april": 4,
		"may": 5,
		"june": 6,
		"july": 7,
		"august": 8,
		"september": 9,
		"october": 10,
		"november": 11,
		"december": 12,
	}

	index, ok := indexes[month]
	if ! ok {
		return 0, fmt.Errorf("Month %s not valid", monthLabel)
	}

	return index, nil
}

func ConvertColumnIndexToLetter(index int64) (string, error) {

	if index > 25 {
		return "", fmt.Errorf("Not allowed to convert index if over 25. It was %d", index)
	}

	runeA := int64([]rune("A")[0])

	return fmt.Sprintf("%c", index + runeA), nil
}


func WriteToCellWithColumnLetter(rowIndex int64, columnLetter, newValue, sheetName, spreadsheetID string, srv *sheets.Service) error {
	cellRange := fmt.Sprintf("%s!%s%d", sheetName, columnLetter, rowIndex)
	valueRange := &sheets.ValueRange{}

	updateValue := []interface{}{newValue}
	valueRange.Values = append(valueRange.Values, updateValue)

	_, err := srv.Spreadsheets.Values.Update(spreadsheetID, cellRange, valueRange).ValueInputOption("RAW").Do()
	if err != nil {
		return fmt.Errorf("Unable to write to cell %s. %v", cellRange, err)
	}

	return nil
}


func WriteToCellWithColumnIndex(rowIndex, columnIndex int64, newValue, sheetName, spreadsheetID string, srv *sheets.Service) error {
	columnLetter, err := ConvertColumnIndexToLetter(columnIndex)
	if err != nil {
		return err
	}
	return WriteToCellWithColumnLetter(rowIndex, columnLetter, newValue, sheetName, spreadsheetID, srv)
}


func InsertRowOrColumn(insertRowNotColumn bool, index, sheetID int64, spreadsheetID string, srv *sheets.Service) error {
	dimension := "COLUMNS"
	if insertRowNotColumn {
		dimension = "ROWS"
	}

	request := sheets.Request{
		InsertDimension: &sheets.InsertDimensionRequest {
			Range: &sheets.DimensionRange {
				SheetId: sheetID,
				Dimension: dimension,
				StartIndex: index,
				EndIndex: index+1,
			},
			InheritFromBefore: true,
		},
	}

	rbb := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{&request},
	}
	_, err := srv.Spreadsheets.BatchUpdate(spreadsheetID, rbb).Context(context.Background()).Do()
	if err != nil {
		return fmt.Errorf("Unable to insert %s %d. %s", dimension, index, err)
	}

	return nil
}


func InsertColumn(index, sheetID int64, spreadsheetID string, srv *sheets.Service) error {
	return InsertRowOrColumn(false, index, sheetID, spreadsheetID, srv)
}


func InsertRow(index, sheetID int64, spreadsheetID string, srv *sheets.Service) error {
	return InsertRowOrColumn(true, index, sheetID, spreadsheetID, srv)
}



func AddColumn(sheetID int64, spreadsheetID string, srv *sheets.Service) {
	request := sheets.Request{
		AppendDimension: &sheets.AppendDimensionRequest {
			Dimension: "COLUMNS",
			Length: 1,
			SheetId: sheetID,
		},
	}

	rbb := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{&request},
	}
	_, err := srv.Spreadsheets.BatchUpdate(spreadsheetID, rbb).Context(context.Background()).Do()
	if err != nil {
		log.Fatalf("Unable to add column to sheet %d. %s", sheetID,  err)
	}
}


func GetSheetIDFromTitle(title string, sheetsData SheetsData) (bool, int64, error) {
	ssResp, err := sheetsData.Service.Spreadsheets.Get(sheetsData.SpreadsheetID).Do()

	if err != nil {
		return false, 0, fmt.Errorf("Error trying to find sheet %s. %v", title, err)
	}

	for _, next := range ssResp.Sheets {
		if next.Properties.Title == title {
			return true, next.Properties.SheetId, nil
		}
	}

	return false, 0, nil
}