package googlesheets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_findRowPositionAndWhetherToInsertARow(t *testing.T) {
	type testCase struct {
		name          string
		values        [][]any
		checkName     string
		wantRow       int
		wantInsertRow bool
	}

	testCases := []testCase{
		{
			name:          "all empty",
			values:        [][]any{},
			checkName:     "Anything",
			wantRow:       0,
			wantInsertRow: false,
		},
		{
			name:          "empty row",
			values:        [][]any{{}},
			checkName:     "Anything",
			wantRow:       0,
			wantInsertRow: false,
		},
		{
			name:          "blank first cell",
			values:        [][]any{{""}},
			checkName:     "Anything",
			wantRow:       0,
			wantInsertRow: false,
		},
		{
			name: "new checkName comes before first",
			values: [][]any{
				{"Second"},
				{"Third"},
			},
			checkName:     "first",
			wantRow:       0,
			wantInsertRow: true,
		},
		{
			name: "new checkName is the same as first",
			values: [][]any{
				{"First"},
				{"Second"},
			},
			checkName:     "first",
			wantRow:       0,
			wantInsertRow: false,
		},
		{
			name: "new checkName comes before second",
			values: [][]any{
				{"First"},
				{"Third"},
			},
			checkName:     "second",
			wantRow:       1,
			wantInsertRow: true,
		},
		{
			name: "new checkName is the same as second",
			values: [][]any{
				{"first"},
				{"second"},
			},
			checkName:     "Second",
			wantRow:       1,
			wantInsertRow: false,
		},
		{
			name: "new checkName comes after others",
			values: [][]any{ // testing once with [][]any, since that is what the Google library uses
				{"first"},
				{"second"},
			},
			checkName:     "Third",
			wantRow:       2,
			wantInsertRow: false,
		},
	}

	for _, tc := range testCases {
		gotRow, gotInsertRow := findRowPositionAndWhetherToInsertARow(tc.checkName, tc.values)

		assert.Equal(t, tc.wantRow, gotRow, "incorrect row number in test: %s", tc.name)
		assert.Equal(t, tc.wantInsertRow, gotInsertRow, "incorrect insert row boolean in test: %s", tc.name)
	}
}
