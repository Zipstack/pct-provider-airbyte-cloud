package plugin

import (
	"fmt"
	"time"

	"github.com/zipstack/pct-plugin-framework/fwhelpers"
	"github.com/zipstack/pct-plugin-framework/schema"

	"github.com/zipstack/pct-provider-airbyte-cloud/api"
)

// Resource implementation.
type sourceFreshdeskResource struct {
	Client *api.Client
}

type sourceFreshdeskResourceModel struct {
	Name                    string                         `pctsdk:"name"`
	SourceId                string                         `pctsdk:"source_id"`
	WorkspaceId             string                         `pctsdk:"workspace_id"`
	ConnectionConfiguration sourceFreshdeskConnConfigModel `pctsdk:"configuration"`
}

type sourceFreshdeskConnConfigModel struct {
	SourceType        string `pctsdk:"source_type"`
	StartDate         string `pctsdk:"start_date"`
	Domain            string `pctsdk:"domain"`
	ApiKey            string `pctsdk:"api_key"`
	RequestsPerMinute int    `pctsdk:"requests_per_minute"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ schema.ResourceService = &sourceFreshdeskResource{}
)

// Helper function to return a resource service instance.
func NewSourceFreshdeskResource() schema.ResourceService {
	return &sourceFreshdeskResource{}
}

// Metadata returns the resource type name.
// It is always provider name + "_" + resource type name.
func (r *sourceFreshdeskResource) Metadata(req *schema.ServiceRequest) *schema.ServiceResponse {
	return &schema.ServiceResponse{
		TypeName: req.TypeName + "_source_freshdesk",
	}
}

// Configure adds the provider configured client to the resource.
func (r *sourceFreshdeskResource) Configure(req *schema.ServiceRequest) *schema.ServiceResponse {
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
func (r *sourceFreshdeskResource) Schema() *schema.ServiceResponse {
	s := &schema.Schema{
		Description: "Source Freshdesk resource for Airbyte",
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
					"domain": &schema.StringAttribute{
						Description: "Domain",
						Required:    true,
					},
					"api_key": &schema.StringAttribute{
						Description: "Api Key",
						Required:    true,
						Sensitive:   true,
					},
					"requests_per_minute": &schema.IntAttribute{
						Description: "Requests Per Minute",
						Optional:    true,
						Required:    true,
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
func (r *sourceFreshdeskResource) Create(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	// Retrieve values from plan
	var plan sourceFreshdeskResourceModel
	err := fwhelpers.UnpackModel(req.PlanContents, &plan)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Generate API request body from plan
	body := api.SourceFreshdesk{}
	body.Name = plan.Name
	body.WorkspaceId = plan.WorkspaceId

	body.ConnectionConfiguration = api.SourceFreshdeskConnConfig{}
	body.ConnectionConfiguration.SourceType = plan.ConnectionConfiguration.SourceType
	body.ConnectionConfiguration.StartDate = plan.ConnectionConfiguration.StartDate
	body.ConnectionConfiguration.ApiKey = plan.ConnectionConfiguration.ApiKey
	body.ConnectionConfiguration.Domain = plan.ConnectionConfiguration.Domain
	body.ConnectionConfiguration.RequestsPerMinute = plan.ConnectionConfiguration.RequestsPerMinute

	// Create new source
	source, err := r.Client.CreateFreshdeskSource(body)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Update resource state with response body
	state := sourceFreshdeskResourceModel{}
	state.Name = source.Name
	state.SourceId = source.SourceId
	state.WorkspaceId = source.WorkspaceId

	state.ConnectionConfiguration = sourceFreshdeskConnConfigModel{}
	state.ConnectionConfiguration.SourceType = source.ConnectionConfiguration.SourceType
	state.ConnectionConfiguration.StartDate = source.ConnectionConfiguration.StartDate
	state.ConnectionConfiguration.ApiKey = source.ConnectionConfiguration.ApiKey
	state.ConnectionConfiguration.Domain = source.ConnectionConfiguration.Domain
	state.ConnectionConfiguration.RequestsPerMinute = source.ConnectionConfiguration.RequestsPerMinute

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
func (r *sourceFreshdeskResource) Read(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	var state sourceFreshdeskResourceModel

	// Get current state
	err := fwhelpers.UnpackModel(req.StateContents, &state)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	res := schema.ServiceResponse{}

	if req.StateID != "" {
		// Query using existing previous state.
		source, err := r.Client.ReadFreshdeskSource(req.StateID)
		if err != nil {
			return schema.ErrorResponse(err)
		}

		// Update state with refreshed value
		state.Name = source.Name
		state.SourceId = source.SourceId
		state.WorkspaceId = source.WorkspaceId

		state.ConnectionConfiguration = sourceFreshdeskConnConfigModel{}
		state.ConnectionConfiguration.SourceType = source.ConnectionConfiguration.SourceType
		state.ConnectionConfiguration.StartDate = source.ConnectionConfiguration.StartDate
		state.ConnectionConfiguration.ApiKey = source.ConnectionConfiguration.ApiKey
		state.ConnectionConfiguration.Domain = source.ConnectionConfiguration.Domain
		state.ConnectionConfiguration.RequestsPerMinute = source.ConnectionConfiguration.RequestsPerMinute

		res.StateID = state.SourceId
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

func (r *sourceFreshdeskResource) Update(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	// Retrieve values from plan
	var plan sourceFreshdeskResourceModel
	err := fwhelpers.UnpackModel(req.PlanContents, &plan)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Generate API request body from plan
	body := api.SourceFreshdesk{}
	body.Name = plan.Name
	body.SourceId = plan.SourceId

	body.ConnectionConfiguration.SourceType = plan.ConnectionConfiguration.SourceType
	body.ConnectionConfiguration.StartDate = plan.ConnectionConfiguration.StartDate
	body.ConnectionConfiguration.ApiKey = plan.ConnectionConfiguration.ApiKey
	body.ConnectionConfiguration.Domain = plan.ConnectionConfiguration.Domain
	body.ConnectionConfiguration.RequestsPerMinute = plan.ConnectionConfiguration.RequestsPerMinute

	// Update existing source
	_, err = r.Client.UpdateFreshdeskSource(body)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Fetch updated items
	source, err := r.Client.ReadFreshdeskSource(req.PlanID)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Update state with refreshed value
	state := sourceFreshdeskResourceModel{}
	state.Name = source.Name
	state.SourceId = source.SourceId
	state.WorkspaceId = source.WorkspaceId

	state.ConnectionConfiguration = sourceFreshdeskConnConfigModel{}
	state.ConnectionConfiguration.SourceType = source.ConnectionConfiguration.SourceType
	state.ConnectionConfiguration.StartDate = source.ConnectionConfiguration.StartDate
	state.ConnectionConfiguration.ApiKey = source.ConnectionConfiguration.ApiKey
	state.ConnectionConfiguration.Domain = source.ConnectionConfiguration.Domain
	state.ConnectionConfiguration.RequestsPerMinute = source.ConnectionConfiguration.RequestsPerMinute

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
func (r *sourceFreshdeskResource) Delete(req *schema.ServiceRequest) *schema.ServiceResponse {
	// Delete existing source
	err := r.Client.DeleteFreshdeskSource(req.StateID)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	return &schema.ServiceResponse{}
}
