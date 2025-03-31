package dns

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/unasra/terraform-provider-nios/internal/utils"
	"net/http"

	niosclient "github.com/unasra/nios-go-client/client"
)

var readableAttributes = "aws_rte53_record_info,cloud_info,comment,creation_time,creator,ddns_principal,ddns_protected,disable,discovered_data,dns_name,extattrs,forbid_reclamation,ipv4addr,last_queried,ms_ad_user_data,name,reclaimable,shared_record_group,ttl,use_ttl,view,zone"

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &RecordaResource{}
var _ resource.ResourceWithImportState = &RecordaResource{}

func NewRecordaResource() resource.Resource {
	return &RecordaResource{}
}

// RecordaResource defines the resource implementation.
type RecordaResource struct {
	client *niosclient.APIClient
}

func (r *RecordaResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + "dns_a_record"
}

func (r *RecordaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "",
		Attributes:          RecordAResourceSchemaAttributes,
	}
}

func (r *RecordaResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*niosclient.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *niosclient.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *RecordaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RecordAModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiRes, httpres, err := r.client.DNSAPI.
		RecordaAPI.
		Post(ctx).
		RecordA(*data.Expand(ctx, &resp.Diagnostics, true)).
		ReturnFields2(readableAttributes).
		ReturnAsObject(1).
		Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create Recorda, got error: %s", err))
		return
	}
	if httpres == nil {
		httpres = nil
	}
	res := apiRes.GetResult()
	data.Ref = types.StringPointerValue(res.Ref)
	data.Id = types.StringValue(utils.ExtractResourceRef(data.Ref.ValueString()))

	// Save data into Terraform state
	data.Flatten(ctx, &res, &resp.Diagnostics)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RecordaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data RecordAModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiRes, httpRes, err := r.client.DNSAPI.
		RecordaAPI.
		RecordaReferenceGet(ctx, data.Id.ValueString()).
		//ReturnFields("ref,aws_rte53_record_info,cloud_info, comment, creation_time, creator, ddns_principal, ddns_protected, disable, discovered_data, dns_name, extattrs, forbid_reclamation, ipv4addr, last_queried, ms_ad_user_data, name, reclaimable, remove_associated_ptr, shared_record_group, ttl, use_ttl, view, zone").
		ReturnFields2(readableAttributes).
		ReturnAsObject(1).
		Execute()
	if err != nil {
		if httpRes != nil && httpRes.StatusCode == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Recorda, got error: %s", err))
		return
	}

	res := apiRes.GetResult()
	data.Flatten(ctx, &res, &resp.Diagnostics)

	//data.Ref = types.StringPointerValue(apiRes.GetResult().Ref)
	//data.Id = types.StringValue(utils.ExtractResourceRef(data.Ref.ValueString()))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RecordaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data RecordAModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiRes, _, err := r.client.DNSAPI.
		RecordaAPI.
		RecordaReferencePut(ctx, data.Id.ValueString()).
		RecordA(*data.Expand(ctx, &resp.Diagnostics, false)).
		ReturnFields2(readableAttributes).
		ReturnAsObject(1).
		Execute()
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update Recorda, got error: %s", err))
		return
	}

	//data.Ref = types.StringValue(apiRes)
	//data.Id = types.StringValue(utils.ExtractResourceRef(data.Ref.ValueString()))

	res := apiRes.GetResult()
	data.Flatten(ctx, &res, &resp.Diagnostics)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RecordaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data RecordAModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpRes, err := r.client.DNSAPI.
		RecordaAPI.
		RecordaReferenceDelete(ctx, data.Id.ValueString()).
		Execute()
	if err != nil {
		if httpRes != nil && httpRes.StatusCode == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete Recorda, got error: %s", err))
		return
	}
}

func (r *RecordaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
