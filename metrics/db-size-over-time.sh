#!/bin/bash
# takes some db-related metrics like file size, table count, total rows count, etc
# calculates some agregates
# appends new line to db-size-over-time.csv

DATABASE_FILE="./chat.db"
if [ ! -f "$DATABASE_FILE" ]; then
  echo "'$DATABASE_FILE' does not exist - metrics cant be collected"
  exit 1
fi
echo "'$DATABASE_FILE' exists..."

CSV_FILE="./metrics/db-size-over-time.csv"
if [ -s "$CSV_FILE" ]; then
  LAST_LINE=$(tail -n 1 "$CSV_FILE")
  LAST_NEXT=$(echo "$LAST_LINE" | cut -d',' -f1)
  NEXT=$((LAST_NEXT + 1))
  echo "next line number will be $NEXT"
else
  CSV_FILE_HEADER="day,timestamp,filesize,table_count,index_count,total_rows_count"
  echo "$CSV_FILE_HEADER" >> "$CSV_FILE"
  NEXT=1
  echo "csv header appended, initial line number is $NEXT"
fi

TABLES_COUNT=0
ROWS_COUNT=0
TABLE_NAMES=$(sqlite3 "$DATABASE_FILE" "SELECT name FROM sqlite_master WHERE type='table';")
echo "counting rows in following tables: [$TABLE_NAMES]"
for TABLE in $TABLE_NAMES; do
  ROW_COUNT=$(sqlite3 "$DATABASE_FILE" "SELECT COUNT(1) FROM $TABLE;")
  ROWS_COUNT=$((TOTAL_TABLE_ROWs_COUNT + ROW_COUNT))
  TABLES_COUNT=$((TABLES_COUNT + 1))
done
echo "counted $ROWS_COUNT rows over $TABLES_COUNT tables"

INDEX_COUNT=$(sqlite3 "$DATABASE_FILE" "SELECT count(1) FROM sqlite_master WHERE type='index';")

FILE_SIZE=$(stat -f%z "$DATABASE_FILE")

AVG_ROW_SIZE=$((FILE_SIZE / ROWS_COUNT))
AVG_CUSTOM_LOOKUP_777=$((ROWS_COUNT / ((TABLES_COUNT + INDEX_COUNT))))

CSV_FILE_ROW="$NEXT,$(date -u +"%Y-%m-%dT%H:%M:%SZ"),$FILE_SIZE,$TABLES_COUNT,$INDEX_COUNT,$ROWS_COUNT"
echo "$CSV_FILE_ROW" >> "$CSV_FILE"
echo "Data appended to $CSV_FILE"