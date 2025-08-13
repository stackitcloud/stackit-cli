#!/bin/bash

# This script is used to publish new packages to the CLI RPM repository
# Usage: ./publish-rpm-packages.sh
set -eo pipefail

ROOT_DIR=$(git rev-parse --show-toplevel)

PACKAGES_BUCKET_URL="https://packages.stackit.cloud"
RPM_REPO_PATH="rpm/cli"
RPM_BUCKET_NAME="distribution"
CUSTOM_KEYRING_FILE="rpm-keyring.gpg"
GORELEASER_PACKAGES_FOLDER="dist/"
TEMP_DIR=$(mktemp -d)

# We need to disable the key database daemon (keyboxd)
# This can be done by removing "use-keyboxd" from ~/.gnupg/common.conf (see https://github.com/gpg/gnupg/blob/master/README)
echo -n >~/.gnupg/common.conf

# Create a local mirror of the current state of the remote RPM repository
printf ">>> Creating mirror \n"
curl ${PACKAGES_BUCKET_URL}/${RPM_REPO_PATH}/repodata/repomd.xml >${TEMP_DIR}/repomd.xml || echo "No existing repository found, creating new one"

# Create RPM repository structure
mkdir -p ${TEMP_DIR}/rpm-repo/RPMS

# Copy existing RPMs from remote repository (if any)
printf "\n>>> Downloading existing RPMs \n"
aws s3 sync s3://${RPM_BUCKET_NAME}/${RPM_REPO_PATH}/RPMS/ ${TEMP_DIR}/rpm-repo/RPMS/ --endpoint-url https://object.storage.eu01.onstackit.cloud || echo "No existing RPMs found"

# Copy new generated .rpm packages to the local repo
# Note: GoReleaser already signs these RPM packages with embedded signatures
printf "\n>>> Adding new packages to local repo \n"
cp ${GORELEASER_PACKAGES_FOLDER}/*.rpm ${TEMP_DIR}/rpm-repo/RPMS/

# Create RPM repository metadata using createrepo_c
printf "\n>>> Creating RPM repository metadata \n"
createrepo_c ${TEMP_DIR}/rpm-repo

# Sign the repository metadata using the same GPG key as APT
if [ -n "$GPG_PRIVATE_KEY_FINGERPRINT" ] && [ -n "$GPG_PASSPHRASE" ]; then
    printf "\n>>> Signing repository metadata \n"
    gpg --batch --yes --pinentry-mode loopback --local-user="${GPG_PRIVATE_KEY_FINGERPRINT}" --passphrase="${GPG_PASSPHRASE}" --detach-sign --armor ${TEMP_DIR}/rpm-repo/repodata/repomd.xml
else
    echo ">>> Skipping repository metadata signing (GPG environment variables not set)"
fi

# Upload to S3
printf "\n>>> Uploading to S3 \n"
aws s3 sync ${TEMP_DIR}/rpm-repo/ s3://${RPM_BUCKET_NAME}/${RPM_REPO_PATH}/ --endpoint-url https://object.storage.eu01.onstackit.cloud

# Clean up
rm -rf ${TEMP_DIR}

printf "\n>>> RPM repository published successfully to ${PACKAGES_BUCKET_URL}/${RPM_REPO_PATH} \n" 