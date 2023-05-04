package plugin

import (
	"fmt"
	"time"

	"github.com/zipstack/pct-plugin-framework/fwhelpers"
	"github.com/zipstack/pct-plugin-framework/schema"

	"github.com/zipstack/pct-provider-airbyte-cloud/api"
)

// Resource implementation.
type sourceFacebookMarketingResource struct {
	Client *api.Client
}

type sourceFacebookMarketingResourceModel struct {
	Name                    string                                 `pctsdk:"name"`
	SourceId                string                                 `pctsdk:"source_id"`
	WorkspaceId             string                                 `pctsdk:"workspace_id"`
	ConnectionConfiguration sourceFacebookMarketingConnConfigModel `pctsdk:"configuration"`
}

type sourceFacebookMarketingConnConfigModel struct {
	SourceType  string `pctsdk:"source_type"`
	AccountId   string `pctsdk:"account_id"`
	StartDate   string `pctsdk:"start_date"`
	AccessToken string `pctsdk:"access_token"`

	EndDate              string `pctsdk:"end_date"`
	IncludeDeleted       bool   `pctsdk:"include_deleted"`
	FetchThumbnailImages bool   `pctsdk:"fetch_thumbnail_images"`
	//CustomInsights       any    `pctsdk:"custom_insights"`
	PageSize                   int  `pctsdk:"page_size"`
	InsightsLookbackWindow     int  `pctsdk:"insights_lookback_window"`
	MaxBatchSize               int  `pctsdk:"max_batch_size"`
	ActionBreakdownsAllowEmpty bool `pctsdk:"action_breakdowns_allow_empty"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ schema.ResourceService = &sourceFacebookMarketingResource{}
)

// Helper function to return a resource service instance.
func NewSourceFacebookMarketingResource() schema.ResourceService {
	return &sourceFacebookMarketingResource{}
}

// Metadata returns the resource type name.
// It is always provider name + "_" + resource type name.
func (r *sourceFacebookMarketingResource) Metadata(req *schema.ServiceRequest) *schema.ServiceResponse {
	return &schema.ServiceResponse{
		TypeName: req.TypeName + "_source_facebook_marketing",
	}
}

// Configure adds the provider configured client to the resource.
func (r *sourceFacebookMarketingResource) Configure(req *schema.ServiceRequest) *schema.ServiceResponse {
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
func (r *sourceFacebookMarketingResource) Schema() *schema.ServiceResponse {
	s := &schema.Schema{
		Description: "Source FacebookMarketing resource for Airbyte",
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
					"account_id": &schema.StringAttribute{
						Description: "account id",
						Required:    true,
					},
					"start_date": &schema.StringAttribute{
						Description: "start date",
						Required:    true,
					},
					"access_token": &schema.StringAttribute{
						Description: "Access Token",
						Required:    true,
						Sensitive:   true,
					},
					"end_date": &schema.StringAttribute{
						Description: "end date",
						Optional:    true,
						Required:    true,
					},
					"include_deleted": &schema.BoolAttribute{
						Description: "Include Deleted",
						Optional:    true,
						Required:    true,
					},
					"fetch_thumbnail_images": &schema.BoolAttribute{
						Description: "Fetch Thumbnail Image",
						Optional:    true,
						Required:    true,
					},
					"page_size": &schema.IntAttribute{
						Description: "Page Size",
						Optional:    true,
						Required:    true,
					},
					"insights_lookback_window": &schema.IntAttribute{
						Description: "insights_lookback_window",
						Optional:    true,
						Required:    true,
					},
					"max_batch_size": &schema.IntAttribute{
						Description: "Max Batch Size",
						Optional:    true,
						Required:    true,
					},
					"action_breakdowns_allow_empty": &schema.BoolAttribute{
						Description: "Action Breakdowns Allow Empty",
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
func (r *sourceFacebookMarketingResource) Create(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	// Retrieve values from plan
	var plan sourceFacebookMarketingResourceModel
	err := fwhelpers.UnpackModel(req.PlanContents, &plan)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Generate API request body from plan
	body := api.SourceFacebookMarketing{}
	body.Name = plan.Name
	body.WorkspaceId = plan.WorkspaceId

	body.ConnectionConfiguration = api.SourceFacebookMarketingConnConfig{}
	body.ConnectionConfiguration.SourceType = plan.ConnectionConfiguration.SourceType
	body.ConnectionConfiguration.AccountId = plan.ConnectionConfiguration.AccountId
	body.ConnectionConfiguration.StartDate = plan.ConnectionConfiguration.StartDate
	body.ConnectionConfiguration.AccessToken = plan.ConnectionConfiguration.AccessToken

	body.ConnectionConfiguration.EndDate = plan.ConnectionConfiguration.EndDate
	body.ConnectionConfiguration.IncludeDeleted = plan.ConnectionConfiguration.IncludeDeleted
	body.ConnectionConfiguration.FetchThumbnailImages = plan.ConnectionConfiguration.FetchThumbnailImages
	body.ConnectionConfiguration.PageSize = plan.ConnectionConfiguration.PageSize
	body.ConnectionConfiguration.InsightsLookbackWindow = plan.ConnectionConfiguration.InsightsLookbackWindow
	body.ConnectionConfiguration.MaxBatchSize = plan.ConnectionConfiguration.MaxBatchSize
	body.ConnectionConfiguration.ActionBreakdownsAllowEmpty = plan.ConnectionConfiguration.ActionBreakdownsAllowEmpty

	// Create new source
	source, err := r.Client.CreateFacebookMarketingSource(body)

	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Update resource state with response body
	state := sourceFacebookMarketingResourceModel{}
	state.Name = source.Name
	state.SourceId = source.SourceId
	state.WorkspaceId = source.WorkspaceId

	state.ConnectionConfiguration = sourceFacebookMarketingConnConfigModel{}
	state.ConnectionConfiguration.SourceType = source.ConnectionConfiguration.SourceType
	state.ConnectionConfiguration.AccountId = source.ConnectionConfiguration.AccountId
	state.ConnectionConfiguration.StartDate = source.ConnectionConfiguration.StartDate
	state.ConnectionConfiguration.AccessToken = source.ConnectionConfiguration.AccessToken

	state.ConnectionConfiguration.EndDate = source.ConnectionConfiguration.EndDate
	state.ConnectionConfiguration.IncludeDeleted = source.ConnectionConfiguration.IncludeDeleted
	state.ConnectionConfiguration.FetchThumbnailImages = source.ConnectionConfiguration.FetchThumbnailImages
	state.ConnectionConfiguration.PageSize = source.ConnectionConfiguration.PageSize
	state.ConnectionConfiguration.InsightsLookbackWindow = source.ConnectionConfiguration.InsightsLookbackWindow
	state.ConnectionConfiguration.MaxBatchSize = source.ConnectionConfiguration.MaxBatchSize
	state.ConnectionConfiguration.ActionBreakdownsAllowEmpty = source.ConnectionConfiguration.ActionBreakdownsAllowEmpty

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
func (r *sourceFacebookMarketingResource) Read(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	var state sourceFacebookMarketingResourceModel

	// Get current state
	err := fwhelpers.UnpackModel(req.StateContents, &state)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	res := schema.ServiceResponse{}

	if req.StateID != "" {
		// Query using existing previous state.
		source, err := r.Client.ReadFacebookMarketingSource(req.StateID)
		if err != nil {
			return schema.ErrorResponse(err)
		}

		// Update state with refreshed value
		state.Name = source.Name
		state.SourceId = source.SourceId
		state.WorkspaceId = source.WorkspaceId

		state.ConnectionConfiguration = sourceFacebookMarketingConnConfigModel{}
		state.ConnectionConfiguration.SourceType = source.ConnectionConfiguration.SourceType
		state.ConnectionConfiguration.AccountId = source.ConnectionConfiguration.AccountId
		state.ConnectionConfiguration.StartDate = source.ConnectionConfiguration.StartDate
		state.ConnectionConfiguration.AccessToken = source.ConnectionConfiguration.AccessToken

		state.ConnectionConfiguration.EndDate = source.ConnectionConfiguration.EndDate
		state.ConnectionConfiguration.IncludeDeleted = source.ConnectionConfiguration.IncludeDeleted
		state.ConnectionConfiguration.FetchThumbnailImages = source.ConnectionConfiguration.FetchThumbnailImages
		state.ConnectionConfiguration.PageSize = source.ConnectionConfiguration.PageSize
		state.ConnectionConfiguration.InsightsLookbackWindow = source.ConnectionConfiguration.InsightsLookbackWindow
		state.ConnectionConfiguration.MaxBatchSize = source.ConnectionConfiguration.MaxBatchSize
		state.ConnectionConfiguration.ActionBreakdownsAllowEmpty = source.ConnectionConfiguration.ActionBreakdownsAllowEmpty

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

func (r *sourceFacebookMarketingResource) Update(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	// Retrieve values from plan
	var plan sourceFacebookMarketingResourceModel
	err := fwhelpers.UnpackModel(req.PlanContents, &plan)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Generate API request body from plan
	body := api.SourceFacebookMarketing{}
	body.Name = plan.Name
	body.SourceId = plan.SourceId

	body.ConnectionConfiguration = api.SourceFacebookMarketingConnConfig{}
	body.ConnectionConfiguration.SourceType = plan.ConnectionConfiguration.SourceType
	body.ConnectionConfiguration.AccountId = plan.ConnectionConfiguration.AccountId
	body.ConnectionConfiguration.StartDate = plan.ConnectionConfiguration.StartDate
	body.ConnectionConfiguration.AccessToken = plan.ConnectionConfiguration.AccessToken

	body.ConnectionConfiguration.EndDate = plan.ConnectionConfiguration.EndDate
	body.ConnectionConfiguration.IncludeDeleted = plan.ConnectionConfiguration.IncludeDeleted
	body.ConnectionConfiguration.FetchThumbnailImages = plan.ConnectionConfiguration.FetchThumbnailImages
	body.ConnectionConfiguration.PageSize = plan.ConnectionConfiguration.PageSize
	body.ConnectionConfiguration.InsightsLookbackWindow = plan.ConnectionConfiguration.InsightsLookbackWindow
	body.ConnectionConfiguration.MaxBatchSize = plan.ConnectionConfiguration.MaxBatchSize
	body.ConnectionConfiguration.ActionBreakdownsAllowEmpty = plan.ConnectionConfiguration.ActionBreakdownsAllowEmpty

	// Update existing source
	_, err = r.Client.UpdateFacebookMarketingSource(body)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Fetch updated items
	source, err := r.Client.ReadFacebookMarketingSource(req.PlanID)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Update state with refreshed value
	state := sourceFacebookMarketingResourceModel{}
	state.Name = source.Name
	state.SourceId = source.SourceId
	state.WorkspaceId = source.WorkspaceId

	state.ConnectionConfiguration = sourceFacebookMarketingConnConfigModel{}
	state.ConnectionConfiguration.SourceType = source.ConnectionConfiguration.SourceType
	state.ConnectionConfiguration.AccountId = source.ConnectionConfiguration.AccountId
	state.ConnectionConfiguration.StartDate = source.ConnectionConfiguration.StartDate
	state.ConnectionConfiguration.AccessToken = source.ConnectionConfiguration.AccessToken

	state.ConnectionConfiguration.EndDate = source.ConnectionConfiguration.EndDate
	state.ConnectionConfiguration.IncludeDeleted = source.ConnectionConfiguration.IncludeDeleted
	state.ConnectionConfiguration.FetchThumbnailImages = source.ConnectionConfiguration.FetchThumbnailImages
	state.ConnectionConfiguration.PageSize = source.ConnectionConfiguration.PageSize
	state.ConnectionConfiguration.InsightsLookbackWindow = source.ConnectionConfiguration.InsightsLookbackWindow
	state.ConnectionConfiguration.MaxBatchSize = source.ConnectionConfiguration.MaxBatchSize
	state.ConnectionConfiguration.ActionBreakdownsAllowEmpty = source.ConnectionConfiguration.ActionBreakdownsAllowEmpty

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
func (r *sourceFacebookMarketingResource) Delete(req *schema.ServiceRequest) *schema.ServiceResponse {
	// Delete existing source
	err := r.Client.DeleteFacebookMarketingSource(req.StateID)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	return &schema.ServiceResponse{}
}
