#!/usr/bin/env bash

# Exit if any command below fails
set -e

# Build binaries
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
"$DIR"/build.sh

# Export env vars
if [ "$1" = "prod" ]; then
  export NODEPING_TOKEN=${PROD_NODEPING_TOKEN}
  export CONTACT_GROUP_NAME=${PROD_CONTACT_GROUP_NAME}
  export COUNT_LIMIT=${PROD_COUNT_LIMIT}
  export PERIOD=${PROD_PERIOD}
  export SPREADSHEET_ID=${PROD_SPREADSHEET_ID}

  export GOOGLE_AUTH_CLIENT_EMAIL=${PROD_GOOGLE_AUTH_CLIENT_EMAIL}
  export GOOGLE_AUTH_PRIVATE_KEY_ID=${PROD_GOOGLE_AUTH_PRIVATE_KEY_ID}
  export GOOGLE_AUTH_PRIVATE_KEY=${PROD_GOOGLE_AUTH_PRIVATE_KEY}
  export GOOGLE_AUTH_TOKEN_URI=${PROD_GOOGLE_AUTH_TOKEN_URI}
elif [ "$1" = "dev" ]; then
  export NODEPING_TOKEN=${DEV_NODEPING_TOKEN}
  export CONTACT_GROUP_NAME=${DEV_CONTACT_GROUP_NAME}
  export COUNT_LIMIT=${DEV_COUNT_LIMIT}
  export PERIOD=${DEV_PERIOD}
  export SPREADSHEET_ID=${DEV_SPREADSHEET_ID}

  export GOOGLE_AUTH_CLIENT_EMAIL=${DEV_GOOGLE_AUTH_CLIENT_EMAIL}
  export GOOGLE_AUTH_PRIVATE_KEY_ID=${DEV_GOOGLE_AUTH_PRIVATE_KEY_ID}
  export GOOGLE_AUTH_PRIVATE_KEY=${DEV_GOOGLE_AUTH_PRIVATE_KEY}
  export GOOGLE_AUTH_TOKEN_URI=${DEV_GOOGLE_AUTH_TOKEN_URI}
else
  echo "invalid stage $1"
  exit 1
fi

echo "Deploying app-monitoring-archiver as lambda (stage=$1)..."
serverless deploy --verbose --stage "$1"
