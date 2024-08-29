#!/bin/bash

# This script is used to publish new packages to the CLI APT repository
# Usage: ./publish-apt-packages.sh
set -eo pipefail

ROOT_DIR=$(git rev-parse --show-toplevel)

PACKAGES_BUCKET_URL="https://packages.stackit.cloud"
PUBLIC_KEY_FILE_PATH="keys/key.gpg"
APT_REPO_PATH="apt/cli"
APT_BUCKET_NAME="distribution"
CUSTOM_KEYRING_FILE="aptly-keyring.gpg"
DISTRIBUTION="stackit"
APTLY_CONFIG_FILE_PATH="./.aptly.conf"
GORELEASER_PACKAGES_FOLDER="dist/"

# We need to disable the key database daemon (keyboxd)
# This can be done by removing "use-keyboxd" from ~/.gnupg/common.conf (see https://github.com/gpg/gnupg/blob/master/README)
echo -n >~/.gnupg/common.conf

# Create a local mirror of the current state of the remote APT repository
printf ">>> Creating mirror \n"
curl ${PACKAGES_BUCKET_URL}/${PUBLIC_KEY_FILE_PATH} >public.asc
gpg --no-default-keyring --keyring=${CUSTOM_KEYRING_FILE} --import public.asc
aptly mirror create -config "${APTLY_CONFIG_FILE_PATH}" -keyring="${CUSTOM_KEYRING_FILE}" current "${PACKAGES_BUCKET_URL}/${APT_REPO_PATH}" ${DISTRIBUTION}

# Update the mirror to the latest state
printf "\n>>> Updating mirror \n"
aptly mirror update -keyring="${CUSTOM_KEYRING_FILE}" current

# Create a snapshot of the mirror
printf "\n>>> Creating snapshop from mirror \n"
aptly snapshot create current-snapshot from mirror current

# Create a new fresh local APT repo
printf "\n>>> Creating fresh local repo \n"
aptly repo create -distribution="${DISTRIBUTION}" new-repo

# Add new generated .deb packages to the new local repo
printf "\n>>> Adding new packages to local repo \n"
aptly repo add new-repo ${GORELEASER_PACKAGES_FOLDER}

# Create a snapshot of the local repo
printf "\n>>> Creating snapshot of local repo \n"
aptly snapshot create new-snapshot from repo new-repo

# Merge new-snapshot into current-snapshot creating a new snapshot updated-snapshot
printf "\n>>> Merging snapshots \n"
aptly snapshot pull -no-remove -architectures="amd64,i386,arm64" current-snapshot new-snapshot updated-snapshot ${DISTRIBUTION}

# Publish the new snapshot to the remote repo
printf "\n>>> Publishing updated snapshot \n"
aptly publish snapshot -keyring="${CUSTOM_KEYRING_FILE}" -gpg-key="${GPG_PRIVATE_KEY_FINGERPRINT}" -passphrase "${GPG_PASSPHRASE}" -config "${APTLY_CONFIG_FILE_PATH}" updated-snapshot "s3:${APT_BUCKET_NAME}:${APT_REPO_PATH}"
