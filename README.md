# app-monitoring-archiver
Gets app monitoring results and saves them to Google Sheets

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