#!/bin/bash
set -xeuo pipefail
# set -x          # Print commands before running them
# set -e          # Exit on Error
# set -u          # Exit on Uninitialized Variable
# set -o pipefail # Exit on Pipe Fail
# set -n          # Check syntax without running

echo "Performing checks and cleanup..."
go clean && \
go mod tidy && \
go vet ./... && \
staticcheck ./...

echo "Running tests..."
time go test ./...
echo "Logging metrics..."
./metrics/db-size-over-time.sh
# WIP
#./metrics/project-size-over-time.sh

# echo "Backing up db file..."
# if [ $(./backup-db.sh -in chat.db -out ./backup) -eq 0 ]; then
#   echo "Database backed up"
# else
#   echo "Database backup failed"
#   exit 1
# fi
# if [ -f chat.db ]; then 
#     echo "Dropping db file..."
#     rm chat.db
#     (ls chat.db && echo "...Dropped db successfully.") || echo "...Data not dropped."
# else 
#     echo "No db file found..."
# fi

echo "Building tailwind..."
~/code/bin/tailwindcss -i static/css/input.css -o static/css/tailwind.css

echo "Starting server..."
go run main.go