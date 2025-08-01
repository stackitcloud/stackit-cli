#!/bin/bash

set -euo pipefail

RPM_OUTPUT_DIR="dist"
TEMP_DIR=$(mktemp -d)
GPG_PRIVATE_KEY_FINGERPRINT="${GPG_PRIVATE_KEY_FINGERPRINT:?Set GPG_PRIVATE_KEY_FINGERPRINT}"
GPG_PASSPHRASE="${GPG_PASSPHRASE:?Set GPG_PASSPHRASE}"
AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID:?Set AWS_ACCESS_KEY_ID}"
AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY:?Set AWS_SECRET_ACCESS_KEY}"

# Test environment S3 bucket
S3_BUCKET="distribution-test"
S3_ENDPOINT="https://object.storage.eu01.onstackit.cloud"
RPM_REPO_PATH="rpm/cli"

echo ">>> Preparing RPM repository structure..."
mkdir -p "$TEMP_DIR/rpm-repo/RPMS"

echo ">>> Copying built RPMs..."
cp "$RPM_OUTPUT_DIR"/*.rpm "$TEMP_DIR/rpm-repo/RPMS/"

echo ">>> Creating RPM repository metadata..."
createrepo_c "$TEMP_DIR/rpm-repo"

echo ">>> Signing repository metadata..."
gpg --batch --yes --pinentry-mode loopback \
    --local-user="$GPG_PRIVATE_KEY_FINGERPRINT" \
    --passphrase="$GPG_PASSPHRASE" \
    --detach-sign --armor "$TEMP_DIR/rpm-repo/repodata/repomd.xml"

echo ">>> Uploading to test bucket..."
aws s3 sync "$TEMP_DIR/rpm-repo/" "s3://$S3_BUCKET/$RPM_REPO_PATH/" \
    --endpoint-url "$S3_ENDPOINT" \
    --delete

rm -rf "$TEMP_DIR"
echo ">>> RPM repo published to test bucket: $S3_BUCKET/$RPM_REPO_PATH" 