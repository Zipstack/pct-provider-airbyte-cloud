package main

import (
	"github.com/zipstack/pct-plugin-framework/schema"
	"github.com/zipstack/pct-plugin-framework/server"

	"github.com/zipstack/pct-provider-airbyte-cloud/plugin"
)

// Set while building the compiled binary.
var version string

func main() {
	server.Serve(version, plugin.NewProvider, []func() schema.ResourceService{

		plugin.NewSourcePipedriveResource,
		plugin.NewSourceStripeResource,
		plugin.NewSourceAmplitudeResource,
		plugin.NewSourceShopifyResource,
		plugin.NewSourceFreshdeskResource,
		plugin.NewSourceZendeskSupportResource,
		plugin.NewSourceHubspotResource,
		plugin.NewSourceGoogleAnalyticsV4Resource,
		plugin.NewSourceGoogleSheetsResource,
		plugin.NewSourceFacebookMarketingResource,
		plugin.NewDestinationMysqlResource,
		plugin.NewConnectionResource,
	})
}
