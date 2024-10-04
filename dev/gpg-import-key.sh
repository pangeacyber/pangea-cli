# This command is needed to initialize gpg
import_gpg_key(){
    echo "$GPG_PRIVATE_KEY" | gpg --import 2>&1 | awk '/gpg: key/{gsub(/:/, "", $3); print $3; exit}'
}

if [ -z "$GPG_PRIVATE_KEY" ]; then
    echo "GPG_PRIVATE_KEY env var is empty."
    exit 1
fi

if [ -z "$GPG_KEY_ID" ]; then
    gpg --list-secret-keys
    echo "import key..."
    export GPG_KEY_ID=$(import_gpg_key)
    echo "imported."
else
    echo "GPG_KEY_ID is already set."
fi

echo "GPG_KEY_ID:" $GPG_KEY_ID
