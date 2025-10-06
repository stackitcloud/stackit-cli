#!/usr/bin/env bash

# This script is used to publish new packages to the CLI APT repository
# Usage: ./publish-apt-packages.sh
set -eo pipefail

PACKAGES_BUCKET_URL="https://distribution-test.object.storage.eu01.onstackit.cloud"
PUBLIC_KEY_FILE_PATH="keys/key.gpg"
APT_REPO_PATH="apt/cli"
APT_BUCKET_NAME="distribution-test"
CUSTOM_KEYRING_FILE="aptly-keyring.gpg"
DISTRIBUTION="stackit"
APTLY_CONFIG_FILE_PATH="./.aptly.conf"
GORELEASER_PACKAGES_FOLDER="dist/"

# We need to disable the key database daemon (keyboxd)
# This can be done by removing "use-keyboxd" from ~/.gnupg/common.conf (see https://github.com/gpg/gnupg/blob/master/README)
echo -n >~/.gnupg/common.conf

BOOTSTRAP_ONLY=0

# Try to create a local mirror of the current remote APT repository
printf ">>> Creating mirror (if remote exists)\n"
curl -fsSL ${PACKAGES_BUCKET_URL}/${PUBLIC_KEY_FILE_PATH} -o public.asc
gpg --no-default-keyring --keyring=${CUSTOM_KEYRING_FILE} --import public.asc
if aptly mirror create -config "${APTLY_CONFIG_FILE_PATH}" -keyring="${CUSTOM_KEYRING_FILE}" current "${PACKAGES_BUCKET_URL}/${APT_REPO_PATH}" ${DISTRIBUTION}; then
  printf "\n>>> Updating mirror \n"
  aptly mirror update -keyring="${CUSTOM_KEYRING_FILE}" -max-tries=5 current
  printf "\n>>> Creating snapshot from mirror \n"
  aptly snapshot create current-snapshot from mirror current
else
  printf "\n>>> No existing remote repository found (bootstrap)\n"
  BOOTSTRAP_ONLY=1
fi

# Create a new fresh local APT repo
printf "\n>>> Creating fresh local repo \n"
aptly repo create -distribution="${DISTRIBUTION}" new-repo

# Add new generated .deb packages to the new local repo
printf "\n>>> Adding new packages to local repo \n"
aptly repo add new-repo ${GORELEASER_PACKAGES_FOLDER}

# Create a snapshot of the local repo
printf "\n>>> Creating snapshot of local repo \n"
aptly snapshot create new-snapshot from repo new-repo

UPDATED_SNAPSHOT="new-snapshot"
if [ "$BOOTSTRAP_ONLY" -eq 0 ]; then
  # Merge new-snapshot into current-snapshot creating a new snapshot updated-snapshot
  printf "\n>>> Merging snapshots \n"
  aptly snapshot pull -no-remove -architectures="amd64,i386,arm64" current-snapshot new-snapshot updated-snapshot ${DISTRIBUTION}
  UPDATED_SNAPSHOT="updated-snapshot"
else
  printf "\n>>> Bootstrap mode: publishing new packages as initial snapshot \n"
fi

# Publish the new snapshot to the remote repo
printf "\n>>> Publishing updated snapshot \n"
aptly publish snapshot -keyring="${CUSTOM_KEYRING_FILE}" -gpg-key="${GPG_PRIVATE_KEY_FINGERPRINT}" -passphrase "${GPG_PASSPHRASE}" -config "${APTLY_CONFIG_FILE_PATH}" "$UPDATED_SNAPSHOT" "s3:${APT_BUCKET_NAME}:${APT_REPO_PATH}"
