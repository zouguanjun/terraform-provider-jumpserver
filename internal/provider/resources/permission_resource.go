package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"jumpserver/internal/jumpserver"
)

var (
	_ resource.Resource                = &PermissionResource{}
	_ resource.ResourceWithConfigure   = &PermissionResource{}
	_ resource.ResourceWithImportState = &PermissionResource{}
)

func NewPermissionResource() resource.Resource {
	return &PermissionResource{}
}

type PermissionResource struct {
	client *jumpserver.Client
}

type PermissionResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Users       types.Set    `tfsdk:"users"`
	UserGroups  types.Set    `tfsdk:"user_groups"`
	Assets      types.Set    `tfsdk:"assets"`
	AssetGroups types.Set    `tfsdk:"asset_groups"`
	Actions     types.Set    `tfsdk:"actions"`
	Comment     types.String `tfsdk:"comment"`
}

func (r *PermissionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_permission"
}

func (r *PermissionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages permission rules in JumpServer",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The unique identifier of the permission",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the permission",
			},
			"users": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of user IDs to grant permission",
			},
			"user_groups": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of user group IDs to grant permission",
			},
			"assets": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of asset IDs to grant access to",
			},
			"asset_groups": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of asset group IDs to grant access to",
			},
			"actions": schema.SetAttribute{
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(stringvalidator.OneOf("connect", "upload", "download", "clipboard", "delete_file")),
				},
				Description: "Allowed actions (e.g., 'connect', 'upload', 'download')",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Description: "Additional comments about the permission",
			},
		},
	}
}

func (r *PermissionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*jumpserver.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *jumpserver.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = client
}

func toStringSet(set basetypes.SetValue) []string {
	var result []string
	if !set.IsUnknown() && !set.IsNull() {
		set.ElementsAs(context.Background(), &result, false)
	}
	return result
}

func (r *PermissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PermissionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &jumpserver.CreatePermissionRequest{
		Name:        plan.Name.ValueString(),
		Users:       toStringSet(plan.Users),
		UserGroups:  toStringSet(plan.UserGroups),
		Assets:      toStringSet(plan.Assets),
		AssetGroups: toStringSet(plan.AssetGroups),
		Actions:     toStringSet(plan.Actions),
		Comment:     plan.Comment.ValueString(),
	}

	permission, err := r.client.CreatePermission(createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating permission",
			fmt.Sprintf("Could not create permission: %s", err),
		)
		return
	}

	plan.ID = types.StringValue(permission.ID)
	plan.Name = types.StringValue(permission.Name)
	plan.Comment = types.StringValue(permission.Comment)
	plan.Users, _ = types.SetValueFrom(ctx, types.StringType, permission.GetUserIDs())
	plan.UserGroups, _ = types.SetValueFrom(ctx, types.StringType, permission.UserGroups)
	plan.Assets, _ = types.SetValueFrom(ctx, types.StringType, permission.GetAssetIDs())
	plan.AssetGroups, _ = types.SetValueFrom(ctx, types.StringType, permission.AssetGroups)
	plan.Actions, _ = types.SetValueFrom(ctx, types.StringType, permission.GetActionValues())

	tflog.Trace(ctx, "created permission", map[string]any{"id": permission.ID})

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *PermissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PermissionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	permission, err := r.client.GetPermission(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading permission",
			fmt.Sprintf("Could not read permission: %s", err),
		)
		return
	}

	state.Name = types.StringValue(permission.Name)
	state.Comment = types.StringValue(permission.Comment)
	state.Users, _ = types.SetValueFrom(ctx, types.StringType, permission.GetUserIDs())
	state.UserGroups, _ = types.SetValueFrom(ctx, types.StringType, permission.UserGroups)
	state.Assets, _ = types.SetValueFrom(ctx, types.StringType, permission.GetAssetIDs())
	state.AssetGroups, _ = types.SetValueFrom(ctx, types.StringType, permission.AssetGroups)
	state.Actions, _ = types.SetValueFrom(ctx, types.StringType, permission.GetActionValues())

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *PermissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PermissionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &jumpserver.UpdatePermissionRequest{
		Name:        plan.Name.ValueString(),
		Users:       toStringSet(plan.Users),
		UserGroups:  toStringSet(plan.UserGroups),
		Assets:      toStringSet(plan.Assets),
		AssetGroups: toStringSet(plan.AssetGroups),
		Actions:     toStringSet(plan.Actions),
		Comment:     plan.Comment.ValueString(),
	}

	permission, err := r.client.UpdatePermission(plan.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating permission",
			fmt.Sprintf("Could not update permission: %s", err),
		)
		return
	}

	plan.Name = types.StringValue(permission.Name)
	plan.Comment = types.StringValue(permission.Comment)
	plan.Users, _ = types.SetValueFrom(ctx, types.StringType, permission.GetUserIDs())
	plan.UserGroups, _ = types.SetValueFrom(ctx, types.StringType, permission.UserGroups)
	plan.Assets, _ = types.SetValueFrom(ctx, types.StringType, permission.GetAssetIDs())
	plan.AssetGroups, _ = types.SetValueFrom(ctx, types.StringType, permission.AssetGroups)
	plan.Actions, _ = types.SetValueFrom(ctx, types.StringType, permission.GetActionValues())

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *PermissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state PermissionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeletePermission(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting permission",
			fmt.Sprintf("Could not delete permission: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "deleted permission", map[string]any{"id": state.ID.ValueString()})
}

func (r *PermissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
