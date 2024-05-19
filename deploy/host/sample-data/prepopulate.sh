#!/bin/bash
set -eu

SAMPLE_DATA=${1:-}
REAL_DATA=${2:-}

if [ -z "${SAMPLE_DATA}" ] || [ -z "${REAL_DATA}" ]; then
    echo "Usage: $0 <sample-data> <real-data>"
    exit 1
fi

if [ -n "$(ls -A ${REAL_DATA})" ]; then
    echo "Data already exists in ${REAL_DATA}"
    exit 0
fi

if ! [ -d "${SAMPLE_DATA}" ]; then
    echo "No sample data to copy"
    exit 0
fi

cp -r "${SAMPLE_DATA}/." "${REAL_DATA}/"
exit 0