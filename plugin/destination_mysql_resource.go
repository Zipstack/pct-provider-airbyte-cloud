package plugin

import (
	"fmt"
	"time"

	"github.com/zipstack/pct-plugin-framework/fwhelpers"
	"github.com/zipstack/pct-plugin-framework/schema"

	"github.com/zipstack/pct-provider-airbyte-cloud/api"
)

// Resource implementation.
type destinationMysqlResource struct {
	Client *api.Client
}

type destinationMysqlResourceModel struct {
	Name                    string                          `pctsdk:"name"`
	DestinationId           string                          `pctsdk:"destination_id"`
	WorkspaceId             string                          `pctsdk:"workspace_id"`
	ConnectionConfiguration destinationMysqlConnConfigModel `pctsdk:"configuration"`
}

type destinationMysqlConnConfigModel struct {
	DestinationType string `pctsdk:"destination_type"`
	Host            string `pctsdk:"host"`
	Username        string `pctsdk:"username"`
	Password        string `pctsdk:"password"`
	Database        string `pctsdk:"database"`
	Port            int    `pctsdk:"port"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ schema.ResourceService = &destinationMysqlResource{}
)

// Helper function to return a resource service instance.
func NewDestinationMysqlResource() schema.ResourceService {
	return &destinationMysqlResource{}
}

// Metadata returns the resource type name.
// It is always provider name + "_" + resource type name.
func (r *destinationMysqlResource) Metadata(req *schema.ServiceRequest) *schema.ServiceResponse {
	return &schema.ServiceResponse{
		TypeName: req.TypeName + "_destination_mysql",
	}
}

// Configure adds the provider configured client to the resource.
func (r *destinationMysqlResource) Configure(req *schema.ServiceRequest) *schema.ServiceResponse {
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
func (r *destinationMysqlResource) Schema() *schema.ServiceResponse {
	s := &schema.Schema{
		Description: "Destination mysql resource for Airbyte",
		Attributes: map[string]schema.Attribute{
			"name": &schema.StringAttribute{
				Description: "Name",
				Required:    true,
			},
			"destination_id": &schema.StringAttribute{
				Description: "Destination ID",
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
					"destination_type": &schema.StringAttribute{
						Description: "Destination Type",
						Required:    true,
					},
					"port": &schema.IntAttribute{
						Description: "Port",
						Required:    true,
					},
					"host": &schema.StringAttribute{
						Description: "Host",
						Required:    true,
					},
					"username": &schema.StringAttribute{
						Description: "Username",
						Required:    true,
					},
					"password": &schema.StringAttribute{
						Description: "Password",
						Required:    true,
						Sensitive:   true,
					},
					"database": &schema.StringAttribute{
						Description: "Database",
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
func (r *destinationMysqlResource) Create(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	// Retrieve values from plan
	var plan destinationMysqlResourceModel
	err := fwhelpers.UnpackModel(req.PlanContents, &plan)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Generate API request body from plan
	body := api.DestinationMysql{}
	body.Name = plan.Name
	body.WorkspaceId = plan.WorkspaceId

	body.ConnectionConfiguration = api.DestinationMysqlConnConfig{}
	body.ConnectionConfiguration.DestinationType = plan.ConnectionConfiguration.DestinationType
	body.ConnectionConfiguration.Port = plan.ConnectionConfiguration.Port
	body.ConnectionConfiguration.Username = plan.ConnectionConfiguration.Username
	body.ConnectionConfiguration.Password = plan.ConnectionConfiguration.Password
	body.ConnectionConfiguration.Host = plan.ConnectionConfiguration.Host
	body.ConnectionConfiguration.Database = plan.ConnectionConfiguration.Database

	// Create new destination
	destination, err := r.Client.CreateMysqlDestination(body)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Update resource state with response body
	state := destinationMysqlResourceModel{}
	state.Name = destination.Name
	state.DestinationId = destination.DestinationId
	state.WorkspaceId = destination.WorkspaceId

	state.ConnectionConfiguration = destinationMysqlConnConfigModel{}
	state.ConnectionConfiguration.Host = plan.ConnectionConfiguration.Host
	state.ConnectionConfiguration.Port = plan.ConnectionConfiguration.Port
	state.ConnectionConfiguration.Username = plan.ConnectionConfiguration.Username
	state.ConnectionConfiguration.Password = plan.ConnectionConfiguration.Password
	state.ConnectionConfiguration.DestinationType = plan.ConnectionConfiguration.DestinationType
	state.ConnectionConfiguration.Database = plan.ConnectionConfiguration.Database

	// Set refreshed state
	stateEnc, err := fwhelpers.PackModel(nil, &state)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	return &schema.ServiceResponse{
		StateID:          state.DestinationId,
		StateContents:    stateEnc,
		StateLastUpdated: time.Now().Format(time.RFC850),
	}
}

// Read resource information
func (r *destinationMysqlResource) Read(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	var state destinationMysqlResourceModel

	// Get current state
	err := fwhelpers.UnpackModel(req.StateContents, &state)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	res := schema.ServiceResponse{}

	if req.StateID != "" {
		// Query using existing previous state.
		destination, err := r.Client.ReadMysqlDestination(req.StateID)
		if err != nil {
			return schema.ErrorResponse(err)
		}

		// Update state with refreshed value
		state.Name = destination.Name
		state.DestinationId = destination.DestinationId
		state.WorkspaceId = destination.WorkspaceId

		res.StateID = state.DestinationId
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

func (r *destinationMysqlResource) Update(req *schema.ServiceRequest) *schema.ServiceResponse {
	// logger := fwhelpers.GetLogger()

	// Retrieve values from plan
	var plan destinationMysqlResourceModel
	err := fwhelpers.UnpackModel(req.PlanContents, &plan)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Generate API request body from plan
	body := api.DestinationMysql{}
	body.Name = plan.Name
	body.DestinationId = plan.DestinationId

	body.ConnectionConfiguration = api.DestinationMysqlConnConfig{}
	body.ConnectionConfiguration.DestinationType = plan.ConnectionConfiguration.DestinationType
	body.ConnectionConfiguration.Port = plan.ConnectionConfiguration.Port
	body.ConnectionConfiguration.Host = plan.ConnectionConfiguration.Host
	body.ConnectionConfiguration.Username = plan.ConnectionConfiguration.Username
	body.ConnectionConfiguration.Password = plan.ConnectionConfiguration.Password
	body.ConnectionConfiguration.Database = plan.ConnectionConfiguration.Database

	// Update existing destination
	_, err = r.Client.UpdateMysqlDestination(body)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Fetch updated items
	destination, err := r.Client.ReadMysqlDestination(req.PlanID)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	// Update state with refreshed value
	state := destinationMysqlResourceModel{}
	state.Name = destination.Name
	state.DestinationId = destination.DestinationId
	state.WorkspaceId = destination.WorkspaceId

	state.ConnectionConfiguration = destinationMysqlConnConfigModel{}
	state.ConnectionConfiguration.DestinationType = plan.ConnectionConfiguration.DestinationType
	state.ConnectionConfiguration.Port = plan.ConnectionConfiguration.Port
	state.ConnectionConfiguration.Host = plan.ConnectionConfiguration.Host
	state.ConnectionConfiguration.Username = plan.ConnectionConfiguration.Username
	state.ConnectionConfiguration.Password = plan.ConnectionConfiguration.Password
	state.ConnectionConfiguration.Database = plan.ConnectionConfiguration.Database

	// Set refreshed state
	stateEnc, err := fwhelpers.PackModel(nil, &state)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	return &schema.ServiceResponse{
		StateID:          state.DestinationId,
		StateContents:    stateEnc,
		StateLastUpdated: time.Now().Format(time.RFC850),
	}
}

// Delete deletes the resource and removes the state on success.
func (r *destinationMysqlResource) Delete(req *schema.ServiceRequest) *schema.ServiceResponse {
	// Delete existing destination
	err := r.Client.DeleteMysqlDestination(req.StateID)
	if err != nil {
		return schema.ErrorResponse(err)
	}

	return &schema.ServiceResponse{}
}
