package plugin

import (
	"fmt"

	"github.com/zipstack/pct-plugin-framework/fwhelpers"
	"github.com/zipstack/pct-plugin-framework/schema"

	"github.com/zipstack/pct-provider-airbyte-cloud/api"
)

// Provider implementation.
type Provider struct {
	Client           *api.Client
	ResourceServices map[string]string
}

// Model maps the provider state as per schema.
type ProviderModel struct {
	Host          string `cty:"host"`
	Authorization string `cty:"authorization"`
}

// Ensure the implementation satisfies the expected interfaces
var (
	_ schema.ProviderService = &Provider{}
)

// Helper function to return a provider service instance.
func NewProvider() schema.ProviderService {
	return &Provider{}
}

// Metadata returns the provider type name.
func (p *Provider) Metadata(req *schema.ServiceRequest) *schema.ServiceResponse {
	return &schema.ServiceResponse{
		TypeName: "airbyte",
	}
}

// Schema defines the provider-level schema for configuration data.
func (p *Provider) Schema() *schema.ServiceResponse {
	s := &schema.Schema{
		Description: "Airbyte provider plugin",
		Attributes: map[string]schema.Attribute{
			"host": &schema.StringAttribute{
				Description: "URI for Airbyte API. May also be provided via AIRBYTE_HOST environment variable.",
				Required:    true,
			},

			"authorization": &schema.StringAttribute{
				Description: "Bearer Token for airbyte provider",
				Required:    true,
				Sensitive:   true,
			},
		},
	}

	sEnc, err := fwhelpers.Encode(s)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	return &schema.ServiceResponse{
		SchemaContents: sEnc,
	}
}

func (p *Provider) Configure(req *schema.ServiceRequest) *schema.ServiceResponse {
	var pm ProviderModel
	err := fwhelpers.UnpackModel(req.ConfigContents, &pm)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	if pm.Host == "" || pm.Authorization == "" {
		return schema.ErrorResponse(fmt.Errorf(
			"invalid host or credentials received.\n" +
				"Provider is unable to create Airbyte API client.",
		))
	}

	if p.Client == nil {
		client, err := api.NewClient(pm.Host, pm.Authorization)
		if err != nil {
			return schema.ErrorResponse(err)
		}
		p.Client = client
	}

	// Make API creds available for Resource type Configure methods.
	creds := map[string]string{
		"host":          pm.Host,
		"authorization": pm.Authorization,
	}
	cEnc, err := fwhelpers.Encode(creds)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	return &schema.ServiceResponse{
		ResourceData: cEnc,
	}
}

func (p *Provider) Resources() *schema.ServiceResponse {
	return &schema.ServiceResponse{
		ResourceServices: p.ResourceServices,
	}
}

func (p *Provider) UpdateResourceServices(resServices map[string]string) {
	if resServices != nil {
		p.ResourceServices = resServices
	}
}
