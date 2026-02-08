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
echo "Creating GitHub Release via API..."
echo "========================================="

# 检查标签是否存在
echo "Checking if tag exists..."
TAG_EXISTS=$(curl -s -H "Authorization: token ${GITHUB_TOKEN}" \
  https://api.github.com/repos/${REPO}/git/ref/tags/${TAG_NAME} | grep -o '"ref"' || echo "")

if [ -z "$TAG_EXISTS" ]; then
    echo "Tag does not exist, creating..."

    # 创建 tag reference
    COMMIT_SHA=$(curl -s -H "Authorization: token ${GITHUB_TOKEN}" \
      https://api.github.com/repos/${REPO}/git/refs/heads/main | grep -oP '"sha":\s*"\K[^"]+')

    echo "Latest commit SHA: ${COMMIT_SHA}"

    # 创建 tag object
    TAG_RESPONSE=$(curl -s -X POST \
      -H "Authorization: token ${GITHUB_TOKEN}" \
      -H "Accept: application/vnd.github.v3+json" \
      https://api.github.com/repos/${REPO}/git/tags \
      -d "{
        \"tag\": \"${TAG_NAME}\",
        \"message\": \"Release v1.0.0 - Initial stable release\",
        \"object\": \"${COMMIT_SHA}\",
        \"type\": \"commit\",
        \"tagger\": {
          \"name\": \"zouguanjun\",
          \"email\": \"2023760206@qq.com\",
          \"date\": \"$(date -u +"%Y-%m-%dT%H:%M:%SZ")\"
        }
      }")

    echo "Tag created: ${TAG_RESPONSE}"

    # 创建 tag reference
    curl -s -X POST \
      -H "Authorization: token ${GITHUB_TOKEN}" \
      -H "Accept: application/vnd.github.v3+json" \
      https://api.github.com/repos/${REPO}/git/refs \
      -d "{
        \"ref\": \"refs/tags/${TAG_NAME}\",
        \"sha\": \"$(echo ${TAG_RESPONSE} | grep -oP '"sha":\s*"\K[^"]+')\"
      }"

    echo "Tag reference created."
else
    echo "Tag ${TAG_NAME} already exists."
fi

echo ""
echo "========================================="
echo "Creating Release..."
echo "========================================="

# 读取 Release Notes
if [ -f "RELEASE_NOTES.md" ]; then
    RELEASE_BODY=$(cat RELEASE_NOTES.md | sed 's/"/\\"/g' | tr '\n' '\\n' | sed 's/\\/\\\\/g')
else
    RELEASE_BODY="Initial stable release of JumpServer Terraform Provider v1.0.0"
fi

# 创建 Release
echo "Creating release..."
RELEASE_RESPONSE=$(curl -s -X POST \
  -H "Authorization: token ${GITHUB_TOKEN}" \
  -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/repos/${REPO}/releases \
  -d "{
    \"tag_name\": \"${TAG_NAME}\",
    \"name\": \"${RELEASE_TITLE}\",
    \"body\": \"${RELEASE_BODY}\",
    \"draft\": false,
    \"prerelease\": false
  }")

echo "Release response: ${RELEASE_RESPONSE}"

# 获取 Release ID
RELEASE_ID=$(echo "${RELEASE_RESPONSE}" | grep -oP '"id":\s*\K[0-9]+' || echo "")

if [ -z "$RELEASE_ID" ]; then
    echo "Error: Could not create release. Please check if the token has repo scope permissions."
    exit 1
fi

echo "Release created with ID: ${RELEASE_ID}"

# 上传二进制文件
echo ""
echo "========================================="
echo "Uploading binary to Release..."
echo "========================================="

BINARY_FILE="release/clean/terraform-provider-jumpserver_1.0.0_linux_amd64.zip"

if [ ! -f "${BINARY_FILE}" ]; then
    echo "Error: Binary file not found at ${BINARY_FILE}"
    exit 1
fi

echo "Uploading ${BINARY_FILE}..."
UPLOAD_RESPONSE=$(curl -s -X POST \
  -H "Authorization: token ${GITHUB_TOKEN}" \
  -H "Content-Type: application/zip" \
  -H "Content-Length: $(stat -f%z "${BINARY_FILE}" 2>/dev/null || stat -c%s "${BINARY_FILE}")" \
  --data-binary @"${BINARY_FILE}" \
  https://uploads.github.com/repos/${REPO}/releases/${RELEASE_ID}/assets?name=$(basename ${BINARY_FILE}))

echo "Upload response: ${UPLOAD_RESPONSE}"

echo ""
echo "========================================="
echo "Release completed successfully!"
echo "========================================="
echo "Release URL: https://github.com/${REPO}/releases/tag/${TAG_NAME}"
