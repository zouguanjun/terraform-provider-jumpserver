# 完整测试配置文件 - 测试所有 JumpServer Provider 功能
# 使用组织 ID: 00000000-0000-0000-000000000002

terraform {
  required_providers {
    jumpserver = {
      source  = "fit2cloud/jumpserver"
    }
  }
}

provider "jumpserver" {
  endpoint   = "https://10.1.14.25"
  key_id     = "bb3cd120-585a-40ad-9a20-09991187ccc9"
  key_secret = "XPJ39dgHt6pPtlgmOsMsuU8FMPwBjRxAye0b"
  org_id     = "00000000-0000-0000-000000000002"
  insecure_skip_verify = true
}

# ============================================
# 数据源测试
# ============================================

# 测试平台数据源
data "jumpserver_platform" "linux" {
  id = "Linux"
}

# 测试默认节点数据源
data "jumpserver_node" "default" {
  id = "/Default"
}

# ============================================
# 资源测试
# ============================================

# 1. 创建资产 (Asset)
resource "jumpserver_asset" "test_server" {
  name     = "测试服务器-tf-${formatdate("YYYYMMDDhhmmss", timestamp())}"
  address  = "192.168.200.100"
  platform = data.jumpserver_platform.linux.name
  nodes    = [data.jumpserver_node.default.id]
  is_active = true
  comment  = "Terraform 创建的测试资产"
}

# 2. 创建用户 (User)
resource "jumpserver_user" "test_user" {
  username = "tf_test_user_${formatdate("YYYYMMDDhhmmss", timestamp())}"
  name     = "Terraform测试用户"
  email    = "tf_test_${formatdate("YYYYMMDDhhmmss", timestamp())}@example.com"
  is_active = true
  comment  = "通过 Terraform 创建的测试用户"
}

# 3. 创建账号 (Account) - 使用密码类型
resource "jumpserver_account" "test_password_account" {
  username    = "root"
  asset       = jumpserver_asset.test_server.id
  secret      = "Test@Password12345"
  secret_type = "password"
  comment     = "Terraform 创建的密码账号"
}

# 4. 创建权限 (Permission)
resource "jumpserver_permission" "test_permission" {
  name    = "测试权限-${formatdate("YYYYMMDDhhmmss", timestamp())}"
  users   = [jumpserver_user.test_user.id]
  assets  = [jumpserver_asset.test_server.id]
  actions = ["connect"]
  comment = "Terraform 创建的测试权限"
}

# ============================================
# 输出测试结果
# ============================================

output "asset_info" {
  value = {
    id       = jumpserver_asset.test_server.id
    name     = jumpserver_asset.test_server.name
    address  = jumpserver_asset.test_server.address
    platform = jumpserver_asset.test_server.platform
  }
}

output "user_info" {
  value = {
    id       = jumpserver_user.test_user.id
    username = jumpserver_user.test_user.username
    name     = jumpserver_user.test_user.name
    email    = jumpserver_user.test_user.email
  }
}

output "account_info" {
  value = {
    id         = jumpserver_account.test_password_account.id
    username   = jumpserver_account.test_password_account.username
    secret_type = jumpserver_account.test_password_account.secret_type
    asset_id   = jumpserver_asset.test_server.id
  }
}

output "permission_info" {
  value = {
    id      = jumpserver_permission.test_permission.id
    name    = jumpserver_permission.test_permission.name
    users   = jumpserver_permission.test_permission.users
    assets  = jumpserver_permission.test_permission.assets
    actions = jumpserver_permission.test_permission.actions
  }
}

output "platform_info" {
  value = {
    id   = data.jumpserver_platform.linux.id
    name = data.jumpserver_platform.linux.name
  }
}

output "node_info" {
  value = {
    id        = data.jumpserver_node.default.id
    full_name = data.jumpserver_node.default.full_name
  }
}
