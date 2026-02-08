# Deployment Guide

## Quick Installation

### Local Build
```bash
git clone https://github.com/your-org/terraform-provider-jumpserver
cd terraform-provider-jumpserver
make build
make install
```

### Manual Install
```bash
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/your-org/jumpserver/1.0.0/linux_amd64/
cp terraform-provider-jumpserver ~/.terraform.d/plugins/registry.terraform.io/your-org/jumpserver/1.0.0/linux_amd64/
```

## CI/CD Integration

### GitHub Actions
```yaml
- name: Install Provider
  run: |
    curl -L -o terraform-provider-jumpserver https://github.com/your-org/terraform-provider-jumpserver/releases/download/v1.0.0/terraform-provider-jumpserver_1.0.0_linux_amd64.zip
    chmod +x terraform-provider-jumpserver
```

### GitLab CI
```yaml
before_script:
  - apk add curl
  - curl -L -o terraform-provider-jumpserver https://github.com/your-org/terraform-provider-jumpserver/releases/download/v1.0.0/terraform-provider-jumpserver_1.0.0_linux_amd64.zip
  - chmod +x terraform-provider-jumpserver
```

## Production Best Practices

1. **Credential Management**: Use environment variables or secret stores
2. **Remote State**: Use S3/Consul with locking
3. **Workspaces**: Separate dev/staging/prod environments
4. **Monitoring**: Enable TF_LOG for debugging
5. **Testing**: Run `terraform plan` before apply
