package plugin

import (
	"fmt"
	"time"

	"github.com/zipstack/pct-plugin-framework/fwhelpers"
	"github.com/zipstack/pct-plugin-framework/schema"

	"github.com/zipstack/pct-provider-airbyte-cloud/api"
)

// Resource implementation.
type sourceGoogleSheetsResource struct {
	Client *api.Client
}

type sourceGoogleSheetsResourceModel struct {
	Name                    string                            `cty:"name"`
	SourceId                string                            `cty:"source_id"`
	WorkspaceId             string                            `cty:"workspace_id"`
	ConnectionConfiguration sourceGoogleSheetsConnConfigModel `cty:"configuration"`
}

type sourceGoogleSheetsConnConfigModel struct {
	SourceType    string                      `cty:"source_type"`
	RowBatchSize  int                         `cty:"row_batch_size"`
	SpreadsheetId string                      `cty:"spreadsheet_id"`
	Credentials   googleSheetsCredConfigModel `cty:"credentials"`
}
type googleSheetsCredConfigModel struct {
	AuthType           string `cty:"auth_type"`
	ServiceAccountInfo string `cty:"service_account_info"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ schema.ResourceService = &sourceGoogleSheetsResource{}
)

// Helper function to return a resource service instance.
func NewSourceGoogleSheetsResource() schema.ResourceService {
	return &sourceGoogleSheetsResource{}
}

// Metadata returns the resource type name.
// It is always provider name + "_" + resource type name.
func (r *sourceGoogleSheetsResource) Metadata(req *schema.ServiceRequest) *schema.ServiceResponse {
	return &schema.ServiceResponse{
		TypeName: req.TypeName + "_source_google_sheets",
	}
}

// Configure adds the provider configured client to the resource.
func (r *sourceGoogleSheetsResource) Configure(req *schema.ServiceRequest) *schema.ServiceResponse {
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
func (r *sourceGoogleSheetsResource) Schema() *schema.ServiceResponse {
	s := &schema.Schema{
		Description: "Source GoogleSheets resource for Airbyte",
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
					"row_batch_size": &schema.IntAttribute{
						Description: "Row Batch Size",
						Required:    false,
					},
					"spreadsheet_id": &schema.StringAttribute{
						Description: "Spreadsheet Id",
						Required:    true,
					},
					"credentials": &schema.MapAttribute{
						Description: "credentials",
						Required:    true,
						Attributes: map[string]schema.Attribute{
							"auth_type": &schema.StringAttribute{
								Description: "Auth Type",
								Required:    true,
							},
							"service_account_info": &schema.StringAttribute{
								Description: "Service Account Info",
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
func (r *sourceGoogleSheetsResource) Create(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	// Retrieve values from plan
	var plan sourceGoogleSheetsResourceModel
	err := fwhelpers.UnpackModel(req.PlanContents, &plan)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Generate API request body from plan
	body := api.SourceGoogleSheets{}
	body.Name = plan.Name
	body.WorkspaceId = plan.WorkspaceId

	body.ConnectionConfiguration = api.SourceGoogleSheetsConnConfig{}
	body.ConnectionConfiguration.SourceType = plan.ConnectionConfiguration.SourceType
	body.ConnectionConfiguration.RowBatchSize = plan.ConnectionConfiguration.RowBatchSize
	body.ConnectionConfiguration.SpreadsheetId = plan.ConnectionConfiguration.SpreadsheetId

	body.ConnectionConfiguration.Credentials = api.GoogleSheetsCredConfigModel{}
	body.ConnectionConfiguration.Credentials.AuthType = plan.ConnectionConfiguration.Credentials.AuthType
	body.ConnectionConfiguration.Credentials.ServiceAccountInfo = plan.ConnectionConfiguration.Credentials.ServiceAccountInfo

	// Create new source
	source, err := r.Client.CreateGoogleSheetsSource(body)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Update resource state with response body
	state := sourceGoogleSheetsResourceModel{}
	state.Name = source.Name
	state.SourceId = source.SourceId
	state.WorkspaceId = source.WorkspaceId

	state.ConnectionConfiguration = sourceGoogleSheetsConnConfigModel{}
	state.ConnectionConfiguration.SourceType = source.ConnectionConfiguration.SourceType
	state.ConnectionConfiguration.RowBatchSize = source.ConnectionConfiguration.RowBatchSize
	state.ConnectionConfiguration.SpreadsheetId = source.ConnectionConfiguration.SpreadsheetId

	state.ConnectionConfiguration.Credentials = googleSheetsCredConfigModel{}
	state.ConnectionConfiguration.Credentials.AuthType = source.ConnectionConfiguration.Credentials.AuthType
	state.ConnectionConfiguration.Credentials.ServiceAccountInfo = source.ConnectionConfiguration.Credentials.ServiceAccountInfo

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
func (r *sourceGoogleSheetsResource) Read(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	var state sourceGoogleSheetsResourceModel

	// Get current state
	err := fwhelpers.UnpackModel(req.StateContents, &state)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	res := schema.ServiceResponse{}

	if req.StateID != "" {
		// Query using existing previous state.
		source, err := r.Client.ReadGoogleSheetsSource(req.StateID)
		if err != nil {
			return schema.ErrorResponse(err)
		}

		// Update state with refreshed value
		state.Name = source.Name
		state.SourceId = source.SourceId
		state.WorkspaceId = source.WorkspaceId

		state.ConnectionConfiguration = sourceGoogleSheetsConnConfigModel{}
		state.ConnectionConfiguration.SourceType = source.ConnectionConfiguration.SourceType
		state.ConnectionConfiguration.RowBatchSize = source.ConnectionConfiguration.RowBatchSize
		state.ConnectionConfiguration.SpreadsheetId = source.ConnectionConfiguration.SpreadsheetId

		state.ConnectionConfiguration.Credentials = googleSheetsCredConfigModel{}
		state.ConnectionConfiguration.Credentials.AuthType = source.ConnectionConfiguration.Credentials.AuthType
		state.ConnectionConfiguration.Credentials.ServiceAccountInfo = source.ConnectionConfiguration.Credentials.ServiceAccountInfo

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

func (r *sourceGoogleSheetsResource) Update(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	// Retrieve values from plan
	var plan sourceGoogleSheetsResourceModel
	err := fwhelpers.UnpackModel(req.PlanContents, &plan)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Generate API request body from plan
	body := api.SourceGoogleSheets{}
	body.Name = plan.Name
	body.SourceId = plan.SourceId

	body.ConnectionConfiguration = api.SourceGoogleSheetsConnConfig{}
	body.ConnectionConfiguration.SourceType = plan.ConnectionConfiguration.SourceType
	body.ConnectionConfiguration.RowBatchSize = plan.ConnectionConfiguration.RowBatchSize
	body.ConnectionConfiguration.SpreadsheetId = plan.ConnectionConfiguration.SpreadsheetId

	body.ConnectionConfiguration.Credentials = api.GoogleSheetsCredConfigModel{}
	body.ConnectionConfiguration.Credentials.AuthType = plan.ConnectionConfiguration.Credentials.AuthType
	body.ConnectionConfiguration.Credentials.ServiceAccountInfo = plan.ConnectionConfiguration.Credentials.ServiceAccountInfo

	// Update existing source
	_, err = r.Client.UpdateGoogleSheetsSource(body)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Fetch updated items
	source, err := r.Client.ReadGoogleSheetsSource(req.PlanID)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Update state with refreshed value
	state := sourceGoogleSheetsResourceModel{}
	state.Name = source.Name
	state.SourceId = source.SourceId
	state.WorkspaceId = source.WorkspaceId

	state.ConnectionConfiguration = sourceGoogleSheetsConnConfigModel{}
	state.ConnectionConfiguration.SourceType = source.ConnectionConfiguration.SourceType
	state.ConnectionConfiguration.RowBatchSize = source.ConnectionConfiguration.RowBatchSize
	state.ConnectionConfiguration.SpreadsheetId = source.ConnectionConfiguration.SpreadsheetId

	state.ConnectionConfiguration.Credentials = googleSheetsCredConfigModel{}
	state.ConnectionConfiguration.Credentials.AuthType = source.ConnectionConfiguration.Credentials.AuthType
	state.ConnectionConfiguration.Credentials.ServiceAccountInfo = source.ConnectionConfiguration.Credentials.ServiceAccountInfo

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
func (r *sourceGoogleSheetsResource) Delete(req *schema.ServiceRequest) *schema.ServiceResponse {
	// Delete existing source
	err := r.Client.DeleteGoogleSheetsSource(req.StateID)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	return &schema.ServiceResponse{}
}
