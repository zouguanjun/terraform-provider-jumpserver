# Terraform Provider JumpServer v1.0.0 - 手动发布指南

## ✅ 已完成的工作

1. **编译优化的二进制文件**: `bin/terraform-provider-jumpserver_v1.0.0_x4` (18MB)
2. **打包 Release 文件**: `release/clean/terraform-provider-jumpserver_1.0.0_linux_amd64.zip` (6.1MB)
3. **创建 Git 标签**: `v1.0.0` (本地)
4. **准备 Release 说明**: `RELEASE_NOTES.md`

## 🔑 第一步: 获取有效的 GitHub Token

### 方法 1: 创建 Personal Access Token

1. 访问: https://github.com/settings/tokens
2. 点击 "Generate new token (classic)"
3. 勾选权限:
   - `repo` (完整的仓库访问权限)
   - `workflow` (如需 GitHub Actions)
4. 点击 "Generate token"
5. 复制生成的 token (格式: `ghp_xxxxxxxxxxxxxxxxx`)

### 方法 2: 使用 GitHub CLI (gh)

```bash
# 安装 GitHub CLI
# Ubuntu/Debian:
sudo apt install gh

# 登录
gh auth login

# 创建 token
gh auth token
```

## 📦 第二步: 推送标签到 GitHub

### 方法 1: 使用 Git Credential Helper

```bash
git config --global credential.helper store

# 推送时会提示输入用户名和 token
git push origin v1.0.0
```

### 方法 2: 使用 Personal Access Token

```bash
# 替换 YOUR_TOKEN 为你的实际 token
git push https://YOUR_TOKEN@github.com/zouguanjun/terraform-provider-jumpserver.git v1.0.0
```

### 方法 3: 使用 GitHub CLI

```bash
gh auth login
git push origin v1.0.0
```

## 🚀 第三步: 在 GitHub 创建 Release

### 方法 1: 通过 Web 界面 (最简单)

1. 访问: https://github.com/zouguanjun/terraform-provider-jumpserver/releases/new
2. 选择标签: `v1.0.0`
3. 标题: `v1.0.0 - Initial Stable Release`
4. 描述: 复制 `RELEASE_NOTES.md` 的内容
5. 勾选 "Set as the latest release"
6. 勾选 "Set as a pre-release" (如果这是测试版本,否则不选)
7. 点击 "Generate release notes" 或手动粘贴说明
8. 上传附件: 选择 `release/clean/terraform-provider-jumpserver_1.0.0_linux_amd64.zip`
9. 点击 "Publish release"

### 方法 2: 使用 GitHub CLI

```bash
# 创建 Release 并上传二进制文件
gh release create v1.0.0 \
  --title "v1.0.0 - Initial Stable Release" \
  --notes-file RELEASE_NOTES.md \
  release/clean/terraform-provider-jumpserver_1.0.0_linux_amd64.zip
```

### 方法 3: 使用 API (curl)

```bash
# 创建 Release
GITHUB_TOKEN="your_token_here"

# 创建 Release
RELEASE_RESPONSE=$(curl -s -X POST \
  -H "Authorization: token ${GITHUB_TOKEN}" \
  -H "Accept: application/vnd.github.v3+json" \
  https://api.github.com/repos/zouguanjun/terraform-provider-jumpserver/releases \
  -d '{
    "tag_name": "v1.0.0",
    "name": "v1.0.0 - Initial Stable Release",
    "body": "Initial stable release",
    "draft": false,
    "prerelease": false
  }')

# 获取 Release ID
RELEASE_ID=$(echo "${RELEASE_RESPONSE}" | grep -oP '"id":\s*\K[0-9]+')

# 上传二进制文件
curl -X POST \
  -H "Authorization: token ${GITHUB_TOKEN}" \
  -H "Content-Type: application/zip" \
  --data-binary @release/clean/terraform-provider-jumpserver_1.0.0_linux_amd64.zip \
  https://uploads.github.com/repos/zouguanjun/terraform-provider-jumpserver/releases/${RELEASE_ID}/assets?name=terraform-provider-jumpserver_1.0.0_linux_amd64.zip
```

## ✅ 验证 Release

发布完成后,访问以下地址验证:

- Release 页面: https://github.com/zouguanjun/terraform-provider-jumpserver/releases/tag/v1.0.0
- 检查是否包含二进制文件下载链接
- 检查 Release 说明是否正确显示

## 📊 当前文件状态

```
bin/
├── terraform-provider-jumpserver_v1.0.0_x4        # 优化后的二进制 (18MB)
└── terraform-provider-jumpserver                   # 原始二进制 (25MB)

release/clean/
├── terraform-provider-jumpserver_v1.0.0_x4         # 二进制文件
└── terraform-provider-jumpserver_1.0.0_linux_amd64.zip  # 打包文件 (6.1MB)

.git/
└── refs/tags/
    └── v1.0.0                                      # 本地标签已创建
```

## 🎯 推荐操作

1. **最快方式**: 使用 Web 界面创建 Release (方法1)
2. **自动化方式**: 使用 GitHub CLI (方法2)
3. **CI/CD 方式**: 使用 API (方法3,适合脚本化)

## 🔍 故障排除

### Token 权限问题
确保 token 包含 `repo` 权限,否则无法推送标签和创建 Release。

### 推送失败
```bash
# 检查远程仓库
git remote -v

# 检查标签
git tag -l

# 查看标签详情
git show v1.0.0
```

### Release 创建失败
- 检查仓库是否存在: https://github.com/zouguanjun/terraform-provider-jumpserver
- 确认 token 有效性: curl -H "Authorization: token YOUR_TOKEN" https://api.github.com/user

## 📞 需要帮助?

如果遇到问题,请检查:
1. GitHub token 是否有 `repo` 权限
2. 仓库是否可访问
3. 标签是否已推送
4. 二进制文件是否存在

---

**准备好后,请提供有效的 GitHub token,我可以帮你自动化完成所有步骤!**
