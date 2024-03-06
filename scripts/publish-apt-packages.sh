#!/bin/bash

# This script is used to publish new packages to the CLI APT repository
# Usage: ./publish-apt-packages.sh
set -eo pipefail

ROOT_DIR=$(git rev-parse --show-toplevel)

OBJECT_STORAGE_ENDPOINT="https://object.storage.eu01.onstackit.cloud"
APT_BUCKET_NAME="stackit-cli-apt"
PUBLIC_KEY_BUCKET_NAME="stackit-public-key"
PUBLIC_KEY_FILE="key.gpg"
CUSTIM_KEYRING_PIPELINE_FOLDER="/Users/runner/.gnupg"
CUSTOM_KEYRING_FILE="custom-keyring.gpg"
DISTRIBUTION="stackit"
APTLY_CONFIG_FILE_PATH="./.aptly.conf"
GORELEASER_PACKAGES_FOLDER="dist/"

# Create a local mirror of the current state of the remote APT repository
printf ">>> Creating mirror \n"
curl ${OBJECT_STORAGE_ENDPOINT}/${PUBLIC_KEY_BUCKET_NAME}/${PUBLIC_KEY_FILE} >public.asc
gpg --no-default-keyring --keyring ${CUSTIM_KEYRING_PIPELINE_FOLDER}/${CUSTOM_KEYRING_FILE} --import public.asc
aptly mirror create -keyring="${CUSTOM_KEYRING_FILE}" current "${OBJECT_STORAGE_ENDPOINT}/${APT_BUCKET_NAME}" ${DISTRIBUTION} # Folder (CUSTIM_KEYRING_PIPELINE_FOLDER) is appended automatic in the aptly command

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
aptly publish switch -gpg-key="${GPG_PRIVATE_KEY_ID}" -passphrase "${GPG_PASSPHRASE}" -config "${APTLY_CONFIG_FILE_PATH}" ${DISTRIBUTION} "s3:${APT_BUCKET_NAME}:" updated-snapshot
