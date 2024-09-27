package morpheusv3

import (
	"context"

	"github.com/gomorpheus/terraform-provider-morpheus/morpheus"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = &MorpheusProvider{}

type MorpheusProvider struct {
	version string
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
	var data morpheus.Config

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	resp.DataSourceData, _ = data.Client()
	if resp.Diagnostics.HasError() {
		return
	}
}

func (p *MorpheusProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource {
			return nil
		},
	}
}

func (p *MorpheusProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCloudDataSource,
	}
}
