terraform {
  required_providers {
    jumpserver = {
      source  = "fit2cloud/jumpserver"
      version = ">= 1.0.0"
    }
  }
}

provider "jumpserver" {
  endpoint  = var.jumpserver_endpoint
  key_id    = var.jumpserver_key_id
  key_secret = var.jumpserver_key_secret
}

# Create a platform data source for reference
data "jumpserver_platform" "linux" {
  id = "Linux"
}

# Create a node data source for reference
data "jumpserver_node" "default" {
  id = "Default"
}

# Create an asset
resource "jumpserver_asset" "web_server" {
  name     = "Web Server"
  address  = "192.168.1.100"
  platform = data.jumpserver_platform.linux.name
  nodes    = [data.jumpserver_node.default.id]
  is_active = true
  comment  = "Main web server"
}

# Create a user
resource "jumpserver_user" "developer" {
  username = "developer"
  name     = "Developer User"
  email    = "developer@example.com"
  is_active = true
  comment  = "Developer account"
}

# Create an account on the asset
resource "jumpserver_account" "admin" {
  username    = "admin"
  asset       = jumpserver_asset.web_server.id
  secret      = var.admin_password
  secret_type = "password"
  comment     = "Admin account"
}

# Create a permission
resource "jumpserver_permission" "dev_access" {
  name    = "Developer Access"
  users   = [jumpserver_user.developer.id]
  assets  = [jumpserver_asset.web_server.id]
  actions = ["connect", "upload", "download"]
  comment = "Allow developers to connect and transfer files"
}

output "asset_id" {
  value = jumpserver_asset.web_server.id
}

output "user_id" {
  value = jumpserver_user.developer.id
}

output "permission_id" {
  value = jumpserver_permission.dev_access.id
}
