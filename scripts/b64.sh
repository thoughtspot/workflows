#!/usr/bin/env bash

set -euo pipefail

usage() {
  echo "Usage: $0 [encode|decode] <file-path>"
  exit 1
}

if [ "$#" -ne 2 ]; then
  usage
fi

MODE="$1"
FILE="$2"

if [ ! -f "$FILE" ]; then
  echo "Error: File not found: $FILE"
  exit 1
fi

# Extract filename without path and extension
FILENAME_BASE="$(basename "$FILE")"
FILENAME_NOEXT="${FILENAME_BASE%.*}"
ENCODED_FILE="${FILENAME_NOEXT}-b64-encoded.txt"

OS="$(uname)"
if [[ "$OS" == "Darwin" ]]; then
  if [[ "$MODE" == "encode" ]]; then
    base64 -i "$FILE" > "$ENCODED_FILE"
    echo "Encoded file created: $ENCODED_FILE"
  elif [[ "$MODE" == "decode" ]]; then
    base64 -d -i "$FILE"
  else
    usage
  fi
elif [[ "$OS" == "Linux" ]]; then
  if [[ "$MODE" == "encode" ]]; then
    base64 "$FILE" > "$ENCODED_FILE"
    echo "Encoded file created: $ENCODED_FILE"
  elif [[ "$MODE" == "decode" ]]; then
    base64 -d "$FILE"
  else
    usage
  fi
else
  echo "Unsupported platform: $OS"
  exit 1
fi
