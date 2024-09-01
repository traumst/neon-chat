#!/bin/bash
# creates ./backup folder if needed
# dumps chat.db to ./backup/chat.db-backup-2024-13-35T76:90:65Z.sql

ISO_TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
BACKUP_FOLDER_RELATIVE_PATH="backup"
DATABASE_FILE="chat.db"
BACKUP_FILE="./$BACKUP_FOLDER_RELATIVE_PATH/$DATABASE_FILE-backup-$ISO_TIMESTAMP.sql"

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