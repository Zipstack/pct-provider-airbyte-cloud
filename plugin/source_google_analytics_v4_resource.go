package plugin

import (
	"fmt"
	"time"

	"github.com/zipstack/pct-plugin-framework/fwhelpers"
	"github.com/zipstack/pct-plugin-framework/schema"

	"github.com/zipstack/pct-provider-airbyte-cloud/api"
)

// Resource implementation.
type sourceGoogleAnalyticsV4Resource struct {
	Client *api.Client
}

type sourceGoogleAnalyticsV4ResourceModel struct {
	Name                    string                                 `cty:"name"`
	SourceId                string                                 `cty:"source_id"`
	WorkspaceId             string                                 `cty:"workspace_id"`
	ConnectionConfiguration sourceGoogleAnalyticsV4ConnConfigModel `cty:"configuration"`
}

type sourceGoogleAnalyticsV4ConnConfigModel struct {
	SourceType    string                           `cty:"source_type"`
	StartDate     string                           `cty:"start_date"`
	WindowInDays  int                              `cty:"window_in_days"`
	ViewId        string                           `cty:"view_id"`
	CustomReports string                           `cty:"custom_reports"`
	Credentials   googleAnalyticsV4CredConfigModel `cty:"credentials"`
}
type googleAnalyticsV4CredConfigModel struct {
	AuthType        string `cty:"auth_type"`
	CredentialsJson string `cty:"credentials_json"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ schema.ResourceService = &sourceGoogleAnalyticsV4Resource{}
)

// Helper function to return a resource service instance.
func NewSourceGoogleAnalyticsV4Resource() schema.ResourceService {
	return &sourceGoogleAnalyticsV4Resource{}
}

// Metadata returns the resource type name.
// It is always provider name + "_" + resource type name.
func (r *sourceGoogleAnalyticsV4Resource) Metadata(req *schema.ServiceRequest) *schema.ServiceResponse {
	return &schema.ServiceResponse{
		TypeName: req.TypeName + "_source_google_analytics_v4",
	}
}

// Configure adds the provider configured client to the resource.
func (r *sourceGoogleAnalyticsV4Resource) Configure(req *schema.ServiceRequest) *schema.ServiceResponse {
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
func (r *sourceGoogleAnalyticsV4Resource) Schema() *schema.ServiceResponse {
	s := &schema.Schema{
		Description: "Source GoogleAnalyticsV4 resource for Airbyte",
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
					"view_id": &schema.StringAttribute{
						Description: "View Id",
						Required:    true,
					},
					"window_in_days": &schema.IntAttribute{
						Description: "window in days",
						Required:    false,
					},
					"custom_reports": &schema.StringAttribute{
						Description: "custom reports",
						Required:    false,
					},
					"credentials": &schema.MapAttribute{
						Description: "credentials",
						Required:    true,
						Attributes: map[string]schema.Attribute{
							"auth_type": &schema.StringAttribute{
								Description: "Auth Type",
								Required:    true,
							},
							"credentials_json": &schema.StringAttribute{
								Description: "Credential Json",
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
func (r *sourceGoogleAnalyticsV4Resource) Create(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	// Retrieve values from plan
	var plan sourceGoogleAnalyticsV4ResourceModel
	err := fwhelpers.UnpackModel(req.PlanContents, &plan)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Generate API request body from plan
	body := api.SourceGoogleAnalyticsV4{}
	body.Name = plan.Name
	body.WorkspaceId = plan.WorkspaceId

	body.ConnectionConfiguration = api.SourceGoogleAnalyticsV4ConnConfig{}
	body.ConnectionConfiguration.SourceType = plan.ConnectionConfiguration.SourceType
	body.ConnectionConfiguration.StartDate = plan.ConnectionConfiguration.StartDate
	body.ConnectionConfiguration.CustomReports = plan.ConnectionConfiguration.CustomReports
	body.ConnectionConfiguration.ViewId = plan.ConnectionConfiguration.ViewId
	body.ConnectionConfiguration.WindowInDays = plan.ConnectionConfiguration.WindowInDays
	body.ConnectionConfiguration.Credentials = api.GoogleAnalyticsV4CredConfigModel{}
	body.ConnectionConfiguration.Credentials.AuthType = plan.ConnectionConfiguration.Credentials.AuthType
	body.ConnectionConfiguration.Credentials.CredentialsJson = plan.ConnectionConfiguration.Credentials.CredentialsJson

	// Create new source
	source, err := r.Client.CreateGoogleAnalyticsV4Source(body)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Update resource state with response body
	state := sourceGoogleAnalyticsV4ResourceModel{}
	state.Name = source.Name
	state.SourceId = source.SourceId
	state.WorkspaceId = source.WorkspaceId

	state.ConnectionConfiguration = sourceGoogleAnalyticsV4ConnConfigModel{}
	state.ConnectionConfiguration.SourceType = source.ConnectionConfiguration.SourceType
	state.ConnectionConfiguration.StartDate = source.ConnectionConfiguration.StartDate
	state.ConnectionConfiguration.CustomReports = source.ConnectionConfiguration.CustomReports
	state.ConnectionConfiguration.ViewId = source.ConnectionConfiguration.ViewId
	state.ConnectionConfiguration.WindowInDays = source.ConnectionConfiguration.WindowInDays
	state.ConnectionConfiguration.Credentials = googleAnalyticsV4CredConfigModel{}
	state.ConnectionConfiguration.Credentials.AuthType = source.ConnectionConfiguration.Credentials.AuthType
	state.ConnectionConfiguration.Credentials.CredentialsJson = source.ConnectionConfiguration.Credentials.CredentialsJson

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
func (r *sourceGoogleAnalyticsV4Resource) Read(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	var state sourceGoogleAnalyticsV4ResourceModel

	// Get current state
	err := fwhelpers.UnpackModel(req.StateContents, &state)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	res := schema.ServiceResponse{}

	if req.StateID != "" {
		// Query using existing previous state.
		source, err := r.Client.ReadGoogleAnalyticsV4Source(req.StateID)
		if err != nil {
			return schema.ErrorResponse(err)
		}

		// Update state with refreshed value
		state.Name = source.Name
		state.SourceId = source.SourceId
		state.WorkspaceId = source.WorkspaceId

		state.ConnectionConfiguration = sourceGoogleAnalyticsV4ConnConfigModel{}
		state.ConnectionConfiguration.SourceType = source.ConnectionConfiguration.SourceType
		state.ConnectionConfiguration.StartDate = source.ConnectionConfiguration.StartDate
		state.ConnectionConfiguration.CustomReports = source.ConnectionConfiguration.CustomReports
		state.ConnectionConfiguration.ViewId = source.ConnectionConfiguration.ViewId
		state.ConnectionConfiguration.WindowInDays = source.ConnectionConfiguration.WindowInDays
		state.ConnectionConfiguration.Credentials = googleAnalyticsV4CredConfigModel{}
		state.ConnectionConfiguration.Credentials.AuthType = source.ConnectionConfiguration.Credentials.AuthType
		state.ConnectionConfiguration.Credentials.CredentialsJson = source.ConnectionConfiguration.Credentials.CredentialsJson

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

func (r *sourceGoogleAnalyticsV4Resource) Update(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	// Retrieve values from plan
	var plan sourceGoogleAnalyticsV4ResourceModel
	err := fwhelpers.UnpackModel(req.PlanContents, &plan)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Generate API request body from plan
	body := api.SourceGoogleAnalyticsV4{}
	body.Name = plan.Name
	body.SourceId = plan.SourceId

	body.ConnectionConfiguration = api.SourceGoogleAnalyticsV4ConnConfig{}
	body.ConnectionConfiguration.SourceType = plan.ConnectionConfiguration.SourceType
	body.ConnectionConfiguration.StartDate = plan.ConnectionConfiguration.StartDate
	body.ConnectionConfiguration.CustomReports = plan.ConnectionConfiguration.CustomReports
	body.ConnectionConfiguration.ViewId = plan.ConnectionConfiguration.ViewId
	body.ConnectionConfiguration.WindowInDays = plan.ConnectionConfiguration.WindowInDays
	body.ConnectionConfiguration.Credentials = api.GoogleAnalyticsV4CredConfigModel{}
	body.ConnectionConfiguration.Credentials.AuthType = plan.ConnectionConfiguration.Credentials.AuthType
	body.ConnectionConfiguration.Credentials.CredentialsJson = plan.ConnectionConfiguration.Credentials.CredentialsJson

	// Update existing source
	_, err = r.Client.UpdateGoogleAnalyticsV4Source(body)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Fetch updated items
	source, err := r.Client.ReadGoogleAnalyticsV4Source(req.PlanID)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Update state with refreshed value
	state := sourceGoogleAnalyticsV4ResourceModel{}
	state.Name = source.Name
	state.SourceId = source.SourceId
	state.WorkspaceId = source.WorkspaceId

	state.ConnectionConfiguration = sourceGoogleAnalyticsV4ConnConfigModel{}
	state.ConnectionConfiguration.SourceType = source.ConnectionConfiguration.SourceType
	state.ConnectionConfiguration.StartDate = source.ConnectionConfiguration.StartDate
	state.ConnectionConfiguration.CustomReports = source.ConnectionConfiguration.CustomReports
	state.ConnectionConfiguration.ViewId = source.ConnectionConfiguration.ViewId
	state.ConnectionConfiguration.WindowInDays = source.ConnectionConfiguration.WindowInDays
	state.ConnectionConfiguration.Credentials = googleAnalyticsV4CredConfigModel{}
	state.ConnectionConfiguration.Credentials.AuthType = source.ConnectionConfiguration.Credentials.AuthType
	state.ConnectionConfiguration.Credentials.CredentialsJson = source.ConnectionConfiguration.Credentials.CredentialsJson

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
func (r *sourceGoogleAnalyticsV4Resource) Delete(req *schema.ServiceRequest) *schema.ServiceResponse {
	// Delete existing source
	err := r.Client.DeleteGoogleAnalyticsV4Source(req.StateID)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	return &schema.ServiceResponse{}
}
