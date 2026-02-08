# JumpServer Terraform Provider Development Guide

## Development Setup

1. Clone the repository:
```bash
git clone https://github.com/your-org/terraform-provider-jumpserver
cd terraform-provider-jumpserver
```

2. Install dependencies:
```bash
make deps
```

3. Build the provider:
```bash
make build
```

## Project Structure

```
terraform-provider-jumpserver/
├── main.go                          # Provider entry point
├── go.mod                           # Go module definition
├── Makefile                         # Build automation
├── internal/
│   ├── jumpserver/                  # JumpServer API client
│   │   ├── client.go                # HTTP client with auth
│   │   ├── asset.go                 # Asset operations
│   │   ├── account.go               # Account operations
│   │   ├── permission.go            # Permission operations
│   │   ├── platform.go              # Platform operations
│   │   ├── node.go                  # Node operations
│   │   └── user.go                  # User operations
│   └── provider/                    # Terraform provider implementation
│       ├── provider.go              # Provider configuration
│       ├── resources/               # Resource implementations
│       │   ├── asset_resource.go
│       │   ├── account_resource.go
│       │   ├── permission_resource.go
│       │   └── user_resource.go
│       └── data_sources/            # Data source implementations
│           ├── asset_data_source.go
│           ├── platform_data_source.go
│           ├── node_data_source.go
│           └── user_data_source.go
└── examples/                        # Example configurations
```

## Adding a New Resource

1. Create the resource file in `internal/provider/resources/`
2. Implement the required interfaces: `resource.Resource`, `resource.ResourceWithConfigure`, `resource.ResourceWithImportState`
3. Add API client methods in `internal/jumpserver/`
4. Register the resource in `provider.go`

## Adding a New Data Source

1. Create the data source file in `internal/provider/data_sources/`
2. Implement the required interfaces: `datasource.DataSource`, `datasource.DataSourceWithConfigure`
3. Add API client methods in `internal/jumpserver/`
4. Register the data source in `provider.go`

## Testing

### Unit Tests
```bash
make test
```

### Acceptance Tests
```bash
export JUMPSERVER_ENDPOINT="https://jumpserver.example.com"
export JUMPSERVER_KEY_ID="your-key-id"
export JUMPSERVER_KEY_SECRET="your-key-secret"
make testacc
```

## Code Quality

### Format Code
```bash
make fmt
```

### Run Linter
```bash
make lint
```

### Vet Code
```bash
make vet
```

## Building for Release

### Linux
```bash
GOOS=linux GOARCH=amd64 make build
```

### macOS
```bash
GOOS=darwin GOARCH=amd64 make build
```

### Windows
```bash
GOOS=windows GOARCH=amd64 make build
```

## Authentication

The provider uses JumpServer's HMAC-SHA256 authentication. The authentication flow:

1. Client creates HTTP request
2. Generate signature string: `(request-target): method path\ndate: timestamp`
3. Create HMAC-SHA256 signature using the secret key
4. Add Authorization header with signature

## State Synchronization

The provider maintains Terraform state synchronization with JumpServer:

- **Create**: Creates resource in JumpServer and stores ID in state
- **Read**: Fetches current state from JumpServer and updates Terraform state
- **Update**: Updates resource in JumpServer based on plan
- **Delete**: Removes resource from JumpServer and state

## Error Handling

All API errors are captured and surfaced as Terraform diagnostics:
- Use `resp.Diagnostics.AddError()` for blocking errors
- Use `resp.Diagnostics.AddWarning()` for non-blocking warnings
- Include context in error messages for easier debugging

## Sensitive Data

Sensitive fields are marked with `Sensitive: true` in the schema:
- Secrets in accounts are never displayed in plans or state
- Provider configuration fields (key_id, key_secret) are sensitive

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linting
5. Submit a pull request
