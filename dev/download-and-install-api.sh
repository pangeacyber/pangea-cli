#!/bin/bash

[[ "$(command -v cosign)" ]] || { echo "cosign is not installed" 1>&2 ; exit 1; }

# Define variables
FILE_NAME=$(echo "pangea-$(uname -s)-$(uname -m).tar.gz" | tr '[:upper:]' '[:lower:]')  # lowercase name
SIGNATURE_FILENAME=$FILE_NAME.sig
PUBLIC_KEY_FILENAME=cosign.pub
REPO=pangeacyber/pangea-cli

# First argument is the release tag, if not set will be default to latest
RELEASE_TAG=$1

# Check github token is set
if [ -z "$GITHUB_TOKEN" ]; then
  echo "Error: Need to set "GITHUB_TOKEN" to have access to this release"
  exit 1
fi

# If no release tag was passed as argument, fetch the latest
if [ -z "$RELEASE_TAG" ]; then
  LATEST_RELEASE=$(curl -H "Authorization: token ${GITHUB_TOKEN}" -s https://api.github.com/repos/${REPO}/releases/latest)
  RELEASE_TAG=$(echo $LATEST_RELEASE | jq -r .tag_name)
fi
echo "Target release: $RELEASE_TAG"

# Build download url
DOWNLOAD_URL="https://api.github.com/repos/${REPO}/releases/tags/$RELEASE_TAG"
echo "Download URL: $DOWNLOAD_URL"

download_file() {
  local filename="$1"
  # Get the asset ID of the file
  ASSET_ID=$(curl -s -H "Authorization: token ${GITHUB_TOKEN}" $DOWNLOAD_URL | jq -r ".assets[] | select(.name == \"$filename\") | .id")
  if [ -z "$ASSET_ID" ]; then
    echo "Error: Could not find asset with name $filename in release $RELEASE_TAG"
    exit 1
  fi

  # Download the file
  echo "Downloading $filename from GitHub..."
  curl -L -H "Authorization: token ${GITHUB_TOKEN}" -H "Accept: application/octet-stream" "https://api.github.com/repos/${REPO}/releases/assets/$ASSET_ID" -o $filename
  echo "Download $filename completed."
}

# Download files
download_file $FILE_NAME

if command -v cosign &> /dev/null; then
  echo "Verify signature..."
  download_file $PUBLIC_KEY_FILENAME
  download_file $SIGNATURE_FILENAME
  cosign verify-blob --key $PUBLIC_KEY_FILENAME --signature $SIGNATURE_FILENAME $FILE_NAME

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
