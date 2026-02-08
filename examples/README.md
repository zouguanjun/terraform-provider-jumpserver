# JumpServer Terraform Provider Examples

This directory contains example Terraform configurations demonstrating how to use the JumpServer Terraform Provider.

## Basic Example

The `basic/` directory contains a complete example showing:

- How to configure the JumpServer provider
- How to create assets
- How to create users
- How to create accounts on assets
- How to configure permissions

### Running the Example

1. Copy the example variables file:
```bash
cd basic
cp terraform.tfvars.example terraform.tfvars
```

2. Edit `terraform.tfvars` with your JumpServer credentials and settings

3. Initialize Terraform:
```bash
terraform init
```

4. Review the plan:
```bash
terraform plan
```

5. Apply the configuration:
```bash
terraform apply
```

6. When done, destroy the resources:
```bash
terraform destroy
```

## Notes

- Ensure you have valid JumpServer credentials
- The `secret` field in `jumpserver_account` is sensitive and will not be displayed in plan output
- Always review the plan before applying
