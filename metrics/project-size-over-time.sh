#!/bin/bash

CSV_FILE="./metrics/project-size-over-time.csv"
EXTENSIONS=("*.go" "*.html" "*.js" "*.css" "README.md")
EXTENSIONS_FIND=$(printf " -name %s -o" "${EXTENSIONS[@]}")
EXTENSIONS_FIND=${EXTENSIONS_FIND% -o}

if [ ! -f "$CSV_FILE" ]; then
  echo "day,timestamp,total_files,largest_folder_name,largest_folder_size,greatest_file_depth,min_lines,min_file,max_lines,max_file,median_lines,median_file,total_non_empty_lines" > "$CSV_FILE"
fi

DAY=$(date +%Y-%m-%d)
TIMESTAMP=$(date +%s)

# Find all relevant files and folders
FOLDERS=()
FILES=()
for file in $(find . ! -path "./.git" ! -path "./.git/*" ! -path "./log" ! -path "./log/*")
do
  echo "$file"
  if [ -d "$file" ]; then
    FOLDERS+=("$file")
  elif [ -f "$file" ]; then
    FILES+=("$file")
  fi
done

declare -A FILE_LENGTHS

total_non_empty_lines=0
min_lines=-1
max_lines=-1
min_file=""
max_file=""
file_count=0
declare -a LINE_COUNTS_FILES

for file in $FILES; do
  LINES=$(wc -l < "$file")
  FILE_LENGTHS["$file"]=$LINES
  total_non_empty_lines=$((total_non_empty_lines + LINES))
  LINE_COUNTS_FILES+=("$LINES:$file")
  
  if [ $min_lines -eq -1 ] || [ $LINES -lt $min_lines ]; then
    min_lines=$LINES
    min_file="$file"
  fi
  if [ $max_lines -eq -1 ] || [ $LINES -gt $max_lines ]; then
    max_lines=$LINES
    max_file="$file"
  fi
  file_count=$((file_count + 1))
done


exit 2222

# Total number of relevant files
TOTAL_FOLDERS=$(find "." -type d | wc -l)
TOTAL_FILES=$(find . -type f \( $EXTENSIONS_FIND \) | wc -l)
echo "There are $TOTAL_FILES dirs and files in the project."

exit 111

# Greatest file depth
GREATEST_FILE_DEPTH=$(find "." -type f | awk -F"/" '{print NF-1}' | sort -nr | head -n 1)
echo "The greatest file depth is $GREATEST_FILE_DEPTH."

# Dir with most files
LARGEST_FOLDER=$(find "." -type f \( $EXTENSIONS_FIND \) -exec dirname {} \; | sort | uniq -c | sort -nr | head -n 1)
LARGEST_FOLDER_SIZE=$(echo "$LARGEST_FOLDER" | awk '{print $1}')
LARGEST_FOLDER_NAME=$(echo "$LARGEST_FOLDER" | awk '{print $2}')
echo "Largest folder is $LARGEST_FOLDER_NAME with $LARGEST_FOLDER_SIZE files."

# Total lines count
LINE_COUNTS=()
FILE_NAMES=()
for EXT in "${EXTENSIONS[@]}"; do
  while IFS= read -r -d '' FILE; do
    LINE_COUNT=$(wc -l < "$FILE")
    LINE_COUNTS+=("$LINE_COUNT")
    FILE_NAMES+=("$FILE")
  done < <(find "." -name "$EXT" -print0)
done

TOTAL_NON_EMPTY_LINES=$(printf '%s\n' "${LINE_COUNTS[@]}" | paste -sd+ - | bc)
echo "There are $TOTAL_NON_EMPTY_LINES non-empty lines of code in this project."

MIN_INDEX=$(printf '%s\n' "${!LINE_COUNTS[@]}" | sort -n -k1,1 -k2,2 | head -n 1)
MIN_LINES=${LINE_COUNTS[$MIN_INDEX]}
MIN_FILE=${FILE_NAMES[$MIN_INDEX]}
echo "The smallest file is $MIN_FILE with $MIN_LINES lines."

MAX_INDEX=$(printf '%s\n' "${!LINE_COUNTS[@]}" | sort -n -k1,1 -k2,2 | tail -n 1)
MAX_LINES=${LINE_COUNTS[$MAX_INDEX]}
MAX_FILE=${FILE_NAMES[$MAX_INDEX]}
echo "The largest file is $MAX_FILE with $MAX_LINES lines."

MEDIAN_INDEX=$(printf '%s\n' "${!LINE_COUNTS[@]}" | sort -n -k1,1 -k2,2 | {
  read -r first
  read -r second
  if (( ${#LINE_COUNTS[@]} % 2 == 1 )); then
    for ((i = 2; i <= ${#LINE_COUNTS[@]} / 2; i++)); do
      read -r first
    done
    echo "$first"
  else
    for ((i = 2; i < ${#LINE_COUNTS[@]} / 2; i++)); do
      read -r first
    done
    read -r second
    echo "($first + $second) / 2" | bc
  fi
})
MEDIAN_LINES=${LINE_COUNTS[$MEDIAN_INDEX]}
MEDIAN_FILE=${FILE_NAMES[$MEDIAN_INDEX]}
echo "The median file is $MEDIAN_FILE with $MEDIAN_LINES lines."

exit 666

echo "$DAY,$TIMESTAMP,$TOTAL_FILES,$LARGEST_FOLDER_NAME,$LARGEST_FOLDER_SIZE,$GREATEST_FILE_DEPTH,$MIN_LINES,$MIN_FILE,$MAX_LINES,$MAX_FILE,$MEDIAN_LINES,$MEDIAN_FILE,$TOTAL_NON_EMPTY_LINES" >> "$CSV_FILE"