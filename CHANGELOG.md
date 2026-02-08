# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2025-01-24

### Added
- Initial release of the JumpServer Terraform Provider
- Asset resource: Manage JumpServer assets (servers, devices, etc.)
- Account resource: Manage accounts on assets
- Permission resource: Configure access permissions
- User resource: Manage JumpServer users
- Asset data source: Query asset information
- Platform data source: Query platform information
- Node data source: Query organization node information
- User data source: Query user information
- HMAC-SHA256 authentication support
- Organization ID support for multi-tenant deployments
- Import functionality for all resources

### Features
- Full CRUD operations for all resources
- State synchronization with JumpServer
- Sensitive data handling for secrets
- Comprehensive error handling and logging
- Support for environment variables for provider configuration

[Unreleased]: https://github.com/your-org/terraform-provider-jumpserver/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/your-org/terraform-provider-jumpserver/releases/tag/v1.0.0
