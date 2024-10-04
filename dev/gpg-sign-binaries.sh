#!/bin/bash

# Set the required variables
SIGNATURES_FOLDER="./signatures"
BIN_DIR="./bin"

# Function to check if the GPG_KEY_ID environment variable is set
check_env() {
    if [ -z "$GPG_KEY_ID" ]; then
        echo "Error: GPG_KEY_ID environment variable is not set"
        exit 1
    fi
}

# Main script
main() {
    check_env

    # Create the signatures folder if it doesn't exist
    mkdir -p "${SIGNATURES_FOLDER}"

    # Remove old signatures if it already
    rm ./${SIGNATURES_FOLDER}/*

    # Iterate over files in the BIN_DIR and sign them
    for file in "${BIN_DIR}"/*; do
        if [ -f "$file" ]; then
            echo "Signing $file..."
            gpg -a -u "${GPG_KEY_ID}" --digest-algo SHA256 --output "${SIGNATURES_FOLDER}/$(basename "$file").sig" --detach-sig "$file"
            echo "Signature:"
            cat "${SIGNATURES_FOLDER}/$(basename "$file").sig"
        fi
    done
}

# Run the main script
main
