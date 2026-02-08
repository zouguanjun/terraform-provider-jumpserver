# Quick Start Guide

This guide will help you get started with the JumpServer Terraform Provider in 5 minutes.

## Prerequisites

- Terraform >= 1.0 installed
- JumpServer instance with API access
- JumpServer Access Key (Key ID and Secret)

## Step 1: Install the Provider

### Option A: Local Build

```bash
git clone https://github.com/your-org/terraform-provider-jumpserver
cd terraform-provider-jumpserver
make build
make install
```

### Option B: Use from Terraform Registry

Add to your `terraform` block:

```hcl
terraform {
  required_providers {
    jumpserver = {
      source  = "your-org/jumpserver"
      version = ">= 1.0.0"
    }
  }
}
```

## Step 2: Configure the Provider

Create a `main.tf` file:

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
}
```

Or use environment variables:

```bash
export JUMPSERVER_ENDPOINT="https://jumpserver.example.com"
export JUMPSERVER_KEY_ID="your-access-key-id"
export JUMPSERVER_KEY_SECRET="your-access-key-secret"
```

Then your `main.tf` becomes:

```hcl
terraform {
  required_providers {
    jumpserver = {
      source  = "your-org/jumpserver"
    }
  }
}

provider "jumpserver" {}
```

## Step 3: Initialize Terraform

```bash
terraform init
```

## Step 4: Create Your First Resource

Add to `main.tf`:

```hcl
resource "jumpserver_asset" "my_server" {
  name     = "My First Server"
  address  = "192.168.1.100"
  platform = "Linux"
  is_active = true
}
```

## Step 5: Plan and Apply

```bash
terraform plan
terraform apply
```

## Step 6: Verify

After applying, you can:

1. Check the Terraform state:
```bash
terraform state list
terraform show jumpserver_asset.my_server
```

2. Verify in JumpServer web interface:
   - Navigate to Asset Management
   - Find "My First Server"

## Common Examples

### Create Multiple Assets

```hcl
resource "jumpserver_asset" "web_servers" {
  for_each = {
    "web-1" = "192.168.1.101"
    "web-2" = "192.168.1.102"
    "web-3" = "192.168.1.103"
  }

  name    = each.key
  address = each.value
  platform = "Linux"
}
```

### Use Data Sources

```hcl
data "jumpserver_platform" "linux" {
  id = "Linux"
}

resource "jumpserver_asset" "server" {
  name     = "Server"
  address  = "192.168.1.100"
  platform = data.jumpserver_platform.linux.name
}
```

### Create Complete Access Setup

```hcl
# Create user
resource "jumpserver_user" "developer" {
  username = "developer"
  name     = "Developer"
  email    = "dev@example.com"
}

# Create asset
resource "jumpserver_asset" "server" {
  name     = "Dev Server"
  address  = "192.168.1.100"
  platform = "Linux"
}

# Create permission
resource "jumpserver_permission" "dev_access" {
  name    = "Dev Access"
  users   = [jumpserver_user.developer.id]
  assets  = [jumpserver_asset.server.id]
  actions = ["connect"]
}
```

## Troubleshooting

### Authentication Error

```
Error: Could not configure provider
```

**Solution**: Verify your credentials and endpoint URL are correct.

### Network Error

```
Error: API request failed
```

**Solution**: Check network connectivity to JumpServer API endpoint.

### Resource Not Found

```
Error: Could not read asset
```

**Solution**: The resource may have been deleted outside Terraform. Use `terraform import` or `terraform refresh`.

## Next Steps

- Read the [full documentation](../README.md)
- Explore [examples](../examples/)
- Check [development guide](../DEVELOPMENT.md)

## Getting Help

- Open an issue on GitHub
- Check existing issues for similar problems
- Review Terraform logs with `TF_LOG=DEBUG terraform apply`
