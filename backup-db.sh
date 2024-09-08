#!/bin/bash
# call it like `./backup-db.sh -in chat.db -out ./backup`
# creates ./backup folder if needed
# dumps chat.db to ./backup/chat.db-backup-2024-13-35T76:90:65Z.sql

if [ "$1" != "-in" ] || [ "$3" != "-out" ]; then
  echo "required format: $0 -in <sqlite_db_file> -out <backup_folder>"
  echo "you provided:    $0 $1 $2 $3 $4"
  exit 1
fi
DATABASE_FILE=$2
BACKUP_FOLDER_RELATIVE_PATH=$4
#DATABASE_FILE="chat.db"
#BACKUP_FOLDER_RELATIVE_PATH="./backup"

ISO_TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
BACKUP_FILE="$BACKUP_FOLDER_RELATIVE_PATH/$DATABASE_FILE-backup-$ISO_TIMESTAMP.sql"

if [ ! -d "$BACKUP_FOLDER_RELATIVE_PATH" ]; then
  mkdir "$BACKUP_FOLDER_RELATIVE_PATH"
  echo "'$BACKUP_FOLDER_RELATIVE_PATH' folder created!"
else
  echo "'$BACKUP_FOLDER_RELATIVE_PATH' folder exists..."
fi

sqlite3 "$DATABASE_FILE" .dump > "$BACKUP_FILE"
if [ $? -eq 0 ]; then
  echo "Database backup successful: '$BACKUP_FILE'"
else
  echo "Database backup failed."
fi