package resources

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	_ resource.Resource                = &AccountResource{}
	_ resource.ResourceWithConfigure   = &AccountResource{}
	_ resource.ResourceWithImportState = &AccountResource{}
)

func NewAccountResource() resource.Resource {
	return &AccountResource{}
}

type AccountResource struct {
	client *jumpserver.Client
}

type AccountResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"username"`
	Asset      types.String `tfsdk:"asset"`
	Secret     types.String `tfsdk:"secret"`
	SecretType types.String `tfsdk:"secret_type"`
	Comment    types.String `tfsdk:"comment"`
}

func (r *AccountResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account"
}

func (r *AccountResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an account on a JumpServer asset",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The unique identifier of the account",
			},
			"username": schema.StringAttribute{
				Required:    true,
				Description: "The username",
			},
			"asset": schema.StringAttribute{
				Required:    true,
				Description: "The asset ID to associate the account with",
			},
			"secret": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The password or secret for the account",
			},
			"secret_type": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The secret type (e.g., 'password', 'ssh_key', 'access_key')",
				Validators: []validator.String{
					stringvalidator.OneOf("password", "ssh_key", "access_key", "token"),
				},
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Description: "Additional comments about the account",
			},
		},
	}
}

func (r *AccountResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *AccountResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AccountResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &jumpserver.CreateAccountRequest{
		Name:       plan.Name.ValueString(),
		Asset:      plan.Asset.ValueString(),
		Secret:     plan.Secret.ValueString(),
		SecretType: plan.SecretType.ValueString(),
		Comment:    plan.Comment.ValueString(),
	}

	account, err := r.client.CreateAccount(createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating account",
			fmt.Sprintf("Could not create account: %s", err),
		)
		return
	}

	plan.ID = types.StringValue(account.ID)
	plan.Name = types.StringValue(account.Name)
	plan.Asset = types.StringValue(account.GetAssetID())
	plan.SecretType = types.StringValue(account.GetSecretTypeValue())
	plan.Comment = types.StringValue(account.Comment)

	tflog.Trace(ctx, "created account", map[string]any{"id": plan.ID.ValueString()})

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *AccountResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AccountResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	account, err := r.client.GetAccount(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading account",
			fmt.Sprintf("Could not read account: %s", err),
		)
		return
	}

	state.ID = types.StringValue(account.ID)
	state.Name = types.StringValue(account.Name)
	state.Asset = types.StringValue(account.GetAssetID())
	state.SecretType = types.StringValue(account.GetSecretTypeValue())
	state.Comment = types.StringValue(account.Comment)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *AccountResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AccountResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &jumpserver.UpdateAccountRequest{
		Name:       plan.Name.ValueString(),
		Secret:     plan.Secret.ValueString(),
		SecretType: plan.SecretType.ValueString(),
		Comment:    plan.Comment.ValueString(),
	}

	account, err := r.client.UpdateAccount(plan.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating account",
			fmt.Sprintf("Could not update account: %s", err),
		)
		return
	}

	plan.Name = types.StringValue(account.Name)
	plan.SecretType = types.StringValue(account.GetSecretTypeValue())
	plan.Comment = types.StringValue(account.Comment)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *AccountResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AccountResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAccount(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting account",
			fmt.Sprintf("Could not delete account: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "deleted account", map[string]any{"id": state.ID.ValueString()})
}

func (r *AccountResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
