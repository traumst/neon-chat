#echo "Dropping db file..."
#rm db/chat.db
#$la db/chat.db && echo "...Dropped db successfully." || echo "...Data not dropped."

echo "Building tailwind..."
~/code/bin/tailwindcss -i css/input.css -o css/output.css

echo "Starting server..."
go run main.go