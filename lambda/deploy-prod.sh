#!/usr/bin/env bash

# Exit if any command below fails
set -e
set -x

# build binary
go build -ldflags="-s -w" -o archive-to-google-sheets

# Export env vars
export NODEPING_TOKEN=${PROD_NODEPING_TOKEN}
export CONTACT_GROUP_NAME=${PROD_CONTACT_GROUP_NAME}
export COUNT_LIMIT=${PROD_COUNT_LIMIT}
export PERIOD=${PROD_PERIOD}
export SPREADSHEET_ID=${PROD_SPREADSHEET_ID}

export GOOGLE_AUTH_CLIENT_EMAIL=${PROD_GOOGLE_AUTH_CLIENT_EMAIL}
export GOOGLE_AUTH_PRIVATE_KEY_ID=${PROD_GOOGLE_AUTH_PRIVATE_KEY_ID}
export GOOGLE_AUTH_PRIVATE_KEY=${PROD_GOOGLE_AUTH_PRIVATE_KEY}
export GOOGLE_AUTH_TOKEN_URI=${PROD_GOOGLE_AUTH_TOKEN_URI}

echo "Deploying app-monitoring-archiver as archive-to-google-sheets ..."
serverless deploy -v --stage prod
