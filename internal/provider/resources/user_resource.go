package resources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"jumpserver/internal/jumpserver"
)

var (
	_ resource.Resource                = &UserResource{}
	_ resource.ResourceWithConfigure   = &UserResource{}
	_ resource.ResourceWithImportState = &UserResource{}
)

func NewUserResource() resource.Resource {
	return &UserResource{}
}

type UserResource struct {
	client *jumpserver.Client
}

type UserResourceModel struct {
	ID       types.String `tfsdk:"id"`
	Username types.String `tfsdk:"username"`
	Name     types.String `tfsdk:"name"`
	Email    types.String `tfsdk:"email"`
	IsActive types.Bool   `tfsdk:"is_active"`
	Comment  types.String `tfsdk:"comment"`
}

func (r *UserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *UserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a JumpServer user",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The unique identifier of the user",
			},
			"username": schema.StringAttribute{
				Required:    true,
				Description: "The username for login",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The display name of the user",
			},
			"email": schema.StringAttribute{
				Required:    true,
				Description: "The email address of the user",
			},
			"is_active": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the user is active",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Description: "Additional comments about the user",
			},
		},
	}
}

func (r *UserResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &jumpserver.CreateUserRequest{
		Username: plan.Username.ValueString(),
		Name:     plan.Name.ValueString(),
		Email:    plan.Email.ValueString(),
		IsActive: plan.IsActive.ValueBool(),
		Comment:  plan.Comment.ValueString(),
	}

	user, err := r.client.CreateUser(createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user",
			fmt.Sprintf("Could not create user: %s", err),
		)
		return
	}

	plan.ID = types.StringValue(user.ID)
	plan.Username = types.StringValue(user.Username)
	plan.Name = types.StringValue(user.Name)
	plan.Email = types.StringValue(user.Email)
	plan.IsActive = types.BoolValue(user.IsActive)
	plan.Comment = types.StringValue(user.Comment)

	tflog.Trace(ctx, "created user", map[string]any{"id": plan.ID.ValueString()})

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.GetUser(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading user",
			fmt.Sprintf("Could not read user: %s", err),
		)
		return
	}

	state.ID = types.StringValue(user.ID)
	state.Username = types.StringValue(user.Username)
	state.Name = types.StringValue(user.Name)
	state.Email = types.StringValue(user.Email)
	state.IsActive = types.BoolValue(user.IsActive)
	state.Comment = types.StringValue(user.Comment)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan UserResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &jumpserver.UpdateUserRequest{
		Username: plan.Username.ValueString(),
		Name:     plan.Name.ValueString(),
		Email:    plan.Email.ValueString(),
		IsActive: &[]bool{plan.IsActive.ValueBool()}[0],
		Comment:  plan.Comment.ValueString(),
	}

	user, err := r.client.UpdateUser(plan.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating user",
			fmt.Sprintf("Could not update user: %s", err),
		)
		return
	}

	plan.Name = types.StringValue(user.Name)
	plan.Email = types.StringValue(user.Email)
	plan.IsActive = types.BoolValue(user.IsActive)
	plan.Comment = types.StringValue(user.Comment)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUser(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting user",
			fmt.Sprintf("Could not delete user: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "deleted user", map[string]any{"id": state.ID.ValueString()})
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
