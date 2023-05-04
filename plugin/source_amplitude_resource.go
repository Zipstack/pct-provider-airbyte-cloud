package plugin

import (
	"fmt"
	"time"

	"github.com/zipstack/pct-plugin-framework/fwhelpers"
	"github.com/zipstack/pct-plugin-framework/schema"

	"github.com/zipstack/pct-provider-airbyte-cloud/api"
)

// Resource implementation.
type sourceAmplitudeResource struct {
	Client *api.Client
}

type sourceAmplitudeResourceModel struct {
	Name                    string                         `pctsdk:"name"`
	SourceId                string                         `pctsdk:"source_id"`
	WorkspaceId             string                         `pctsdk:"workspace_id"`
	ConnectionConfiguration sourceAmplitudeConnConfigModel `pctsdk:"configuration"`
}

type sourceAmplitudeConnConfigModel struct {
	SourceType       string `pctsdk:"source_type"`
	StartDate        string `pctsdk:"start_date"`
	DataRegion       string `pctsdk:"data_region"`
	RequestTimeRange int    `pctsdk:"request_time_range"`
	ApiKey           string `pctsdk:"api_key"`
	SecretKey        string `pctsdk:"secret_key"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ schema.ResourceService = &sourceAmplitudeResource{}
)

// Helper function to return a resource service instance.
func NewSourceAmplitudeResource() schema.ResourceService {
	return &sourceAmplitudeResource{}
}

// Metadata returns the resource type name.
// It is always provider name + "_" + resource type name.
func (r *sourceAmplitudeResource) Metadata(req *schema.ServiceRequest) *schema.ServiceResponse {
	return &schema.ServiceResponse{
		TypeName: req.TypeName + "_source_amplitude",
	}
}

// Configure adds the provider configured client to the resource.
func (r *sourceAmplitudeResource) Configure(req *schema.ServiceRequest) *schema.ServiceResponse {
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
func (r *sourceAmplitudeResource) Schema() *schema.ServiceResponse {
	s := &schema.Schema{
		Description: "Source Amplitude resource for Airbyte",
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
					"data_region": &schema.StringAttribute{
						Description: "Date Region",
						Optional:    true,
						Required:    true,
					},
					"request_time_range": &schema.IntAttribute{
						Description: "Required time range",
						Optional:    true,
						Required:    true,
					},
					"secret_key": &schema.StringAttribute{
						Description: "Secret Key",
						Required:    true,
						Sensitive:   true,
					},
					"api_key": &schema.StringAttribute{
						Description: "Api Key",
						Required:    true,
						Sensitive:   true,
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
func (r *sourceAmplitudeResource) Create(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	// Retrieve values from plan
	var plan sourceAmplitudeResourceModel
	err := fwhelpers.UnpackModel(req.PlanContents, &plan)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Generate API request body from plan
	body := api.SourceAmplitude{}
	body.Name = plan.Name
	body.WorkspaceId = plan.WorkspaceId

	body.ConnectionConfiguration = api.SourceAmplitudeConnConfig{}
	body.ConnectionConfiguration.SourceType = plan.ConnectionConfiguration.SourceType
	body.ConnectionConfiguration.StartDate = plan.ConnectionConfiguration.StartDate
	body.ConnectionConfiguration.ApiKey = plan.ConnectionConfiguration.ApiKey
	body.ConnectionConfiguration.SecretKey = plan.ConnectionConfiguration.SecretKey
	body.ConnectionConfiguration.RequestTimeRange = plan.ConnectionConfiguration.RequestTimeRange
	body.ConnectionConfiguration.DataRegion = plan.ConnectionConfiguration.DataRegion

	// Create new source
	source, err := r.Client.CreateAmplitudeSource(body)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Update resource state with response body
	state := sourceAmplitudeResourceModel{}
	state.Name = source.Name
	state.SourceId = source.SourceId
	state.WorkspaceId = source.WorkspaceId

	state.ConnectionConfiguration = sourceAmplitudeConnConfigModel{}
	state.ConnectionConfiguration.SourceType = source.ConnectionConfiguration.SourceType
	state.ConnectionConfiguration.StartDate = source.ConnectionConfiguration.StartDate
	state.ConnectionConfiguration.ApiKey = source.ConnectionConfiguration.ApiKey
	state.ConnectionConfiguration.SecretKey = source.ConnectionConfiguration.SecretKey
	state.ConnectionConfiguration.DataRegion = source.ConnectionConfiguration.DataRegion
	state.ConnectionConfiguration.RequestTimeRange = source.ConnectionConfiguration.RequestTimeRange

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
func (r *sourceAmplitudeResource) Read(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	var state sourceAmplitudeResourceModel

	// Get current state
	err := fwhelpers.UnpackModel(req.StateContents, &state)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	res := schema.ServiceResponse{}

	if req.StateID != "" {
		// Query using existing previous state.
		source, err := r.Client.ReadAmplitudeSource(req.StateID)
		if err != nil {
			return schema.ErrorResponse(err)
		}

		// Update state with refreshed value
		state.Name = source.Name
		state.SourceId = source.SourceId
		state.WorkspaceId = source.WorkspaceId

		state.ConnectionConfiguration = sourceAmplitudeConnConfigModel{}
		state.ConnectionConfiguration.SourceType = source.ConnectionConfiguration.SourceType
		state.ConnectionConfiguration.StartDate = source.ConnectionConfiguration.StartDate
		state.ConnectionConfiguration.ApiKey = source.ConnectionConfiguration.ApiKey
		state.ConnectionConfiguration.SecretKey = source.ConnectionConfiguration.SecretKey
		state.ConnectionConfiguration.RequestTimeRange = source.ConnectionConfiguration.RequestTimeRange
		state.ConnectionConfiguration.DataRegion = source.ConnectionConfiguration.DataRegion

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

func (r *sourceAmplitudeResource) Update(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	// Retrieve values from plan
	var plan sourceAmplitudeResourceModel
	err := fwhelpers.UnpackModel(req.PlanContents, &plan)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Generate API request body from plan
	body := api.SourceAmplitude{}
	body.Name = plan.Name
	body.SourceId = plan.SourceId

	body.ConnectionConfiguration = api.SourceAmplitudeConnConfig{}
	body.ConnectionConfiguration.SourceType = plan.ConnectionConfiguration.SourceType
	body.ConnectionConfiguration.StartDate = plan.ConnectionConfiguration.StartDate
	body.ConnectionConfiguration.ApiKey = plan.ConnectionConfiguration.ApiKey
	body.ConnectionConfiguration.SecretKey = plan.ConnectionConfiguration.SecretKey
	body.ConnectionConfiguration.RequestTimeRange = plan.ConnectionConfiguration.RequestTimeRange
	body.ConnectionConfiguration.DataRegion = plan.ConnectionConfiguration.DataRegion

	// Update existing source
	_, err = r.Client.UpdateAmplitudeSource(body)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Fetch updated items
	source, err := r.Client.ReadAmplitudeSource(req.PlanID)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Update state with refreshed value
	state := sourceAmplitudeResourceModel{}
	state.Name = source.Name
	state.SourceId = source.SourceId
	state.WorkspaceId = source.WorkspaceId

	state.ConnectionConfiguration = sourceAmplitudeConnConfigModel{}
	state.ConnectionConfiguration.SourceType = source.ConnectionConfiguration.SourceType
	state.ConnectionConfiguration.StartDate = source.ConnectionConfiguration.StartDate
	state.ConnectionConfiguration.ApiKey = source.ConnectionConfiguration.ApiKey
	state.ConnectionConfiguration.SecretKey = source.ConnectionConfiguration.SecretKey
	state.ConnectionConfiguration.DataRegion = source.ConnectionConfiguration.DataRegion
	state.ConnectionConfiguration.RequestTimeRange = source.ConnectionConfiguration.RequestTimeRange

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
func (r *sourceAmplitudeResource) Delete(req *schema.ServiceRequest) *schema.ServiceResponse {
	// Delete existing source
	err := r.Client.DeleteAmplitudeSource(req.StateID)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	return &schema.ServiceResponse{}
}
