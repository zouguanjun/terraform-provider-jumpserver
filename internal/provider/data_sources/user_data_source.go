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
	_ datasource.DataSource              = &UserDataSource{}
	_ datasource.DataSourceWithConfigure = &UserDataSource{}
)

func NewUserDataSource() datasource.DataSource {
	return &UserDataSource{}
}

type UserDataSource struct {
	client *jumpserver.Client
}

type UserDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	Username types.String `tfsdk:"username"`
	Name     types.String `tfsdk:"name"`
	Email    types.String `tfsdk:"email"`
	IsActive types.Bool   `tfsdk:"is_active"`
	Comment  types.String `tfsdk:"comment"`
}

func (d *UserDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *UserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves information about a JumpServer user",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier of the user",
			},
			"username": schema.StringAttribute{
				Computed:    true,
				Description: "The username",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The display name",
			},
			"email": schema.StringAttribute{
				Computed:    true,
				Description: "The email address",
			},
			"is_active": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the user is active",
			},
			"comment": schema.StringAttribute{
				Computed:    true,
				Description: "Additional comments",
			},
		},
	}
}

func (d *UserDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config UserDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := d.client.GetUser(config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading user",
			fmt.Sprintf("Could not read user ID %s: %s", config.ID.ValueString(), err),
		)
		return
	}

	config.ID = types.StringValue(user.ID)
	config.Username = types.StringValue(user.Username)
	config.Name = types.StringValue(user.Name)
	config.Email = types.StringValue(user.Email)
	config.IsActive = types.BoolValue(user.IsActive)
	config.Comment = types.StringValue(user.Comment)

	tflog.Trace(ctx, "read user data source", map[string]any{"id": config.ID.ValueString()})

	diags = resp.State.Set(ctx, config)
	resp.Diagnostics.Append(diags...)
}
