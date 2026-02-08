# JumpServer Terraform Provider v1.0.0 Release Notes

## 📦 Release Information
- **Version**: 1.0.0
- **Release Date**: 2026-02-08
- **Status**: Initial Stable Release

## ✨ New Features

### Provider Configuration
- HMAC-SHA256 authentication support
- Customizable JumpServer endpoint
- Timeout configuration

### Resources

#### jumpserver_asset
- Manage JumpServer assets (servers, network devices, etc.)
- Support for IP address and hostname configuration
- Platform type specification (Linux, Windows, etc.)
- Full CRUD operations with import support

#### jumpserver_account
- Manage accounts for assets
- Support for password and SSH key authentication
- Asset association
- Full CRUD operations with import support

#### jumpserver_permission
- Manage user permissions
- Multi-user and multi-asset support
- Action-based permissions (connect, upload, download, etc.)
- Full CRUD operations with import support

#### jumpserver_user
- Manage JumpServer users
- Email, name, and username configuration
- Active/inactive status control
- Full CRUD operations with import support

### Data Sources

#### data.jumpserver_asset
- Query asset information
- Filter by asset name or ID

#### data.jumpserver_platform
- Query platform information
- Filter by platform name or ID

#### data.jumpserver_node
- Query node information
- Filter by node name or ID

#### data.jumpserver_user
- Query user information
- Filter by username or ID

## 🔧 Requirements

- Terraform >= 1.0.0
- Go >= 1.21 (for development)
- JumpServer instance with API access

## 📖 Installation

### Binary Installation
1. Download the release binary for your platform
2. Create the directory: `mkdir -p ~/.terraform.d/plugins/registry.terraform.io/fit2cloud/jumpserver/1.0.0/linux_amd64/`
3. Extract and place the binary in the directory

### From Terraform Registry (Recommended)
Add to your `versions.tf`:

```hcl
terraform {
  required_providers {
    jumpserver = {
      source  = "fit2cloud/jumpserver"
      version = "1.0.0"
    }
  }
}
```

## 🚀 Quick Start

```hcl
terraform {
  required_providers {
    jumpserver = {
      source  = "fit2cloud/jumpserver"
      version = "1.0.0"
    }
  }
}

provider "jumpserver" {
  endpoint   = "https://your-jumpserver.com"
  key_id     = "your-access-key-id"
  key_secret = "your-access-key-secret"
}

resource "jumpserver_user" "example" {
  username = "terraform-user"
  name     = "Terraform User"
  email    = "terraform@example.com"
  is_active = true
}

resource "jumpserver_asset" "example" {
  name    = "Web Server"
  address = "192.168.1.100"
  platform = "Linux"
}

resource "jumpserver_account" "example" {
  username    = "root"
  secret_type = "password"
  secret      = "secure-password"
  asset_id    = jumpserver_asset.example.id
}

resource "jumpserver_permission" "example" {
  name    = "Server Access"
  users   = [jumpserver_user.example.id]
  assets  = [jumpserver_asset.example.id]
  actions = ["connect"]
}
```

## 📚 Documentation

- Full documentation: [Project README](README.md)
- Examples: [examples/](examples/)
- Development guide: [DEVELOPMENT.md](DEVELOPMENT.md)

## 🐛 Known Issues

None reported in this release.

## 🤝 Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework)
- Powered by [JumpServer](https://jumpserver.org/)
