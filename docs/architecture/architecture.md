# JumpServer Terraform Provider Architecture

## Overview

The JumpServer Terraform Provider is built using Terraform Plugin Framework, which provides a modern, type-safe approach to building Terraform providers in Go.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                        Terraform                            │
│                     (CLI / Cloud)                           │
└─────────────────────────┬───────────────────────────────────┘
                          │
                          │ Plugin Protocol
                          │
┌─────────────────────────▼───────────────────────────────────┐
│           JumpServer Terraform Provider                      │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              Provider Configuration                  │  │
│  │  - Endpoint, Key ID, Key Secret, Org ID              │  │
│  └────────────────────┬─────────────────────────────────┘  │
│                       │                                      │
│  ┌────────────────────▼─────────────────────────────────┐  │
│  │           Terraform Plugin Framework                  │  │
│  │  - Resource CRUD Operations                           │  │
│  │  - Data Source Queries                                │  │
│  │  - State Management                                   │  │
│  └────────────────────┬─────────────────────────────────┘  │
│                       │                                      │
│  ┌────────────────────▼─────────────────────────────────┐  │
│  │         JumpServer API Client                         │  │
│  │  - HTTP Client                                        │  │
│  │  - HMAC-SHA256 Authentication                         │  │
│  │  - Request/Response Handling                          │  │
│  └────────────────────┬─────────────────────────────────┘  │
│                       │                                      │
│  ┌────────────────────▼─────────────────────────────────┐  │
│  │     API Layer (asset.go, account.go, etc.)            │  │
│  │  - Asset Operations                                   │  │
│  │  - Account Operations                                 │  │
│  │  - Permission Operations                              │  │
│  │  - User Operations                                    │  │
│  └────────────────────┬─────────────────────────────────┘  │
└───────────────────────┼───────────────────────────────────┘
                        │
                        │ HTTPS
                        │
┌───────────────────────▼───────────────────────────────────┐
│                    JumpServer API                           │
│                                                              │
│  ┌──────────────────┐  ┌──────────────────┐                │
│  │   Asset API      │  │   User API       │                │
│  └──────────────────┘  └──────────────────┘                │
│  ┌──────────────────┐  ┌──────────────────┐                │
│  │  Account API     │  │ Permission API   │                │
│  └──────────────────┘  └──────────────────┘                │
│  ┌──────────────────┐  ┌──────────────────┐                │
│  │  Platform API    │  │   Node API       │                │
│  └──────────────────┘  └──────────────────┘                │
└──────────────────────────────────────────────────────────────┘
```

## Component Overview

### 1. Provider Entry Point (`main.go`)
- Initializes the provider
- Serves the provider via plugin protocol
- Handles debug mode support

### 2. Provider Configuration (`provider.go`)
- Defines provider schema (endpoint, credentials)
- Configures the JumpServer API client
- Registers resources and data sources

### 3. API Client (`internal/jumpserver/client.go`)
- HTTP client with timeout management
- HMAC-SHA256 request signing
- Generic request/response handling
- Error handling and reporting

### 4. Resource Layer (`internal/provider/resources/`)
- Terraform resource implementations
- Schema definitions
- CRUD operations (Create, Read, Update, Delete)
- Import functionality

### 5. Data Source Layer (`internal/provider/data_sources/`)
- Terraform data source implementations
- Query operations
- State read-only access

### 6. API Layer (`internal/jumpserver/*.go`)
- JumpServer API specific operations
- Request/response models
- Data transformation

## Data Flow

### Create Operation

```
1. Terraform Plan
   ↓
2. Provider Validate & Plan
   ↓
3. User applies
   ↓
4. Provider Create()
   ↓
5. API Client POST request
   ↓
6. JumpServer creates resource
   ↓
7. API returns created resource
   ↓
8. Provider stores ID in state
```

### Read Operation

```
1. Terraform refresh
   ↓
2. Provider Read()
   ↓
3. API Client GET request
   ↓
4. JumpServer returns resource
   ↓
5. Provider updates state
```

### Update Operation

```
1. Terraform plan detects changes
   ↓
2. Provider Update()
   ↓
3. API Client PUT request
   ↓
4. JumpServer updates resource
   ↓
5. API returns updated resource
   ↓
6. Provider updates state
```

### Delete Operation

```
1. User removes resource from config
   ↓
2. Provider Delete()
   ↓
3. API Client DELETE request
   ↓
4. JumpServer deletes resource
   ↓
5. Provider removes from state
```

## State Management

The provider maintains a bidirectional sync between Terraform state and JumpServer:

- **Terraform State**: Stores resource IDs and attributes
- **JumpServer**: Single source of truth for actual resource state
- **Sync Strategy**: Read operation ensures state matches JumpServer

## Authentication Flow

```
1. Provider receives config
   ↓
2. Create API client with credentials
   ↓
3. For each API request:
   a. Generate timestamp
   b. Build signature string
   c. Create HMAC-SHA256 signature
   d. Add Authorization header
   e. Send request to JumpServer
```

## Error Handling

- **API Errors**: Converted to Terraform diagnostics
- **Validation Errors**: Caught during plan phase
- **Network Errors**: Retried or surfaced immediately
- **State Errors**: Prevent operations from proceeding

## Concurrency

- Each resource operation is independent
- Provider creates a single API client instance
- HTTP client handles concurrent requests safely
- Rate limiting handled by JumpServer API

## Extensibility

Adding new resources:

1. Add API methods in `internal/jumpserver/`
2. Create resource in `internal/provider/resources/`
3. Register in `provider.go`
4. Add tests

## Security Considerations

- Sensitive fields marked in schema
- Secrets never stored in plan output
- Credentials passed via secure config
- HTTPS enforced for API communication
- HMAC signatures prevent tampering
