# app-monitoring-archiver
Gets app monitoring results and saves them to Google Sheets

## Description
This app gets the previous month's uptime values for each Nodeping check for that are associated with
a particular contact group.

It then adds/inserts them in a Google sheet.

 - The month headings go from cell B2 to the right.
 - Each month's results are in a column starting at row 3.
 - The month columns do not get overwritten - just added to (inserted in alphabetical order).
 - The check names go from A3 down.
 - Each row has the results for one check (beginning at row 3).
 - New rows for checks are inserted in alphabetical order.  (If the existing checks are out of order,
   they will not be corrected.)


## Setup

### Google API
 - Set up a Google API project and authentication credentials using a
service account by following the instuctions at https://flaviocopes.com/google-api-authentication/
 - Give that service account edit permissions on your Google Sheet.

 Note: There is a 100 writes per 100 seconds rate limit on Google Sheets.

### Set environment variables
$ export NODEPING_TOKEN=EG123ABC
$ export SPREADSHEET_ID=EG123ABC

The SPREADSHEET_ID is the middle part of the url for the target Google Sheet when you just browse to it.

