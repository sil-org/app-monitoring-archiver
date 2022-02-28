#!/usr/bin/env bash

# Exit script with error if any step fails.
set -e

# Echo out all commands for monitoring progress
set -x

if [ "$1" = "cli" ]
then
  go build -ldflags="-s -w" -o bin/cli cmd/cli/main.go
else
  go build -ldflags="-s -w" -o bin/lambda cmd/lambda/main.go
fi
