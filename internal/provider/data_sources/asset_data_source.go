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
	_ datasource.DataSource              = &AssetDataSource{}
	_ datasource.DataSourceWithConfigure = &AssetDataSource{}
)

func NewAssetDataSource() datasource.DataSource {
	return &AssetDataSource{}
}

type AssetDataSource struct {
	client *jumpserver.Client
}

type AssetDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Address  types.String `tfsdk:"address"`
	Platform types.String `tfsdk:"platform"`
	Gateway  types.String `tfsdk:"gateway"`
	Nodes    types.List   `tfsdk:"nodes"`
	IsActive types.Bool   `tfsdk:"is_active"`
	Comment  types.String `tfsdk:"comment"`
}

func (d *AssetDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_asset"
}

func (d *AssetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a JumpServer asset",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier of the asset",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the asset",
			},
			"address": schema.StringAttribute{
				Computed:    true,
				Description: "The IP address or hostname of the asset",
			},
			"platform": schema.StringAttribute{
				Computed:    true,
				Description: "The platform name",
			},
			"gateway": schema.StringAttribute{
				Computed:    true,
				Description: "The gateway asset ID",
			},
			"nodes": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "List of organization node IDs",
			},
			"is_active": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the asset is active",
			},
			"comment": schema.StringAttribute{
				Computed:    true,
				Description: "Additional comments",
			},
		},
	}
}

func (d *AssetDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AssetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config AssetDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	asset, err := d.client.GetAsset(config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading asset",
			fmt.Sprintf("Could not read asset ID %s: %s", config.ID.ValueString(), err),
		)
		return
	}

	config.ID = types.StringValue(asset.ID)
	config.Name = types.StringValue(asset.Name)
	config.Address = types.StringValue(asset.Address)
	config.Platform = types.StringValue(asset.Platform.Name)

	var nodeIDs []any
	for _, node := range asset.Nodes {
		nodeIDs = append(nodeIDs, node.ID)
	}
	diags = config.Nodes.ElementsAs(ctx, nodeIDs, false)
	resp.Diagnostics.Append(diags...)

	config.IsActive = types.BoolValue(asset.IsActive)
	config.Comment = types.StringValue(asset.Comment)

	tflog.Trace(ctx, "read asset data source", map[string]any{"id": config.ID.ValueString()})

	diags = resp.State.Set(ctx, config)
	resp.Diagnostics.Append(diags...)
}
