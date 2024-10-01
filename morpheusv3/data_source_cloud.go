package morpheusv3

import (
	"context"
	"fmt"

	"github.com/gomorpheus/morpheus-go-sdk"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

type CloudDataSource struct {
	client *morpheus.Client
}

type CloudDataSourceModel struct {
	ID             types.Int32  `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Code           types.String `tfsdk:"code"`
	Location       types.String `tfsdk:"location"`
	ExternalID     types.String `tfsdk:"external_id"`
	InventoryLevel types.String `tfsdk:"inventory_level"`
	GuidanceMode   types.String `tfsdk:"guidance_mode"`
	TimeZone       types.String `tfsdk:"time_zone"`
	CostingMode    types.String `tfsdk:"costing_mode"`
	Labels         types.Set    `tfsdk:"labels"`
	GroupIDs       types.Set    `tfsdk:"group_ids"`
}

var _ datasource.DataSource = &CloudDataSource{}

func NewCloudDataSource() datasource.DataSource {
	return &CloudDataSource{}
}

func (d *CloudDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud"
}

func (d *CloudDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:        "Provides a Morpheus cloud data source.",
		Blocks:             map[string]schema.Block{},
		DeprecationMessage: "",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				Optional: true,
				Computed: true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the Morpheus cloud",
				Optional:    true,
			},
			"code": schema.StringAttribute{
				Description: "Optional code for use with policies",
				Computed:    true,
			},
			"location": schema.StringAttribute{
				Description: "Optional location for your cloud",
				Computed:    true,
			},
			"external_id": schema.StringAttribute{
				Description: "The external id of the cloud",
				Computed:    true,
			},
			"inventory_level": schema.StringAttribute{
				Description: "The inventory level of the cloud",
				Computed:    true,
			},
			"guidance_mode": schema.StringAttribute{
				Description: "The guidance mode of the cloud",
				Computed:    true,
			},
			"time_zone": schema.StringAttribute{
				Description: "The time zone of the cloud",
				Computed:    true,
			},
			"costing_mode": schema.StringAttribute{
				Description: "The costing mode of the cloud",
				Computed:    true,
			},
			"labels": schema.SetAttribute{
				Description: "The organization labels associated with the cloud",
				Computed:    true,
				ElementType: types.StringType,
			},
			"group_ids": schema.SetAttribute{
				Description: "The ids of the groups granted access to the cloud",
				Computed:    true,
				ElementType: types.Int32Type,
			},
		},
	}
}

func (d *CloudDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*morpheus.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *morpheus.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client

}

//nolint:cyclop,funlen,gocognit,gocyclo // needs refactoring after tests are complete
func (d *CloudDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()
	id := data.ID.ValueInt32()

	// lookup by name if we do not have an id yet
	var morphResponse *morpheus.Response
	var err error
	if id == 0 && name != "" {
		morphResponse, err = d.client.FindCloudByName(name)
	} else if id > 0 {
		morphResponse, err = d.client.GetCloud(int64(id), &morpheus.Request{})
	} else {
		resp.Diagnostics.AddError("Missing data", "Cloud cannot be read without name or id")
		return
	}
	if err != nil {
		if morphResponse != nil {
			resp.Diagnostics.AddError("API Returned Error", fmt.Sprintf("Status code %d, err %v", morphResponse.StatusCode, err))
			return
		} else {
			resp.Diagnostics.AddError("API FAILURE", fmt.Sprintf("%v - %v", morphResponse, err))
			return
		}
	}

	// store resource data
	result := morphResponse.Result.(*morpheus.GetCloudResult)
	cloud := result.Cloud
	if cloud != nil {
		data.ID = types.Int32Value(int32(cloud.ID))
		data.Name = types.StringValue(cloud.Name)
		data.Code = types.StringValue(cloud.Code)
		data.Location = types.StringValue(cloud.Location)
		data.ExternalID = types.StringValue(cloud.ExternalID)
		data.InventoryLevel = types.StringValue(cloud.InventoryLevel)
		data.GuidanceMode = types.StringValue(cloud.GuidanceMode)
		data.TimeZone = types.StringValue(cloud.TimeZone)
		data.CostingMode = types.StringValue(cloud.CostingMode)

		var groupIds []attr.Value
		for _, group := range cloud.Groups {
			groupIds = append(groupIds, types.Int64Value(group.ID))
		}
		groupTypeInt64, diagErr := types.SetValue(types.Int64Type, groupIds)
		data.GroupIDs = groupTypeInt64
		resp.Diagnostics.Append(diagErr...)
		var labels []attr.Value
		for _, label := range cloud.Labels {
			labels = append(labels, types.StringValue(label))
		}
		labelTypeString, diagErr := types.SetValue(types.StringType, labels)
		data.Labels = labelTypeString
		resp.Diagnostics.Append(diagErr...)

		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		resp.Diagnostics.AddError("Not Found", "Cloud not found in response data")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}
