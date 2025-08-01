#!/bin/bash

set -euo pipefail

: "${GPG_PRIVATE_KEY:?GPG_PRIVATE_KEY must be set}"
: "${GPG_PASSPHRASE:?GPG_PASSPHRASE must be set}"

# Create temporary GPG key file
TEMP_KEY_FILE=$(mktemp)
echo "$GPG_PRIVATE_KEY" > "$TEMP_KEY_FILE"

# Import key into GPG keyring
gpg --batch --import "$TEMP_KEY_FILE"

# Set environment variable for GoReleaser to use the key file path
export GPG_PRIVATE_KEY="$TEMP_KEY_FILE"
export GPG_PASSPHRASE

goreleaser release --clean --config goreleaser.rpm.yaml --snapshot

# Clean up
rm -f "$TEMP_KEY_FILE"

echo "RPM build complete. Find RPMs in ./dist/" 