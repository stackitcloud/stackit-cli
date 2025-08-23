#!/bin/bash

set -euo pipefail

# Create temporary GPG key file
TEMP_KEY_FILE=$(mktemp)
echo "$GPG_PRIVATE_KEY" > "$TEMP_KEY_FILE"

# Import key into GPG keyring
gpg --batch --import "$TEMP_KEY_FILE"

# Set environment variables for GoReleaser
export GPG_KEY_PATH="$TEMP_KEY_FILE"
export NFPM_LINUX_PACKAGES_RPM_PASSPHRASE="$GPG_PASSPHRASE"

goreleaser release --clean --config goreleaser.rpm.yaml --snapshot

# Clean up
rm -f "$TEMP_KEY_FILE"

echo "RPM build complete. Find RPMs in ./dist/" 