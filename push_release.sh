#!/bin/bash

export GITHUB_TOKEN="ghp_KNQyUQ0T7eX4Q6N4Z8yJcR0pWqMkLsX1vF2h"

# 更新远程 URL
git remote set-url origin https://${GITHUB_TOKEN}@github.com/zouguanjun/terraform-provider-jumpserver.git

# 推送标签
git push origin v1.0.0

echo "Tag pushed successfully!"
