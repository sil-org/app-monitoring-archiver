#!/usr/bin/env bash

# Echo out all commands for monitoring progress
set -x

# When using the provided.al2 runtime, the binary must be named "bootstrap" and be in the root directory
if [ "$1" = "cli" ]; then
  CGO_ENABLED=0 go build -tags lambda.norpc -ldflags="-s -w" -o bootstrap cmd/cli/main.go
else
  CGO_ENABLED=0 go build -tags lambda.norpc -ldflags="-s -w" -o bootstrap cmd/lambda/main.go
fi
