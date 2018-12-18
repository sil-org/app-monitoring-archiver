#!/usr/bin/env bash

# Exit if any command below fails
set -e
set -x

# build binary
go build -ldflags="-s -w" -o archive-to-google-sheets

STAGE="dev"
if [[ "${CI_BRANCH}" == "master" ]]; then
    STAGE="prod"
fi

echo "Deploying app-monitoring-archiver ..."
serverless deploy -v --stage $STAGE
