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
	_ datasource.DataSource              = &TaskDataSource{}
	_ datasource.DataSourceWithConfigure = &TaskDataSource{}
)

func NewTaskDataSource() datasource.DataSource {
	return &TaskDataSource{}
}

type TaskDataSource struct {
	client *jumpserver.Client
}

type TaskDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	Output   types.String `tfsdk:"output"`
	Finished types.Bool   `tfsdk:"finished"`
	Mark     types.String `tfsdk:"mark"`
}

func (d *TaskDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_task"
}

func (d *TaskDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves task execution information from JumpServer",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "The unique identifier of the task execution",
			},
			"output": schema.StringAttribute{
				Computed:    true,
				Description: "The output data from the command execution",
			},
			"finished": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the task execution is finished",
			},
			"mark": schema.StringAttribute{
				Computed:    true,
				Description: "The mark identifier for the task execution",
			},
		},
	}
}

func (d *TaskDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *TaskDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config TaskDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	task, err := d.client.GetCommandExecution(config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading task execution",
			fmt.Sprintf("Could not read task ID %s: %s", config.ID.ValueString(), err),
		)
		return
	}

	config.ID = types.StringValue(config.ID.ValueString())
	config.Output = types.StringValue(task.Data)
	config.Finished = types.BoolValue(task.End)
	config.Mark = types.StringValue(task.Mark)

	tflog.Trace(ctx, "read task data source", map[string]any{"id": config.ID.ValueString()})

	diags = resp.State.Set(ctx, config)
	resp.Diagnostics.Append(diags...)
}
