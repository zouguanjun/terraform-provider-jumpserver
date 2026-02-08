# Terraform Provider for JumpServer

A Terraform provider for managing JumpServer resources. This provider enables you to use Terraform to declaratively manage JumpServer assets, accounts, permissions, and users.

## Features

- **Asset Management**: Create, read, update, and delete JumpServer assets (servers, devices, etc.)
- **Account Management**: Manage accounts on assets
- **Permission Management**: Configure access permissions for users and user groups
- **User Management**: Create and manage JumpServer users
- **Data Sources**: Query JumpServer resources for use in your Terraform configurations

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.25

## Building the Provider

1. Clone the repository:
```bash
git clone https://github.com/your-org/terraform-provider-jumpserver
cd terraform-provider-jumpserver
```

2. Build the provider:
```bash
make build
```

3. Install the provider locally:
```bash
make install
```

## Usage

### Provider Configuration

```hcl
terraform {
  required_providers {
    jumpserver = {
      source  = "your-org/jumpserver"
      version = ">= 1.0.0"
    }
  }
}

provider "jumpserver" {
  endpoint  = "https://jumpserver.example.com"
  key_id    = "your-access-key-id"
  key_secret = "your-access-key-secret"
  org_id    = "00000000-0000-0000-0000-000000000000"  # Optional
}
```

### Environment Variables

Alternatively, you can use environment variables:

```bash
export JUMPSERVER_ENDPOINT="https://jumpserver.example.com"
export JUMPSERVER_KEY_ID="your-access-key-id"
export JUMPSERVER_KEY_SECRET="your-access-key-secret"
export JUMPSERVER_ORG_ID="00000000-0000-0000-0000-000000000000"
```

### Example: Managing Assets

```hcl
# Create an asset
resource "jumpserver_asset" "web_server" {
  name     = "Web Server"
  address  = "192.168.1.100"
  platform = "Linux"
  nodes    = ["Default/Production"]
  is_active = true
  comment  = "Main web server"
}

# Query an existing asset
data "jumpserver_asset" "existing" {
  id = jumpserver_asset.web_server.id
}
```

### Example: Managing Accounts

```hcl
resource "jumpserver_asset" "server" {
  name    = "Production Server"
  address = "192.168.1.200"
  platform = "Linux"
}

resource "jumpserver_account" "admin" {
  username    = "admin"
  asset       = jumpserver_asset.server.id
  secret      = var.admin_password
  secret_type = "password"
  comment     = "Admin account"
}
```

### Example: Managing Users

```hcl
resource "jumpserver_user" "developer" {
  username = "developer"
  name     = "Developer User"
  email    = "developer@example.com"
  is_active = true
  comment  = "Developer account"
}
```

### Example: Managing Permissions

```hcl
resource "jumpserver_permission" "dev_access" {
  name   = "Developer Access"
  users  = [jumpserver_user.developer.id]
  assets = [jumpserver_asset.server.id]
  actions = ["connect", "upload", "download"]
  comment = "Allow developers to connect and transfer files"
}
```

## Resources

- `jumpserver_asset` - Manage JumpServer assets
- `jumpserver_account` - Manage accounts on assets
- `jumpserver_permission` - Manage access permissions
- `jumpserver_user` - Manage JumpServer users

## Data Sources

- `jumpserver_asset` - Query asset information
- `jumpserver_platform` - Query platform information
- `jumpserver_node` - Query organization node information
- `jumpserver_user` - Query user information

## Authentication

This provider uses JumpServer's Access Key authentication. You can generate access keys in the JumpServer web interface under Settings > Access Keys.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues and questions, please open an issue on GitHub.
