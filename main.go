package main

import (
	"github.com/zipstack/pct-plugin-framework/schema"
	"github.com/zipstack/pct-plugin-framework/server"

	"github.com/zipstack/pct-provider-airbyte-cloud/plugin"
)

// Set while building the compiled binary.
var version string

func main() {
	server.Serve(version, plugin.NewProvider, []func() schema.ResourceService{})
}
