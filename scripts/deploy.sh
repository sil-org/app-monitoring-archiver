#!/usr/bin/env bash

# Exit if any command below fails
set -e

if [ "$1" != "prod" ] && [ "$1" != "dev" ]; then
  echo "invalid stage $1"
  exit 1
fi

# Build binaries
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
"$DIR"/build.sh

echo "Deploying app-monitoring-archiver as lambda (stage=$1)..."
serverless deploy --verbose --stage "$1"
