package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	niosclient "github.com/unasra/nios-go-client/client"
	"github.com/unasra/nios-go-client/option"
	"github.com/unasra/terraform-provider-nios/internal/service/dns"
)

// Ensure NIOSProvider satisfies various provider interfaces.
var _ provider.Provider = &NIOSProvider{}

// NIOSProvider defines the provider implementation.
type NIOSProvider struct {
	version string
	commit  string
}

// NIOSProviderModel describes the provider data model.
type NIOSProviderModel struct {
	NIOSHostURL types.String `tfsdk:"nios_host_url"`
	NIOSAuth    types.String `tfsdk:"nios_auth"`
}

func (p *NIOSProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "nios"
	resp.Version = p.version
}

func (p *NIOSProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The NIOS provider is used to interact with the resources supported by Infoblox NIOS WAPI.",
		Attributes: map[string]schema.Attribute{
			"nios_host_url": schema.StringAttribute{
				Optional: true,
			},
			"nios_auth": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *NIOSProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data NIOSProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	client := niosclient.NewAPIClient(
		option.WithClientName(fmt.Sprintf("terraform/%s#%s", p.version, p.commit)),
		option.WithNIOSAuth(data.NIOSAuth.ValueString()),
		option.WithNIOSHostUrl(data.NIOSHostURL.ValueString()),
	)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *NIOSProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		dns.NewRecordaResource,
	}
}

func (p *NIOSProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		dns.NewRecordaDataSource,
	}
}

func New(version, commit string) func() provider.Provider {
	return func() provider.Provider {
		return &NIOSProvider{
			version: version,
			commit:  commit,
		}
	}
}
