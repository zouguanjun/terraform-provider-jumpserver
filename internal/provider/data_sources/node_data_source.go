package data_sources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"jumpserver/internal/jumpserver"
)

var (
	_ datasource.DataSource              = &NodeDataSource{}
	_ datasource.DataSourceWithConfigure = &NodeDataSource{}
)

func NewNodeDataSource() datasource.DataSource {
	return &NodeDataSource{}
}

type NodeDataSource struct {
	client *jumpserver.Client
}

type NodeDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	FullName types.String `tfsdk:"full_name"`
	Weight   types.Int64  `tfsdk:"weight"`
}

func (d *NodeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_node"
}

func (d *NodeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a JumpServer organization node",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier or full name of the node",
			},
			"full_name": schema.StringAttribute{
				Computed:    true,
				Description: "The full path name of the node",
			},
			"weight": schema.Int64Attribute{
				Computed:    true,
				Description: "The weight of the node",
			},
		},
	}
}

func (d *NodeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*jumpserver.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *jumpserver.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *NodeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config NodeDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	node, err := d.client.GetNodeByFullName(config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading node",
			fmt.Sprintf("Could not read node %s: %s", config.ID.ValueString(), err),
		)
		return
	}

	config.ID = types.StringValue(node.ID)
	config.FullName = types.StringValue(node.FullName)
	config.Weight = types.Int64Value(int64(node.Weight))

	tflog.Trace(ctx, "read node data source", map[string]any{"id": config.ID.ValueString()})

	diags = resp.State.Set(ctx, config)
	resp.Diagnostics.Append(diags...)
}
