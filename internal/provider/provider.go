package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type imagerProvider struct {
	version string
}

type imagerProviderModel struct {
	Token types.String `tfsdk:"token"`
}

type providerData struct {
	hcloudClient *hcloud.Client
}

var _ provider.Provider = (*imagerProvider)(nil)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &imagerProvider{
			version: version,
		}
	}
}

func (p *imagerProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "imager"
	resp.Version = p.version
}

func (p *imagerProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Hetzner Cloud API token. Can also be set via the HCLOUD_TOKEN environment variable.",
			},
		},
	}
}

func (p *imagerProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	token := os.Getenv("HCLOUD_TOKEN")

	var config imagerProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !config.Token.IsNull() && config.Token.ValueString() != "" {
		token = config.Token.ValueString()
	}

	if token == "" {
		resp.Diagnostics.AddError(
			"Missing Hetzner Cloud token",
			"Set the token in the provider configuration or via the HCLOUD_TOKEN environment variable.",
		)
		return
	}

	hcloudClient := hcloud.NewClient(hcloud.WithToken(token))
	resp.ResourceData = &providerData{
		hcloudClient: hcloudClient,
	}
}

func (p *imagerProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewImageResource,
	}
}

func (p *imagerProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}
