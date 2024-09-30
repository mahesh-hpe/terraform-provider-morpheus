package morpheusv3

import (
	"context"

	"github.com/gomorpheus/terraform-provider-morpheus/morpheus"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &MorpheusProvider{}

type MorpheusProvider struct {
	version string
}
type MorpheusProviderModel struct {
	Url             types.String `tfsdk:"url"`
	AccessToken     types.String `tfsdk:"access_token"`
	TenantSubdomain types.String `tfsdk:"tenant_subdomain"`
	Username        types.String `tfsdk:"username"`
	Password        types.String `tfsdk:"password"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MorpheusProvider{
			version: version,
		}
	}
}

func (p *MorpheusProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "morpheus"
	resp.Version = p.version
}

func (p *MorpheusProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Required:    true,
				Description: "The URL of the Morpheus Data Appliance where requests will be directed.",
			},
			"access_token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Access Token of Morpheus user. This can be used instead of authenticating with Username and Password.",
			},

			"tenant_subdomain": schema.StringAttribute{
				Optional:    true,
				Description: "The tenant subdomain used for authentication",
			},

			"username": schema.StringAttribute{
				Optional:    true,
				Description: "Username of Morpheus user for authentication",
			},

			"password": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Password of Morpheus user for authentication",
			},
		},
	}

}

func (p *MorpheusProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var providerData MorpheusProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &providerData)...)
	morphConfig := morpheus.Config{
		Url:             providerData.Url.ValueString(),
		AccessToken:     providerData.AccessToken.ValueString(),
		TenantSubdomain: providerData.TenantSubdomain.String(),
		Username:        providerData.Username.ValueString(),
		Password:        providerData.Password.ValueString(),
	}

	morphClient, _ := morphConfig.Client()
	resp.DataSourceData = morphClient
	resp.ResourceData = morphClient
	if resp.Diagnostics.HasError() {
		return
	}
}

func (p *MorpheusProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *MorpheusProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCloudDataSource,
	}
}
