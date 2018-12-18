#!/usr/bin/env bash

# Exit if any command below fails
set -e
set -x

# build binary
go build -ldflags="-s -w" -o archive-to-google-sheets

# Export env vars
export NODEPING_TOKEN=${DEV_NODEPING_TOKEN}
export CONTACT_GROUP_NAME=${DEV_CONTACT_GROUP_NAME}
export COUNT_LIMIT=${DEV_COUNT_LIMIT}
export PERIOD=${DEV_PERIOD}
export SPREADSHEET_ID=${DEV_SPREADSHEET_ID}

echo "Deploying app-monitoring-archiver ..."
serverless deploy -v --stage dev
