#!/bin/bash

# 直接在 git 命令中嵌入 token
git config --global credential.helper store
git config --global user.name "zouguanjun"

# 方法1: 使用 netrc
cat > ~/.netrc << 'EOF'
machine github.com
login zouguanjun
password ghp_KNQyUQ0T7eX4Q6N4Z8yJcR0pWqMkLsX1vF2h
EOF
chmod 600 ~/.netrc

# 方法2: 更新远程 URL（备用）
git remote set-url origin https://ghp_KNQyUQ0T7eX4Q6N4Z8yJcR0pWqMkLsX1vF2h@github.com/zouguanjun/terraform-provider-jumpserver.git

# 推送标签
git push origin v1.0.0
