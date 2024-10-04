#!/bin/bash

# Define variables
FILE_NAME=$(echo "pangea-$(uname -s)-$(uname -m).tar.gz" | tr '[:upper:]' '[:lower:]')  # lowercase name
SIGNATURE_FILENAME=$FILE_NAME.sig
PUBLIC_KEY_FILENAME=cosign.pub
REPO=pangeacyber/pangea-cli-internal

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

if command -v cosign &> /dev/null; then
  echo "Verify signature..."
  download_file $PUBLIC_KEY_FILENAME
  download_file $SIGNATURE_FILENAME
  cosign verify-blob --key $PUBLIC_KEY_FILENAME -signature $SIGNATURE_FILENAME $FILE_NAME

  # Check the exit code of the cosign command
  if [ $? -ne 0 ]; then
    echo "Error: cosign signature verification failed for $FILE_NAME."
    exit 1
  fi
else
  echo "cosign is not installed. Signature verification skipped. Please install cosign to verify package signature."
fi

# Make dir and uncompress
mkdir -p installer && tar -xzvf $FILE_NAME -C installer
cd installer

echo "Installing..."
./install.sh

echo "Install completed."
