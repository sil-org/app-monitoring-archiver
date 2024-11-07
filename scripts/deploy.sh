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

apt-get update && apt-get -y install curl
curl -o- --location https://slss.io/install | VERSION=$SERVERLESS_VERSION bash
mv $HOME/.serverless/bin/serverless /usr/local/bin
ln -s /usr/local/bin/serverless /usr/local/bin/sls

echo "Deploying app-monitoring-archiver as lambda (stage=$1)..."
serverless deploy --verbose --stage "$1"
