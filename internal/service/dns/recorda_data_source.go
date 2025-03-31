package dns

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"

	niosclient "github.com/unasra/nios-go-client/client"
	"github.com/unasra/nios-go-client/dns"
	"github.com/unasra/terraform-provider-nios/internal/flex"
	"github.com/unasra/terraform-provider-nios/internal/utils"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &RecordaDataSource{}

func NewRecordaDataSource() datasource.DataSource {
	return &RecordaDataSource{}
}

// RecordaDataSource defines the data source implementation.
type RecordaDataSource struct {
	client *niosclient.APIClient
}

func (d *RecordaDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + "dns_a_records"
}

type RecordAModelWithFilter struct {
	Result types.List `tfsdk:"result"`
	Body   types.Map  `tfsdk:"body"`
}

func (m *RecordAModelWithFilter) FlattenResults(ctx context.Context, from []dns.RecordA, diags *diag.Diagnostics) {
	if len(from) == 0 {
		return
	}
	m.Result = flex.FlattenFrameworkListNestedBlock(ctx, from, RecordAAttrTypes, diags, FlattenRecordA)
}

func (d *RecordaDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "",
		Attributes: map[string]schema.Attribute{
			"result": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: utils.DataSourceAttributeMap(RecordAResourceSchemaAttributes, &resp.Diagnostics),
				},
				Computed: true,
			},
			"body": schema.MapAttribute{
				Description: "The body of the GET Call.",
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

func (d *RecordaDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*niosclient.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *niosclient.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *RecordaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RecordAModelWithFilter

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiRes, httpRes, err := d.client.DNSAPI.
		RecordaAPI.
		Get(ctx).
		ReturnAsObject(1).
		ReturnFields2(readableAttributes).
		Body(flex.ExpandFrameworkMapString(ctx, data.Body, &resp.Diagnostics)).
		Execute()
	if err != nil {
		if httpRes != nil && httpRes.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Record, got error: %s", err))
		return
	}

	res := apiRes.ListRecordAResponseObject.GetResult()
	data.FlattenResults(ctx, res, &resp.Diagnostics)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
