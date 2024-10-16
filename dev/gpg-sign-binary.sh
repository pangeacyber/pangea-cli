#!/bin/bash

# Set the required variables
SIGNATURES_FOLDER="./signatures"
filename=$1

# Function to check if the GPG_KEY_ID environment variable is set
check_env() {
    echo "Checking GPG_KEY_ID env var..."
    if [ -z "$GPG_KEY_ID" ]; then
        echo "Error: GPG_KEY_ID environment variable is not set"
        exit 1
    fi
}

# Main script
main() {

    if [ -z "$filename" ]; then
        echo "Error: 'filename' variable is empty"
        exit 1
    fi

    check_env

    echo "Creating ${SIGNATURES_FOLDER} folder..."
    # Create the signatures folder if it doesn't exist
    mkdir -p "${SIGNATURES_FOLDER}"

    echo "Signing $filename file..."
    gpg -a -u "${GPG_KEY_ID}" --digest-algo SHA256 --output "${SIGNATURES_FOLDER}/$(basename "${filename}").sig" --detach-sig "${filename}"
    echo "Signature:"
    cat "${SIGNATURES_FOLDER}/$(basename "${filename}").sig"
}

# Run the main script
main
