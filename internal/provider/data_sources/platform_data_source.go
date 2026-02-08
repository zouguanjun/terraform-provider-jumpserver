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
	_ datasource.DataSource              = &PlatformDataSource{}
	_ datasource.DataSourceWithConfigure = &PlatformDataSource{}
)

func NewPlatformDataSource() datasource.DataSource {
	return &PlatformDataSource{}
}

type PlatformDataSource struct {
	client *jumpserver.Client
}

type PlatformDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	DisplayName types.String `tfsdk:"display_name"`
	Type        types.String `tfsdk:"type"`
	Category    types.String `tfsdk:"category"`
}

func (d *PlatformDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_platform"
}

func (d *PlatformDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a JumpServer platform",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier of the platform",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The platform name",
			},
			"display_name": schema.StringAttribute{
				Computed:    true,
				Description: "The platform display name",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The platform type",
			},
			"category": schema.StringAttribute{
				Computed:    true,
				Description: "The platform category",
			},
		},
	}
}

func (d *PlatformDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PlatformDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config PlatformDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	platform, err := d.client.GetPlatformByName(config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading platform",
			fmt.Sprintf("Could not read platform %s: %s", config.ID.ValueString(), err),
		)
		return
	}

	config.ID = types.StringValue(platform.GetID())
	config.Name = types.StringValue(platform.Name)
	config.DisplayName = types.StringValue(platform.DisplayName)
	config.Type = types.StringValue(platform.GetTypeValue())
	config.Category = types.StringValue(platform.GetCategoryValue())

	tflog.Trace(ctx, "read platform data source", map[string]any{"id": config.ID.ValueString()})

	diags = resp.State.Set(ctx, config)
	resp.Diagnostics.Append(diags...)
}
