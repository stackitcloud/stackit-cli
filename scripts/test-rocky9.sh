#!/bin/bash

# Test script for Rocky Linux 9 RPM repository
# Uses test bucket for RPMs, production bucket for GPG key

set -e

echo "=========================================="
echo "STACKIT CLI Rocky Linux 9 RPM Test"
echo "=========================================="

# Configuration
CONTAINER_NAME="stackit-rocky9-test"
IMAGE="rockylinux:9"

# Test environment S3 bucket (for RPMs)
TEST_S3_BUCKET="distribution-test"
TEST_S3_ENDPOINT="object.storage.eu01.onstackit.cloud"
TEST_RPM_REPO_PATH="rpm/cli"

# Production S3 bucket (for GPG key)
PROD_S3_BUCKET="distribution"
PROD_S3_ENDPOINT="object.storage.eu01.onstackit.cloud"
PROD_GPG_KEY_PATH="keys/key.gpg"

echo "Step 1: Starting Rocky Linux 9 container..."
docker run -d --name $CONTAINER_NAME $IMAGE tail -f /dev/null

echo "Step 2: Installing dependencies..."
docker exec $CONTAINER_NAME bash -c "
    dnf update -y
    dnf install -y --allowerasing curl wget gpg
"

echo "Step 3: Downloading GPG key from production bucket..."
docker exec $CONTAINER_NAME bash -c "
    curl -o /tmp/stackit-gpg-signer.asc 'https://$PROD_S3_BUCKET.$PROD_S3_ENDPOINT/$PROD_GPG_KEY_PATH'
    gpg --import /tmp/stackit-gpg-signer.asc
    echo '✅ GPG key imported'
"

echo "Step 4: Creating repository configuration..."
docker exec $CONTAINER_NAME bash -c "
    cat > /etc/yum.repos.d/stackit-cli.repo << EOF
[stackit-cli]
name=STACKIT CLI Repository
baseurl=https://$TEST_S3_BUCKET.$TEST_S3_ENDPOINT/$TEST_RPM_REPO_PATH
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://$PROD_S3_BUCKET.$PROD_S3_ENDPOINT/$PROD_GPG_KEY_PATH
EOF
    cat /etc/yum.repos.d/stackit-cli.repo
    echo '✅ Repository configuration created'
"

echo "Step 5: Updating package cache..."
docker exec $CONTAINER_NAME bash -c "
    dnf clean all
    dnf repolist
    echo '✅ Package cache updated'
"

echo "Step 6: Installing STACKIT CLI..."
docker exec $CONTAINER_NAME bash -c "
    dnf install -y stackit
    echo '✅ STACKIT CLI installed'
"

echo "Step 7: Verifying installation..."
docker exec $CONTAINER_NAME bash -c "
    if command -v stackit >/dev/null 2>&1; then
        echo '✅ stackit command found: \$(which stackit)'
        echo '✅ Version: \$(stackit version)'
    else
        echo '❌ stackit command not found'
        exit 1
    fi
"

echo "Step 8: Testing basic functionality..."
docker exec $CONTAINER_NAME bash -c "
    echo '=== STACKIT CLI HELP OUTPUT ==='
    stackit --help
    echo '=== END HELP OUTPUT ==='
    echo '✅ Basic functionality test passed'
"

echo "Step 9: Testing package update..."
docker exec $CONTAINER_NAME bash -c "
    dnf check-update stackit || echo 'No updates available (expected for test)'
    echo '✅ Update check completed'
"

echo "Step 10: Uninstalling STACKIT CLI..."
docker exec $CONTAINER_NAME bash -c "
    dnf remove -y stackit
    echo '✅ STACKIT CLI uninstalled'
"

echo "Step 11: Verifying uninstallation..."
docker exec $CONTAINER_NAME bash -c "
    if ! command -v stackit >/dev/null 2>&1; then
        echo '✅ stackit command no longer found'
    else
        echo '❌ stackit command still found: \$(which stackit)'
        exit 1
    fi
"

echo "Step 12: Cleaning up container..."
docker stop $CONTAINER_NAME
docker rm $CONTAINER_NAME

echo "=========================================="
echo "✅ Rocky Linux 9 RPM test completed successfully!"
echo "==========================================" 