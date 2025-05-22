#!/usr/bin/env bash

set -euo pipefail

usage() {
  echo "Usage: $0 <file-path>"
  exit 1
}

if [ "$#" -ne 1 ]; then
  usage
fi

FILE="$1"

if [ ! -f "$FILE" ]; then
  echo "Error: File not found: $FILE"
  exit 1
fi

OS="$(uname)"
if [[ "$OS" == "Darwin" ]]; then
  # macOS: must use -i for input file
  base64 -i "$FILE"
elif [[ "$OS" == "Linux" ]]; then
  # Linux: supports passing file as positional arg
  base64 "$FILE"
else
  echo "Unsupported platform: $OS"
  exit 1
fi
