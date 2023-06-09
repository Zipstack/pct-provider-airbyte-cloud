package plugin

import (
	"fmt"
	"time"

	"github.com/zipstack/pct-plugin-framework/fwhelpers"
	"github.com/zipstack/pct-plugin-framework/schema"

	"github.com/zipstack/pct-provider-airbyte-cloud/api"
)

// Resource implementation.
type sourceShopifyResource struct {
	Client *api.Client
}

type sourceShopifyResourceModel struct {
	Name                    string                  `pctsdk:"name"`
	SourceId                string                  `pctsdk:"source_id"`
	WorkspaceId             string                  `pctsdk:"workspace_id"`
	ConnectionConfiguration SourceShopifyConnConfig `pctsdk:"configuration"`
}

type SourceShopifyConnConfig struct {
	SourceType  string                 `pctsdk:"source_type"`
	StartDate   string                 `pctsdk:"start_date"`
	Shop        string                 `pctsdk:"shop"`
	Credentials ShopifyCredConfigModel `pctsdk:"credentials"`
}
type ShopifyCredConfigModel struct {
	AuthMethod  string `pctsdk:"auth_method"`
	ApiPassword string `pctsdk:"api_password"`
	// ClientSecret string `pctsdk:"client_secret"`
	// AccessToken  string `pctsdk:"access_token"`
	// ClientId     string `pctsdk:"client_id"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ schema.ResourceService = &sourceShopifyResource{}
)

// Helper function to return a resource service instance.
func NewSourceShopifyResource() schema.ResourceService {
	return &sourceShopifyResource{}
}

// Metadata returns the resource type name.
// It is always provider name + "_" + resource type name.
func (r *sourceShopifyResource) Metadata(req *schema.ServiceRequest) *schema.ServiceResponse {
	return &schema.ServiceResponse{
		TypeName: req.TypeName + "_source_shopify",
	}
}

// Configure adds the provider configured client to the resource.
func (r *sourceShopifyResource) Configure(req *schema.ServiceRequest) *schema.ServiceResponse {
	if req.ResourceData == "" {
		return schema.ErrorResponse(fmt.Errorf("no data provided to configure resource"))
	}

	var creds map[string]string
	err := fwhelpers.Decode(req.ResourceData, &creds)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	client, err := api.NewClient(
		creds["host"], creds["authorization"],
	)
	if err != nil {
		return schema.ErrorResponse(fmt.Errorf("malformed data provided to configure resource"))
	}

	r.Client = client

	return &schema.ServiceResponse{}
}

// Schema defines the schema for the resource.
func (r *sourceShopifyResource) Schema() *schema.ServiceResponse {
	s := &schema.Schema{
		Description: "Source Shopify resource for Airbyte",
		Attributes: map[string]schema.Attribute{
			"name": &schema.StringAttribute{
				Description: "Name",
				Required:    true,
			},
			"source_id": &schema.StringAttribute{
				Description: "Source ID",
				Required:    false,
				Computed:    true,
			},
			"workspace_id": &schema.StringAttribute{
				Description: "Workspace ID",
				Required:    true,
			},
			"configuration": &schema.MapAttribute{
				Description: "Connection configuration",
				Required:    true,
				//Sensitive:   true,
				Attributes: map[string]schema.Attribute{
					"source_type": &schema.StringAttribute{
						Description: "Source Type",
						Required:    true,
					},
					"start_date": &schema.StringAttribute{
						Description: "Start Date",
						Required:    true,
					},
					"shop": &schema.StringAttribute{
						Description: "Shop",
						Required:    true,
					},
					"credentials": &schema.MapAttribute{
						Description: "credentials",
						Required:    true,
						Attributes: map[string]schema.Attribute{
							"auth_method": &schema.StringAttribute{
								Description: "auth_method",
								Required:    true,
							},
							"api_password": &schema.StringAttribute{
								Description: "api_password",
								Required:    true,
								Sensitive:   true,
							},
						},
					},
				},
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

// Create a new resource
func (r *sourceShopifyResource) Create(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	// Retrieve values from plan
	var plan sourceShopifyResourceModel
	err := fwhelpers.UnpackModel(req.PlanContents, &plan)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Generate API request body from plan
	body := api.SourceShopify{}
	body.Name = plan.Name
	body.WorkspaceId = plan.WorkspaceId

	body.ConnectionConfiguration = api.SourceShopifyConnConfig{}
	body.ConnectionConfiguration.SourceType = plan.ConnectionConfiguration.SourceType
	body.ConnectionConfiguration.StartDate = plan.ConnectionConfiguration.StartDate
	body.ConnectionConfiguration.Shop = plan.ConnectionConfiguration.Shop
	body.ConnectionConfiguration.Credentials = api.ShopifyCredConfigModel{}
	body.ConnectionConfiguration.Credentials.ApiPassword = plan.ConnectionConfiguration.Credentials.ApiPassword
	body.ConnectionConfiguration.Credentials.AuthMethod = plan.ConnectionConfiguration.Credentials.AuthMethod

	// Create new source
	source, err := r.Client.CreateShopifySource(body)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Update resource state with response body
	state := sourceShopifyResourceModel{}
	state.Name = source.Name
	state.SourceId = source.SourceId
	state.WorkspaceId = source.WorkspaceId

	state.ConnectionConfiguration = SourceShopifyConnConfig{}
	state.ConnectionConfiguration.SourceType = plan.ConnectionConfiguration.SourceType
	state.ConnectionConfiguration.StartDate = plan.ConnectionConfiguration.StartDate
	state.ConnectionConfiguration.Shop = plan.ConnectionConfiguration.Shop
	state.ConnectionConfiguration.Credentials = ShopifyCredConfigModel{}
	state.ConnectionConfiguration.Credentials.ApiPassword = plan.ConnectionConfiguration.Credentials.ApiPassword
	state.ConnectionConfiguration.Credentials.AuthMethod = plan.ConnectionConfiguration.Credentials.AuthMethod

	// Set refreshed state
	stateEnc, err := fwhelpers.PackModel(nil, &state)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	return &schema.ServiceResponse{
		StateID:          state.SourceId,
		StateContents:    stateEnc,
		StateLastUpdated: time.Now().Format(time.RFC850),
	}
}

// Read resource information
func (r *sourceShopifyResource) Read(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	var state sourceShopifyResourceModel

	// Get current state
	err := fwhelpers.UnpackModel(req.StateContents, &state)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	res := schema.ServiceResponse{}

	if req.StateID != "" {
		// Query using existing previous state.
		source, err := r.Client.ReadShopifySource(req.StateID)
		if err != nil {
			return schema.ErrorResponse(err)
		}

		// Update state with refreshed value
		state.Name = source.Name
		state.SourceId = source.SourceId
		state.WorkspaceId = source.WorkspaceId

		res.StateID = state.SourceId
		// Retaining other attributes from state itself as Reading resource have only 4 attributes in response
	} else {
		// No previous state exists.
		res.StateID = ""
	}

	// Set refreshed state
	stateEnc, err := fwhelpers.PackModel(nil, &state)
	if err != nil {
		return schema.ErrorResponse(err)
	}
	res.StateContents = stateEnc

	return &res
}

func (r *sourceShopifyResource) Update(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	// Retrieve values from plan
	var plan sourceShopifyResourceModel
	err := fwhelpers.UnpackModel(req.PlanContents, &plan)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Generate API request body from plan
	body := api.SourceShopify{}
	body.Name = plan.Name
	body.SourceId = plan.SourceId

	body.ConnectionConfiguration = api.SourceShopifyConnConfig{}
	body.ConnectionConfiguration.SourceType = plan.ConnectionConfiguration.SourceType
	body.ConnectionConfiguration.StartDate = plan.ConnectionConfiguration.StartDate
	body.ConnectionConfiguration.Shop = plan.ConnectionConfiguration.Shop
	body.ConnectionConfiguration.Credentials = api.ShopifyCredConfigModel{}
	body.ConnectionConfiguration.Credentials.ApiPassword = plan.ConnectionConfiguration.Credentials.ApiPassword
	body.ConnectionConfiguration.Credentials.AuthMethod = plan.ConnectionConfiguration.Credentials.AuthMethod

	// Update existing source
	_, err = r.Client.UpdateShopifySource(body)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Fetch updated items
	source, err := r.Client.ReadShopifySource(req.PlanID)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Update state with refreshed value
	state := sourceShopifyResourceModel{}
	state.Name = source.Name
	state.SourceId = source.SourceId
	state.WorkspaceId = source.WorkspaceId

	state.ConnectionConfiguration = SourceShopifyConnConfig{}
	state.ConnectionConfiguration.SourceType = plan.ConnectionConfiguration.SourceType
	state.ConnectionConfiguration.StartDate = plan.ConnectionConfiguration.StartDate
	state.ConnectionConfiguration.Shop = plan.ConnectionConfiguration.Shop
	state.ConnectionConfiguration.Credentials = ShopifyCredConfigModel{}
	state.ConnectionConfiguration.Credentials.ApiPassword = plan.ConnectionConfiguration.Credentials.ApiPassword
	state.ConnectionConfiguration.Credentials.AuthMethod = plan.ConnectionConfiguration.Credentials.AuthMethod

	// Set refreshed state
	stateEnc, err := fwhelpers.PackModel(nil, &state)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	return &schema.ServiceResponse{
		StateID:          state.SourceId,
		StateContents:    stateEnc,
		StateLastUpdated: time.Now().Format(time.RFC850),
	}
}

// Delete deletes the resource and removes the state on success.
func (r *sourceShopifyResource) Delete(req *schema.ServiceRequest) *schema.ServiceResponse {
	// Delete existing source
	err := r.Client.DeleteShopifySource(req.StateID)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	return &schema.ServiceResponse{}
}
