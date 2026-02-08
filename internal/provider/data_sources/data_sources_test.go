package data_sources_test

import (
	"jumpserver/internal/provider"
)

// testAccProtoV5ProviderFactories provides the provider factory for data source tests
var testAccProtoV5ProviderFactories = map[string]func() (interface{}, error){
	"jumpserver": func() (interface{}, error) {
		return provider.New("dev")(), nil
	},
}
