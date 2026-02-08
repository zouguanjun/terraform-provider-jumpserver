package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"jumpserver/internal/jumpserver"
	"jumpserver/internal/provider/data_sources"
	"jumpserver/internal/provider/resources"
)

// Ensure JumpServerProvider satisfies various provider interfaces.
var _ provider.Provider = &JumpServerProvider{}

// JumpServerProvider defines the provider implementation.
type JumpServerProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// JumpServerProviderModel describes the provider data model.
type JumpServerProviderModel struct {
	Endpoint           types.String `tfsdk:"endpoint"`
	KeyID              types.String `tfsdk:"key_id"`
	KeySecret          types.String `tfsdk:"key_secret"`
	OrgID              types.String `tfsdk:"org_id"`
	InsecureSkipVerify types.Bool   `tfsdk:"insecure_skip_verify"`
}

func (p *JumpServerProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "jumpserver"
	resp.Version = p.version
}

func (p *JumpServerProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The JumpServer provider is used to interact with JumpServer resources.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Description: "The JumpServer API endpoint URL. Example: https://jumpserver.example.com",
				Required:    true,
			},
			"key_id": schema.StringAttribute{
				Description: "The JumpServer Access Key ID",
				Required:    true,
				Sensitive:   true,
			},
			"key_secret": schema.StringAttribute{
				Description: "The JumpServer Access Key Secret",
				Required:    true,
				Sensitive:   true,
			},
			"org_id": schema.StringAttribute{
				Description: "The JumpServer organization ID (optional, defaults to default org)",
				Optional:    true,
			},
			"insecure_skip_verify": schema.BoolAttribute{
				Description: "Skip TLS certificate verification (not recommended for production)",
				Optional:    true,
			},
		},
	}
}

func (p *JumpServerProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config JumpServerProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Endpoint.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as endpoint",
		)
		return
	}

	if config.KeyID.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as key_id",
		)
		return
	}

	if config.KeySecret.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as key_secret",
		)
		return
	}

	// Default values to empty strings
	endpoint := config.Endpoint.ValueString()
	keyID := config.KeyID.ValueString()
	keySecret := config.KeySecret.ValueString()
	orgID := ""

	if !config.OrgID.IsNull() {
		orgID = config.OrgID.ValueString()
	}

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("provider"),
			"Missing JumpServer API Endpoint",
			"The provider cannot create the JumpServer API client as there is a missing or empty value for the JumpServer API endpoint. "+
				"Set the endpoint value in the configuration or use the JUMPSERVER_ENDPOINT environment variable.",
		)
		return
	}

	if keyID == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("provider"),
			"Missing JumpServer Access Key ID",
			"The provider cannot create the JumpServer API client as there is a missing or empty value for the JumpServer Access Key ID. "+
				"Set the key_id value in the configuration or use the JUMPSERVER_KEY_ID environment variable.",
		)
		return
	}

	if keySecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("provider"),
			"Missing JumpServer Access Key Secret",
			"The provider cannot create the JumpServer API client as there is a missing or empty value for the JumpServer Access Key Secret. "+
				"Set the key_secret value in the configuration or use the JUMPSERVER_KEY_SECRET environment variable.",
		)
		return
	}

	insecureSkipVerify := !config.InsecureSkipVerify.IsNull() && config.InsecureSkipVerify.ValueBool()

	client := jumpserver.NewClient(&jumpserver.Config{
		Endpoint:           endpoint,
		KeyID:              keyID,
		KeySecret:          keySecret,
		OrgID:              orgID,
		InsecureSkipVerify: insecureSkipVerify,
	})

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *JumpServerProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewAssetResource,
		resources.NewAccountResource,
		resources.NewPermissionResource,
		resources.NewUserResource,
	}
}

func (p *JumpServerProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		data_sources.NewAssetDataSource,
		data_sources.NewPlatformDataSource,
		data_sources.NewNodeDataSource,
		data_sources.NewUserDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &JumpServerProvider{
			version: version,
		}
	}
}
