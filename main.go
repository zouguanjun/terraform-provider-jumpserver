package main

import (
	"context"
	"flag"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"jumpserver/internal/provider"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	providerserver.Serve(
		context.Background(),
		provider.New("1.0.0"),
		providerserver.ServeOpts{
			Address: "registry.terraform.io/jumpserver/jumpserver",
			Debug:   debug,
		},
	)
}
