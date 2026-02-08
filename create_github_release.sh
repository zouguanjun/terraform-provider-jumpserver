#!/bin/bash

# Terraform Provider JumpServer v1.0.0 Release Script
# Owner: zouguanjun

set -e

# Configuration
GITHUB_TOKEN="ghp_KNQyUQ0T7eX4Q6N4Z8yJcR0pWqMkLsX1vF2h"
REPO="zouguanjun/terraform-provider-jumpserver"
VERSION="v1.0.0"
TAG_NAME="v1.0.0"
RELEASE_TITLE="v1.0.0 - Initial Stable Release"

echo "========================================="
echo "Pushing tag to GitHub..."
echo "========================================="

# 更新远程 URL
git remote set-url origin https://${GITHUB_TOKEN}@github.com/${REPO}.git

# 推送标签
echo "Pushing tag ${TAG_NAME}..."
git push origin ${TAG_NAME}

echo ""
echo "========================================="
echo "Creating GitHub Release..."
echo "========================================="

# 使用 curl 创建 Release
echo "Creating release..."
curl -s -X POST \
  -H "Authorization: token ${GITHUB_TOKEN}" \
  -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/repos/${REPO}/releases \
  -d "{
    \"tag_name\": \"${TAG_NAME}\",
    \"name\": \"${RELEASE_TITLE}\",
    \"body\": \"$(cat RELEASE_NOTES.md)\",
    \"draft\": false,
    \"prerelease\": false
  }" > release_response.json

# 获取 Release ID
RELEASE_ID=$(cat release_response.json | grep -oP '"id":\s*\K[0-9]+' || cat release_response.json | jq -r '.id' 2>/dev/null)

echo "Release created with ID: ${RELEASE_ID}"

# 上传二进制文件
echo ""
echo "========================================="
echo "Uploading binary to Release..."
echo "========================================="

BINARY_FILE="release/clean/terraform-provider-jumpserver_1.0.0_linux_amd64.zip"

if [ -f "${BINARY_FILE}" ]; then
    echo "Uploading ${BINARY_FILE}..."
    
    curl -s -X POST \
      -H "Authorization: token ${GITHUB_TOKEN}" \
      -H "Content-Type: application/zip" \
      --data-binary @"${BINARY_FILE}" \
      https://uploads.github.com/repos/${REPO}/releases/${RELEASE_ID}/assets?name=$(basename ${BINARY_FILE})
    
    echo "Binary uploaded successfully!"
else
    echo "Error: Binary file not found at ${BINARY_FILE}"
    exit 1
fi

echo ""
echo "========================================="
echo "Release completed successfully!"
echo "========================================="
echo "Release URL: https://github.com/${REPO}/releases/tag/${TAG_NAME}"
