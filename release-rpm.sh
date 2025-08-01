#!/bin/bash

set -euo pipefail

: "${GPG_PRIVATE_KEY:?GPG_PRIVATE_KEY must be set}"
: "${GPG_PASSPHRASE:?GPG_PASSPHRASE must be set}"

export GPG_PRIVATE_KEY
export GPG_PASSPHRASE

gpg --batch --import <<< "$GPG_PRIVATE_KEY"

goreleaser release --clean --config goreleaser.rpm.yaml --skip-publish --skip-validate

echo "RPM build complete. Find RPMs in ./dist/" 