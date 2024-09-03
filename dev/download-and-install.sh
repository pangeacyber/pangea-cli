#!/bin/bash

# Define variables
FILE_NAME=$(echo "pangea-$(uname -s)-$(uname -m).tar.gz" | tr '[:upper:]' '[:lower:]')  # lowercase name
REPO=pangeacyber/pangea-cli

# Download the file
download_file() {
  local filename="$1"
  DOWNLOAD_URL=$(echo "https://github.com/$REPO/releases/latest/download/$filename")
  echo "Downloading from $DOWNLOAD_URL..."
  HTTP_STATUS=$(curl -L -w "%{http_code}" $DOWNLOAD_URL -o $filename -s)

  # Check if the HTTP status code indicates success (2xx)
  if [[ "${HTTP_STATUS}" != 2* ]]; then
    echo "Download failed with HTTP status code ${HTTP_STATUS}."
    exit 1
  fi
  echo "Download completed."
}

download_file $FILE_NAME

# Make dir and uncompress
mkdir -p installer && tar -xzvf $FILE_NAME -C installer
cd installer

echo "Installing..."
./install.sh

echo "Install completed."
