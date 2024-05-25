echo "Running tests"
go test ./...

#echo "Dropping db file..."
rm chat.db
(la chat.db && echo "...Dropped db successfully.") || echo "...Data not dropped."

echo "Building tailwind..."
~/code/bin/tailwindcss -i static/css/input.css -o static/css/tailwind.css

echo "Starting server..."
go run main.go