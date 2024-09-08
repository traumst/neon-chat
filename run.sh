echo "Running tests"
time go test ./...

# creating db backup
# BACKUP_EXIT_CODE=
# if [ $(./backup-db.sh -in chat.db -out ./backup) -eq 0 ]; then
#   echo "Database backed up"
# else
#   echo "Database backup failed"
#   exit 1
# fi
# echo "Dropping db file..."
# rm chat.db
# (la chat.db && echo "...Dropped db successfully.") || echo "...Data not dropped."


echo "logging metrics..."
./metrics/db-size-over-time.sh

echo "Building tailwind..."
~/code/bin/tailwindcss -i static/css/input.css -o static/css/tailwind.css

echo "Starting server..."
go run main.go