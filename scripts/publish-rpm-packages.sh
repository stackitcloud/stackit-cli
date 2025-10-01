#!/usr/bin/env bash

# This script is used to publish new RPM packages to the CLI RPM repository
# Usage: ./publish-rpm-packages.sh
set -eo pipefail

PACKAGES_BUCKET_URL="https://distribution-test.object.storage.eu01.onstackit.cloud"
PUBLIC_KEY_FILE_PATH="keys/key.gpg"
RPM_REPO_PATH="rpm/cli"
RPM_BUCKET_NAME="distribution-test"
CUSTOM_KEYRING_FILE="rpm-keyring.gpg"
DISTRIBUTION="stackit"
GORELEASER_PACKAGES_FOLDER="dist/"

# We need to disable the key database daemon (keyboxd)
# This can be done by removing "use-keyboxd" from ~/.gnupg/common.conf (see https://github.com/gpg/gnupg/blob/master/README)
echo -n >~/.gnupg/common.conf

# Create RPM repository directory structure
printf ">>> Creating RPM repository structure \n"
mkdir -p rpm-repo/x86_64
mkdir -p rpm-repo/i386
mkdir -p rpm-repo/aarch64

# Copy RPM packages and signatures to appropriate architecture directories
printf "\n>>> Copying RPM packages and signatures to architecture directories \n"

# Copy x86_64 packages (amd64)
for file in ${GORELEASER_PACKAGES_FOLDER}*_amd64.rpm*; do
    if [ -f "$file" ]; then
        cp "$file" rpm-repo/x86_64/
        printf "Copied $(basename "$file") to x86_64/\n"
    fi
done

# Copy i386 packages
for file in ${GORELEASER_PACKAGES_FOLDER}*_386.rpm*; do
    if [ -f "$file" ]; then
        cp "$file" rpm-repo/i386/
        printf "Copied $(basename "$file") to i386/\n"
    fi
done

# Copy aarch64 packages (arm64)
for file in ${GORELEASER_PACKAGES_FOLDER}*_arm64.rpm*; do
    if [ -f "$file" ]; then
        cp "$file" rpm-repo/aarch64/
        printf "Copied $(basename "$file") to aarch64/\n"
    fi
done

# Download existing repository metadata if it exists
printf "\n>>> Downloading existing repository metadata \n"
aws s3 sync s3://${RPM_BUCKET_NAME}/${RPM_REPO_PATH}/ rpm-repo/ --delete || echo "No existing repository found, creating new one"

# Create repository metadata for each architecture
printf "\n>>> Creating repository metadata \n"
for arch in x86_64 i386 aarch64; do
    if [ -d "rpm-repo/${arch}" ] && [ "$(ls -A rpm-repo/${arch})" ]; then
        printf "Creating metadata for ${arch}...\n"
        
        # List what we're working with
        printf "Files in ${arch}: $(ls rpm-repo/${arch}/ | tr '\n' ' ')\n"
        
        # Create repository metadata
        createrepo_c --update rpm-repo/${arch}
        
        # Sign the repository metadata
        printf "Signing repository metadata for ${arch}...\n"
        gpg --detach-sign --armor --local-user "${GPG_PRIVATE_KEY_FINGERPRINT}" --passphrase "${GPG_PASSPHRASE}" rpm-repo/${arch}/repodata/repomd.xml
        
        # Verify the signature was created
        if [ -f "rpm-repo/${arch}/repodata/repomd.xml.asc" ]; then
            printf "Repository metadata signed successfully for ${arch}\n"
        else
            printf "WARNING: Repository metadata signature not created for ${arch}\n"
        fi
    else
        printf "No packages found for ${arch}, skipping...\n"
    fi
done

# Upload the updated repository to S3
printf "\n>>> Uploading repository to S3 \n"
aws s3 sync rpm-repo/ s3://${RPM_BUCKET_NAME}/${RPM_REPO_PATH}/ --delete

# Upload the public key
printf "\n>>> Uploading public key \n"
gpg --armor --export "${GPG_PRIVATE_KEY_FINGERPRINT}" > public-key.asc
aws s3 cp public-key.asc s3://${RPM_BUCKET_NAME}/${PUBLIC_KEY_FILE_PATH}

printf "\n>>> RPM repository published successfully! \n"
printf "Repository URL: ${PACKAGES_BUCKET_URL}/${RPM_REPO_PATH}/ \n"
printf "Public key URL: ${PACKAGES_BUCKET_URL}/${PUBLIC_KEY_FILE_PATH} \n"
