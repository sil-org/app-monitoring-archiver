package googlesheets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_findRowPositionAndWhetherToInsertARow(t *testing.T) {

	type testCase struct {
		name          string
		values        [][]interface{}
		checkName     string
		wantRow       int
		wantInsertRow bool
	}

	testCases := []testCase{
		{
			name:          "all empty",
			values:        [][]interface{}{},
			checkName:     "Anything",
			wantRow:       0,
			wantInsertRow: false,
		},
		{
			name:          "empty row",
			values:        [][]interface{}{[]interface{}{}},
			checkName:     "Anything",
			wantRow:       0,
			wantInsertRow: false,
		},
		{
			name:          "blank first cell",
			values:        [][]interface{}{[]interface{}{""}},
			checkName:     "Anything",
			wantRow:       0,
			wantInsertRow: false,
		},
		{
			name: "new checkName comes before first",
			values: [][]interface{}{
				[]interface{}{"Second"},
				[]interface{}{"Third"},
			},
			checkName:     "first",
			wantRow:       0,
			wantInsertRow: true,
		},
		{
			name: "new checkName is the same as first",
			values: [][]interface{}{
				[]interface{}{"First"},
				[]interface{}{"Second"},
			},
			checkName:     "first",
			wantRow:       0,
			wantInsertRow: false,
		},
		{
			name: "new checkName comes before second",
			values: [][]interface{}{
				[]interface{}{"First"},
				[]interface{}{"Third"},
			},
			checkName:     "second",
			wantRow:       1,
			wantInsertRow: true,
		},
		{
			name: "new checkName is the same as second",
			values: [][]interface{}{
				[]interface{}{"first"},
				[]interface{}{"second"},
			},
			checkName:     "Second",
			wantRow:       1,
			wantInsertRow: false,
		},
		{
			name: "new checkName comes after others",
			values: [][]interface{}{
				[]interface{}{"first"},
				[]interface{}{"second"},
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
