#!/bin/sh
#
# Use this script to run your program LOCALLY.

set -e # Exit early if any commands fail

# Compile the Go program
(
  cd "$(dirname "$0")" # Ensure compile steps are run within the repository directory
  go build -o /tmp/redis-server app/*.go
)

# Run the compiled program
exec /tmp/redis-server "$@"
