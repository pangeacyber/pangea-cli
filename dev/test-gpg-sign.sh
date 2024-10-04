#!/bin/bash

# This will be set on a Gitlab env var
# Warning! This is just a develop key, should not be used in production releases
export GPG_PRIVATE_KEY=$(cat dev/gpg.pem)

source ./dev/gpg-import-key.sh
make build-all
./dev/gpg-sign-binaries.sh
