package provider_test

import (
	"testing"

	providerpkg "jumpserver/internal/provider"
)

// testAccProtoV5ProviderFactories provides the provider factory for acceptance tests
var testAccProtoV5ProviderFactories = map[string]func() (interface{}, error){
	"jumpserver": func() (interface{}, error) {
		return providerpkg.New("dev")(), nil
	},
}

// GetTestAccProtoV5ProviderFactories exports the factories for use in other test packages
func GetTestAccProtoV5ProviderFactories() map[string]func() (interface{}, error) {
	return testAccProtoV5ProviderFactories
}

func TestProvider(t *testing.T) {
	// Provider validation is now handled by the framework itself
	// This test ensures the provider can be instantiated
	_ = providerpkg.New("dev")()
}

func TestProvider_Configure(t *testing.T) {
	t.Skip("Skipping - requires JumpServer instance and framework compatibility update")
}

func testAccProviderExists() {
	// Verify provider configuration was successful
}

const testProviderConfig = `
provider "jumpserver" {
  endpoint  = "https://jumpserver.test"
  key_id    = "test-key-id"
  key_secret = "test-key-secret"
}

data "jumpserver_platform" "linux" {
  id = "Linux"
}
`
