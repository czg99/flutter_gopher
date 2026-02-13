#!/bin/bash

cd $(dirname $0)/../

GO_SRC="gosrc"
TIMESTAMP_FILE=".last_timestamp"

# Get the latest timestamp of the source code
get_gosrc_timestamp() {
	NEWEST_FILE=$(find ${GO_SRC} -type f -name "*.go" -printf "%T@ %p\n" 2>/dev/null | sort -nr | head -1)
	if [ -z "${NEWEST_FILE}" ]; then
		echo 0
		return
	fi

	NEWEST_TIMESTAMP=$(echo ${NEWEST_FILE} | cut -d' ' -f1 | cut -d'.' -f1)
	echo ${NEWEST_TIMESTAMP}
}

# Save the latest timestamp of the source code
save_gosrc_timestamp() {
	TIMESTAMP=$(get_gosrc_timestamp)
	if [ $TIMESTAMP -eq 0 ]; then
		return
	fi

	echo ${TIMESTAMP} >"${TIMESTAMP_FILE}"
}

# Check if source code has been updated
check_gosrc_changes() {
	TIMESTAMP=$(get_gosrc_timestamp)
	if [ $TIMESTAMP -eq 0 ]; then
		echo "No"
		return
	fi

	if [ ! -f "${TIMESTAMP_FILE}" ]; then
		echo "Yes"
		return
	fi

	LAST_TIMESTAMP=$(cat "${TIMESTAMP_FILE}")
	if [ "${TIMESTAMP}" -gt "${LAST_TIMESTAMP}" ]; then
		echo "Yes"
	else
		echo "No"
	fi
}

if [ "$1" == "get" ]; then
	get_gosrc_timestamp
elif [ "$1" == "save" ]; then
	save_gosrc_timestamp
elif [ "$1" == "check" ]; then
	check_gosrc_changes
else
	get_gosrc_timestamp
fi