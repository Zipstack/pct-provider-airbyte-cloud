package plugin

import (
	"fmt"
	"time"

	"github.com/zipstack/pct-plugin-framework/fwhelpers"
	"github.com/zipstack/pct-plugin-framework/schema"

	"github.com/zipstack/pct-provider-airbyte-cloud/api"
)

type connectionResource struct {
	Client *api.Client
}

type connectionResourceModel struct {
	Name                             string           `pctsdk:"name"`
	SourceID                         string           `pctsdk:"source_id"`
	DestinationID                    string           `pctsdk:"destination_id"`
	ConnectionID                     string           `pctsdk:"connection_id"`
	DataResidency                    string           `pctsdk:"data_residency"`
	NamespaceDefinition              string           `pctsdk:"namespace_definition"`
	NamespaceFormat                  string           `pctsdk:"namespace_format"`
	NonBreakingSchemaUpdatesBehavior string           `pctsdk:"nonBreakingSchemaUpdatesBehavior"`
	Prefix                           string           `pctsdk:"prefix"`
	Status                           string           `pctsdk:"status"`
	Schedule                         connScheduleData `pctsdk:"schedule"`
	// OperatorConfiguration connOperatorConfig `pctsdk:"operator_configuration"`
}

type connScheduleData struct {
	ScheduleType   string `pctsdk:"schedule_type"`
	CronExpression string `pctsdk:"cron_expression"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ schema.ResourceService = &connectionResource{}
)

// Helper function to return a resource service instance.
func NewConnectionResource() schema.ResourceService {
	return &connectionResource{}
}

// Metadata returns the resource type name.
// It is always provider name + "_" + resource type name.
func (r *connectionResource) Metadata(req *schema.ServiceRequest) *schema.ServiceResponse {
	return &schema.ServiceResponse{
		TypeName: req.TypeName + "_connection",
	}
}

// Configure adds the provider configured client to the resource.
func (r *connectionResource) Configure(req *schema.ServiceRequest) *schema.ServiceResponse {
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
func (r *connectionResource) Schema() *schema.ServiceResponse {
	s := &schema.Schema{
		Description: "Connection resource for Airbyte",
		Attributes: map[string]schema.Attribute{
			"name": &schema.StringAttribute{
				Description: "Name",
				Required:    true,
			},
			"source_id": &schema.StringAttribute{
				Description: "Source ID",
				Required:    true,
			},
			"destination_id": &schema.StringAttribute{
				Description: "Destination ID",
				Required:    true,
			},
			"connection_id": &schema.StringAttribute{
				Description: "Connection ID",
				Computed:    true,
			},
			"status": &schema.StringAttribute{
				Description: "Status",
				Required:    true,
			},
			"schedule": &schema.MapAttribute{
				Description: "Schedule",
				Required:    true,
				Attributes: map[string]schema.Attribute{
					"schedule_type": &schema.StringAttribute{
						Description: "schedule Type",
						Required:    true,
					},
					"cron_expression": &schema.StringAttribute{
						Description: "cron Expression",
						Required:    true,
						Optional:    true,
					},
				},
			},
			"data_residency": &schema.StringAttribute{
				Description: "Data Residency",
				Required:    true,
				Optional:    true,
			},
			"namespace_definition": &schema.StringAttribute{
				Description: "namespace definition",
				Required:    true,
				Optional:    true,
			},
			"namespace_format": &schema.StringAttribute{
				Description: "namespace Format",
				Required:    true,
				Optional:    true,
			},
			"prefix": &schema.StringAttribute{
				Description: "prefix",
				Required:    true,
				Optional:    true,
			},
			"nonBreakingSchemaUpdatesBehavior": &schema.StringAttribute{
				Description: "nonBreakingSchemaUpdatesBehavior",
				Required:    true,
				Optional:    true,
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
func (r *connectionResource) Create(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	var plan connectionResourceModel
	err := fwhelpers.UnpackModel(req.PlanContents, &plan)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	body := api.ConnectionResource{}
	body.Name = plan.Name
	body.SourceID = plan.SourceID
	body.DestinationID = plan.DestinationID

	body.DataResidency = plan.DataResidency
	body.Status = plan.Status
	body.NamespaceDefinition = plan.NamespaceDefinition
	body.NamespaceFormat = plan.NamespaceFormat
	body.NonBreakingSchemaUpdatesBehavior = plan.NonBreakingSchemaUpdatesBehavior
	body.Prefix = plan.Prefix

	body.Schedule = api.ConnScheduleData{}
	body.Schedule.ScheduleType = plan.Schedule.ScheduleType
	body.Schedule.CronExpression = plan.Schedule.CronExpression

	connection, err := r.Client.CreateConnectionResource(body)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Map response body to schema and populate Computed attribute values
	state := connectionResourceModel{}

	state.Name = connection.Name
	state.ConnectionID = connection.ConnectionID
	state.SourceID = connection.SourceID
	state.DestinationID = connection.DestinationID

	state.DataResidency = connection.DataResidency
	state.Status = connection.Status
	state.NamespaceDefinition = connection.NamespaceDefinition
	state.NamespaceFormat = connection.NamespaceFormat
	state.NonBreakingSchemaUpdatesBehavior = connection.NonBreakingSchemaUpdatesBehavior
	state.Prefix = connection.Prefix

	state.Schedule = connScheduleData{}
	state.Schedule.ScheduleType = connection.Schedule.ScheduleType
	state.Schedule.CronExpression = connection.Schedule.CronExpression

	stateEnc, err := fwhelpers.PackModel(nil, &state)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	return &schema.ServiceResponse{
		StateID:          state.ConnectionID,
		StateContents:    stateEnc,
		StateLastUpdated: time.Now().Format(time.RFC850),
	}
}

// Read resource information
func (r *connectionResource) Read(req *schema.ServiceRequest) *schema.ServiceResponse {
	var state connectionResourceModel

	// Get current state
	err := fwhelpers.UnpackModel(req.StateContents, &state)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	res := schema.ServiceResponse{}

	if req.StateID != "" {
		// Query using existing previous state.
		connection, err := r.Client.ReadConnectionResource(req.StateID)
		if err != nil {
			return schema.ErrorResponse(err)
		}

		state = connectionResourceModel{}

		// Update state with refreshed value
		state.Name = connection.Name
		state.ConnectionID = connection.ConnectionID
		state.SourceID = connection.SourceID
		state.DestinationID = connection.DestinationID

		state.DataResidency = connection.DataResidency
		state.Status = connection.Status
		state.NamespaceDefinition = connection.NamespaceDefinition
		state.NamespaceFormat = connection.NamespaceFormat
		state.NonBreakingSchemaUpdatesBehavior = connection.NonBreakingSchemaUpdatesBehavior
		state.Prefix = connection.Prefix

		state.Schedule = connScheduleData{}
		state.Schedule.ScheduleType = connection.Schedule.ScheduleType
		state.Schedule.CronExpression = connection.Schedule.CronExpression

		res.StateID = connection.ConnectionID
	} else {
		// No previous state exists.
		res.StateID = ""
		res.StateLastUpdated = ""
	}

	// Set refreshed state
	stateEnc, err := fwhelpers.PackModel(nil, &state)
	if err != nil {
		return schema.ErrorResponse(err)
	}
	res.StateContents = stateEnc

	return &res
}

func (r *connectionResource) Update(req *schema.ServiceRequest) *schema.ServiceResponse {
	var plan connectionResourceModel
	err := fwhelpers.UnpackModel(req.PlanContents, &plan)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Generate API request body from plan
	body := api.ConnectionResource{}

	body.Name = plan.Name
	body.ConnectionID = plan.ConnectionID
	body.SourceID = plan.SourceID
	body.DestinationID = plan.DestinationID

	body.DataResidency = plan.DataResidency
	body.Status = plan.Status
	body.NamespaceDefinition = plan.NamespaceDefinition
	body.NamespaceFormat = plan.NamespaceFormat
	body.NonBreakingSchemaUpdatesBehavior = plan.NonBreakingSchemaUpdatesBehavior
	body.Prefix = plan.Prefix

	body.Schedule = api.ConnScheduleData{}
	body.Schedule.ScheduleType = plan.Schedule.ScheduleType
	body.Schedule.CronExpression = plan.Schedule.CronExpression

	// Update existing source
	_, err = r.Client.UpdateConnectionResource(body)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Fetch updated items
	connection, err := r.Client.ReadConnectionResource(req.PlanID)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Update state with refreshed value
	state := connectionResourceModel{}

	state.Name = connection.Name
	state.ConnectionID = connection.ConnectionID
	state.SourceID = connection.SourceID
	state.DestinationID = connection.DestinationID

	state.DataResidency = connection.DataResidency
	state.Status = connection.Status
	state.NamespaceDefinition = connection.NamespaceDefinition
	state.NamespaceFormat = connection.NamespaceFormat
	state.NonBreakingSchemaUpdatesBehavior = connection.NonBreakingSchemaUpdatesBehavior
	state.Prefix = connection.Prefix

	state.Schedule = connScheduleData{}
	state.Schedule.ScheduleType = connection.Schedule.ScheduleType
	state.Schedule.CronExpression = connection.Schedule.CronExpression

	// Set refreshed state
	stateEnc, err := fwhelpers.PackModel(nil, &state)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	return &schema.ServiceResponse{
		StateID:          state.ConnectionID,
		StateContents:    stateEnc,
		StateLastUpdated: time.Now().Format(time.RFC850),
	}
}

// Delete deletes the resource and removes the state on success.
func (r *connectionResource) Delete(req *schema.ServiceRequest) *schema.ServiceResponse {
	// Delete existing source
	err := r.Client.DeleteConnectionResource(req.StateID)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	return &schema.ServiceResponse{}
}
