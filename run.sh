#echo "Dropping db file..."
rm db/chat.db
#$la db/chat.db && echo "...Dropped db successfully." || echo "...Data not dropped."
echo "Starting server..."
go run main.go