#!/bin/bash
#
if [[ ${#args[@]} -eq 0 ]]; then
  exit 1
fi
if [ $1 == () ]; then
  exit 2
fi
array=("$@:2")
# expects following args
#  $1            - input args array
#  args=("$@:2") - search values
contains_values() {
  local array=$1
  local lookup=("value1")
  local value2="value2"
  local found_value1=0
  local found_value2=0

  for element in "${array[@]}"; do
    if [[ "$element" == "$value1" ]]; then
      found_value1=1
    elif [[ "$element" == "$value2" ]]; then
      found_value2=1
    fi
  done

  if [[ $found_value1 -eq 1 && $found_value2 -eq 1 ]]; then
    return 0
  else
    return 1
  fi
}
is_excluded() {
  local EXCLUDES=("tailwind.css" "README.md")
  for EXCLUDE in "${EXCLUDES[@]}"; do
    if [[ "$1" == *"$EXCLUDE"* ]]; then
      echo 1
      return
    fi
  done
  echo 0
}
dirs_except_pattern() {
  find "$1" -type d \
    ! -path "./.git" \
    ! -path "./.git/*" \
    ! -path "./tmp" \
    ! -path "./tmp/*" \
    ! -path "./log" \
    ! -path "./log/*" \
    ! -path "./_bugz" \
  ;
}
files_by_pattern() {
  find "$1" -type f \
    -name "*.go" -o \
    -name "*.html" -o \
    -name "*.css" ! -name "tailwind.css" -o \
    -name "*.js" -o \
    -name "*.md" ! -name "README.md" \
  ;
}
all_except_pattern() {
  local dirs_array=($(dirs_except_pattern "$1"))
  local files_array=($(files_by_pattern "$1"))
  echo "${dirs_array[@]}" "${files_array[@]}"
}
non_empty_lines() {
  grep -v '^\s*$' "$1"
}

# Take relevant metrics
TOTAL_FILES=0
TOTAL_DIRS=0
TOTAL_LINES=0
LONGEST_FILE=""
LONGEST_LINES=0
DEEPEST_NEST=0
DEEPEST_FILE=0
declare -a DIR_SIZES; # in direct files + folders decendants
declare -a FILE_SIZES # in non-empty text lines
take_metrics() {
  #
  # mean is all numbers in the set divided by the set count. 
  # median is the middle value when a data set is ordered from least to greatest. 
  # mode is the number that occurs most often in a data set.
  # 
  echo "processing metrics at $(pwd)"
  for file in $(all_except_pattern ".")
  do
    echo "...glancing at [$file]"
    if [ -d "$file" ]; then # directory
      local NESTED_FILE_COUNT=$(files_by_pattern "$file" | wc -l)
      if [ $NESTED_FILE_COUNT -gt 0 ]; then
        echo "   skip dir without relevant files [$file]"
        TOTAL_DIRS=$((TOTAL_DIRS + 1))
        continue
      fi
      local ALL_BELOW=$(all_except_pattern "$file" | wc -l)
      local FILES_BELOW=$(files_by_pattern "$file" | wc -l)
      DIR_SIZES+=("$file:count:$ALL_BELOW")
      TOTAL_DIRS=$((TOTAL_DIRS + 1))
      echo "    all below: $ALL_BELOW"
      echo "  files below: $FILES_BELOW"
    elif [ -f "$file" ]; then # file
      local LINE_COUNT=$(wc -l < "$file" | awk '{print $1}')
      local USEFUL_COUNT=$(non_empty_lines "$file" | wc -l | awk '{print $1}')
      local EXCLUDED_STATUS=$(is_excluded "$file")
      if [[ $EXCLUDED_STATUS -eq 0 ]]; then
        TOTAL_FILES=$((TOTAL_FILES + 1))
        if [[ $LONGEST_LINES -eq 0 ]] || [[ $LINE_COUNT -gt $LONGEST_LINES ]]; then
          LONGEST_LINES="$LINE_COUNT"
          LONGEST_FILE="$file"
        else 
          echo "   minmax unchanged"
        fi
      else
        echo "   minmax skip"
      fi
      TOTAL_LINES=$((TOTAL_LINES + USEFUL_COUNT))
      FILE_SIZES+=("$file:lines:$USEFUL_COUNT")
      FILE_DEPTH=$(echo "$file" | tr -cd '/' | wc -c | awk '{print $1}')
      echo "        depth: $FILE_DEPTH"
      if [[ $FILE_DEPTH -gt $DEEPEST_NEST ]]; then
        DEEPEST_NEST=$FILE_DEPTH
        DEEPEST_FILE="$file"
      fi
      echo "   line count: $LINE_COUNT"
      echo " useful count: $USEFUL_COUNT"
    fi
  done
}

# MAIN
CSV_FILE="./metrics/project-size-over-time.csv"
HEADER=("date" 
  "timestamp" 
  "total_files" 
  "total_non_empty_lines"
  "greatest_file_depth" 
  "largest_folder_name" 
  "largest_folder_size" 
  "min_file_lines" 
  "min_file_name" 
  "mid_file_lines" #20-80% size median
  "mid_file_file"  #20-80% size median
  "max_file_lines"  
  "max_file_name"
)
CSV_HEADER=$(printf "%s," "${HEADER[@]}")
CSV_HEADER=${CSV_HEADER% ","}
if [ ! -f "$CSV_FILE" ]; then
  echo "$CSV_HEADER" > "$CSV_FILE"
fi
echo "---------------------------------------------"
echo "   CSV_FILE: $CSV_FILE"
echo " CSV_HEADER: $CSV_HEADER"
echo "---------------------------------------------"
echo "---------------------------------------------"

take_metrics

echo "---------------------------------------------"
echo "    WORKING_DIR: '$(pwd)'"
echo "     TOTAL_DIRS: $TOTAL_DIRS"
echo "    TOTAL_FILES: $TOTAL_FILES matching: [ "*.go", "*.html", "*.css", "*.js", "*.md" ]"
echo "    TOTAL_LINES: $TOTAL_LINES"
echo "   LONGEST_FILE: $LONGEST_FILE"
echo "  LONGEST_LINES: $LONGEST_LINES"
echo "   DEEPEST_FILE: $DEEPEST_FILE"
echo "   DEEPEST_NEST: $DEEPEST_NEST"
echo "---------------------------------------------"
echo "---------------------------------------------"

# TODO averages

DATE=$(date +%Y-%m-%d)
TIMESTAMP=$(date +%s)

#echo "$DATE,$TIMESTAMP,$TOTAL_FILES,$TOTAL_LINES,  
  # "greatest_file_depth" 
  # "largest_folder_name" 
  # "largest_folder_size" 
  # "min_file_lines" 
  # "min_file_name" 
  # "mid_file_lines" #20-80% size median
  # "mid_file_file"  #20-80% size median
  # "max_file_lines"  
  # "max_file_name"
# " >> "$CSV_FILE"

exit 1