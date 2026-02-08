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

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &AssetResource{}
	_ resource.ResourceWithConfigure   = &AssetResource{}
	_ resource.ResourceWithImportState = &AssetResource{}
)

// NewAssetResource is a helper function to simplify the provider implementation.
func NewAssetResource() resource.Resource {
	return &AssetResource{}
}

// AssetResource defines the resource implementation.
type AssetResource struct {
	client *jumpserver.Client
}

// AssetResourceModel describes the resource data model.
type AssetResourceModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Address  types.String `tfsdk:"address"`
	Platform types.String `tfsdk:"platform"`
	Nodes    types.List   `tfsdk:"nodes"`
	IsActive types.Bool   `tfsdk:"is_active"`
	Comment  types.String `tfsdk:"comment"`
}

func (r *AssetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_asset"
}

func (r *AssetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a JumpServer asset (server, device, etc.)",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The unique identifier of the asset",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the asset",
			},
			"address": schema.StringAttribute{
				Required:    true,
				Description: "The IP address or hostname of the asset",
			},
			"platform": schema.StringAttribute{
				Required:    true,
				Description: "The platform ID (e.g., '1' for Linux)",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},

			"nodes": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "List of organization node IDs to associate the asset with",
			},
			"is_active": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether the asset is active",
			},
			"comment": schema.StringAttribute{
				Optional:    true,
				Description: "Additional comments about the asset",
			},
		},
	}
}

func (r *AssetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*jumpserver.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *jumpserver.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *AssetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AssetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert nodes list
	var nodes []string
	diags = plan.Nodes.ElementsAs(ctx, &nodes, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create asset - parse platform as integer
	platformID := 1 // default value
	if plan.Platform.ValueString() != "" {
		fmt.Sscanf(plan.Platform.ValueString(), "%d", &platformID)
	}

	// Convert nodes to NodeRequest format
	var nodeReqs []jumpserver.NodeRequest
	for _, nodeID := range nodes {
		nodeReqs = append(nodeReqs, jumpserver.NodeRequest{PK: nodeID})
	}

	createReq := &jumpserver.CreateAssetRequest{
		Name:     plan.Name.ValueString(),
		Address:  plan.Address.ValueString(),
		Platform: jumpserver.PlatformRequest{PK: platformID},
		Nodes:    nodeReqs,
		IsActive: plan.IsActive.ValueBool(),
		Comment:  plan.Comment.ValueString(),
	}

	asset, err := r.client.CreateAsset(createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating asset",
			fmt.Sprintf("Could not create asset, unexpected error: %s", err),
		)
		return
	}

	// Map response body to model
	plan.ID = types.StringValue(asset.ID)
	plan.Name = types.StringValue(asset.Name)
	// Use Addrs from response if available, otherwise use Address
	if asset.Addrs != "" {
		plan.Address = types.StringValue(asset.Addrs)
	} else {
		plan.Address = types.StringValue(asset.Address)
	}
	// Keep platform as name (not ID) for consistency with schema
	plan.Platform = types.StringValue(asset.Platform.Name)

	// Convert nodes back to list
	var nodeIDs []string
	for _, node := range asset.Nodes {
		nodeIDs = append(nodeIDs, node.ID)
	}
	nodeList, diags := types.ListValueFrom(ctx, types.StringType, nodeIDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Nodes = nodeList

	plan.IsActive = types.BoolValue(asset.IsActive)
	plan.Comment = types.StringValue(asset.Comment)

	// Log creation
	tflog.Trace(ctx, "created asset", map[string]any{"id": plan.ID.ValueString()})

	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AssetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AssetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	asset, err := r.client.GetAsset(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading asset",
			fmt.Sprintf("Could not read asset ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}

	// Map response body to model
	state.ID = types.StringValue(asset.ID)
	state.Name = types.StringValue(asset.Name)
	// Use Addrs from response if available, otherwise use Address
	if asset.Addrs != "" {
		state.Address = types.StringValue(asset.Addrs)
	} else {
		state.Address = types.StringValue(asset.Address)
	}
	// Keep platform as name (not ID) for consistency with schema
	state.Platform = types.StringValue(asset.Platform.Name)

	// Convert nodes to list
	var nodeIDs []string
	for _, node := range asset.Nodes {
		nodeIDs = append(nodeIDs, node.ID)
	}
	nodeList, diags := types.ListValueFrom(ctx, types.StringType, nodeIDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.Nodes = nodeList

	state.IsActive = types.BoolValue(asset.IsActive)
	state.Comment = types.StringValue(asset.Comment)

	// Set state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AssetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AssetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert nodes list
	var nodes []string
	diags = plan.Nodes.ElementsAs(ctx, &nodes, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update asset - parse platform as integer
	platformID := 1 // default value
	if plan.Platform.ValueString() != "" {
		fmt.Sscanf(plan.Platform.ValueString(), "%d", &platformID)
	}

	// Convert nodes to NodeRequest format
	var nodeReqs []jumpserver.NodeRequest
	for _, nodeID := range nodes {
		nodeReqs = append(nodeReqs, jumpserver.NodeRequest{PK: nodeID})
	}

	updateReq := &jumpserver.UpdateAssetRequest{
		Name:     plan.Name.ValueString(),
		Address:  plan.Address.ValueString(),
		Platform: jumpserver.PlatformRequest{PK: platformID},
		Nodes:    nodeReqs,
		IsActive: &[]bool{plan.IsActive.ValueBool()}[0],
		Comment:  plan.Comment.ValueString(),
	}

	asset, err := r.client.UpdateAsset(plan.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating asset",
			fmt.Sprintf("Could not update asset ID %s: %s", plan.ID.ValueString(), err),
		)
		return
	}

	// Map response body to model
	plan.Name = types.StringValue(asset.Name)
	// Use Addrs from response if available, otherwise use Address
	if asset.Addrs != "" {
		plan.Address = types.StringValue(asset.Addrs)
	} else {
		plan.Address = types.StringValue(asset.Address)
	}
	// Keep platform as name (not ID) for consistency with schema
	plan.Platform = types.StringValue(asset.Platform.Name)

	// Convert nodes back to list
	var nodeIDs []string
	for _, node := range asset.Nodes {
		nodeIDs = append(nodeIDs, node.ID)
	}
	nodeList, diags := types.ListValueFrom(ctx, types.StringType, nodeIDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.Nodes = nodeList

	plan.IsActive = types.BoolValue(asset.IsActive)
	plan.Comment = types.StringValue(asset.Comment)

	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *AssetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AssetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAsset(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting asset",
			fmt.Sprintf("Could not delete asset ID %s: %s", state.ID.ValueString(), err),
		)
		return
	}

	// Log deletion
	tflog.Trace(ctx, "deleted asset", map[string]any{"id": state.ID.ValueString()})
}

func (r *AssetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
