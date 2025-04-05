package dns

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/unasra/nios-go-client/dns"

	"github.com/unasra/terraform-provider-nios/internal/flex"
)

type RecordAModel struct {
	Ref                 types.String `tfsdk:"ref"`
	AwsRte53RecordInfo  types.String `tfsdk:"aws_rte53_record_info"`
	CloudInfo           types.String `tfsdk:"cloud_info"`
	Comment             types.String `tfsdk:"comment"`
	CreationTime        types.Int32  `tfsdk:"creation_time"`
	Creator             types.String `tfsdk:"creator"`
	DdnsPrincipal       types.String `tfsdk:"ddns_principal"`
	DdnsProtected       types.Bool   `tfsdk:"ddns_protected"`
	Disable             types.Bool   `tfsdk:"disable"`
	DiscoveredData      types.String `tfsdk:"discovered_data"`
	DnsName             types.String `tfsdk:"dns_name"`
	Extattrs            types.Map    `tfsdk:"extattrs"`
	ForbidReclamation   types.Bool   `tfsdk:"forbid_reclamation"`
	Ipv4addr            types.String `tfsdk:"ipv4addr"`
	LastQueried         types.String `tfsdk:"last_queried"`
	MsAdUserData        types.String `tfsdk:"ms_ad_user_data"`
	Name                types.String `tfsdk:"name"`
	Reclaimable         types.Bool   `tfsdk:"reclaimable"`
	RemoveAssociatedPtr types.Bool   `tfsdk:"remove_associated_ptr"`
	SharedRecordGroup   types.String `tfsdk:"shared_record_group"`
	Ttl                 types.Int32  `tfsdk:"ttl"`
	UseTtl              types.Bool   `tfsdk:"use_ttl"`
	View                types.String `tfsdk:"view"`
	Zone                types.String `tfsdk:"zone"`
}

var RecordAAttrTypes = map[string]attr.Type{
	"ref":                   types.StringType,
	"aws_rte53_record_info": types.StringType,
	"cloud_info":            types.StringType,
	"comment":               types.StringType,
	"creation_time":         types.Int32Type,
	"creator":               types.StringType,
	"ddns_principal":        types.StringType,
	"ddns_protected":        types.BoolType,
	"disable":               types.BoolType,
	"discovered_data":       types.StringType,
	"dns_name":              types.StringType,
	"extattrs":              types.MapType{ElemType: types.MapType{ElemType: types.StringType}},
	"forbid_reclamation":    types.BoolType,
	"ipv4addr":              types.StringType,
	"last_queried":          types.StringType,
	"ms_ad_user_data":       types.StringType,
	"name":                  types.StringType,
	"reclaimable":           types.BoolType,
	"remove_associated_ptr": types.BoolType,
	"shared_record_group":   types.StringType,
	"ttl":                   types.Int32Type,
	"use_ttl":               types.BoolType,
	"view":                  types.StringType,
	"zone":                  types.StringType,
}

var RecordAResourceSchemaAttributes = map[string]schema.Attribute{
	"ref": schema.StringAttribute{
		Computed: true,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The reference to the object.",
	},
	"aws_rte53_record_info": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "Aws Route 53 record information.",
	},
	"cloud_info": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "Structure containing all cloud API related information for this object.",
	},
	"comment": schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "Comment for the record; maximum 256 characters.",
	},
	"creation_time": schema.Int32Attribute{
		Computed:            true,
		MarkdownDescription: "The time of the record creation in Epoch seconds format.",
	},
	"creator": schema.StringAttribute{
		Optional: true,
		Computed: true,
		Validators: []validator.String{
			stringvalidator.OneOf("STATIC", "DYNAMIC"),
		},
		Default:             stringdefault.StaticString("STATIC"),
		MarkdownDescription: "The record creator.",
	},
	"ddns_principal": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "The GSS-TSIG principal that owns this record.",
	},
	"ddns_protected": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Determines if the DDNS updates for this record are allowed or not.",
	},
	"disable": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Determines if the record is disabled or not. False means that the record is enabled.",
	},
	"discovered_data": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "The discovered data for this A record.",
	},
	"dns_name": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "The name for an A record in punycode format.",
	},
	"extattrs": schema.MapAttribute{
		Optional:    true,
		Computed:    true,
		ElementType: types.MapType{ElemType: types.StringType},
		//Default:             mapdefault.StaticValue(types.MapNull(types.MapType{ElemType: types.StringType})),
		MarkdownDescription: "Extensible attributes associated with the object.",
	},
	"forbid_reclamation": schema.BoolAttribute{
		Optional:            true,
		Computed:            true,
		Default:             booldefault.StaticBool(false),
		MarkdownDescription: "Determines if the reclamation is allowed for the record or not.",
	},
	"ipv4addr": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The IPv4 Address of the record.",
	},
	"last_queried": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "The time of the last DNS query in Epoch seconds format.",
	},
	"ms_ad_user_data": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "The Microsoft Active Directory user related information.",
	},
	"name": schema.StringAttribute{
		Required:            true,
		MarkdownDescription: "The Name of the record.",
	},
	"reclaimable": schema.BoolAttribute{
		Computed:            true,
		MarkdownDescription: "Determines if the record is reclaimable or not.",
	},
	"remove_associated_ptr": schema.BoolAttribute{
		Optional:            true,
		MarkdownDescription: "Whether to remove associated PTR records while deleting the A record.",
	},
	"shared_record_group": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "The shared record group this record belongs to.",
	},
	"ttl": schema.Int32Attribute{
		Optional:            true,
		MarkdownDescription: "Time-to-live value of the record, in seconds.",
		Validators: []validator.Int32{
			int32validator.AlsoRequires(path.MatchRoot("use_ttl")),
		},
	},
	"use_ttl": schema.BoolAttribute{
		Optional: true,
		Computed: true,
		Default:  booldefault.StaticBool(false),
		Validators: []validator.Bool{
			boolvalidator.AlsoRequires(path.MatchRoot("ttl")),
		},
		MarkdownDescription: "Flag to indicate whether the TTL value should be used for the A record.",
	},
	"view": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "View that this record is part of.",
	},
	"zone": schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "The zone in which the record resides.",
	},
}

func (m *RecordAModel) Expand(ctx context.Context, diags *diag.Diagnostics, isCreate bool) *dns.RecordA {
	if m == nil {
		return nil
	}
	to := &dns.RecordA{
		Comment:             flex.ExpandStringPointer(m.Comment),
		Creator:             flex.ExpandStringPointer(m.Creator),
		DdnsPrincipal:       flex.ExpandStringPointer(m.DdnsPrincipal),
		DdnsProtected:       flex.ExpandBoolPointer(m.DdnsProtected),
		Disable:             flex.ExpandBoolPointer(m.Disable),
		Extattrs:            flex.ExpandFrameworkMapOfMapString(ctx, m.Extattrs, diags),
		ForbidReclamation:   flex.ExpandBoolPointer(m.ForbidReclamation),
		Ipv4addr:            flex.ExpandString(m.Ipv4addr),
		Name:                flex.ExpandString(m.Name),
		RemoveAssociatedPtr: flex.ExpandBoolPointer(m.RemoveAssociatedPtr),
		Ttl:                 flex.ExpandInt32Pointer(m.Ttl),
		UseTtl:              flex.ExpandBoolPointer(m.UseTtl),
	}
	if isCreate {
		to.View = flex.ExpandStringPointer(m.View)
	}
	return to
}

func FlattenRecordA(ctx context.Context, from *dns.RecordA, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(RecordAAttrTypes)
	}
	m := RecordAModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, RecordAAttrTypes, m)
	diags.Append(d...)
	return t
}

func (m *RecordAModel) Flatten(ctx context.Context, from *dns.RecordA, diags *diag.Diagnostics) {
	if from == nil {
		return
	}
	if m == nil {
		*m = RecordAModel{}
	}
	m.Ref = flex.FlattenStringPointer(from.Ref)
	m.AwsRte53RecordInfo = flex.FlattenStringPointer(from.AwsRte53RecordInfo)
	m.CloudInfo = flex.FlattenStringPointer(from.CloudInfo)
	m.Comment = flex.FlattenStringPointer(from.Comment)
	m.CreationTime = flex.FlattenInt32Pointer(from.CreationTime)
	m.Creator = flex.FlattenStringPointer(from.Creator)
	m.DdnsPrincipal = flex.FlattenStringPointer(from.DdnsPrincipal)
	m.DdnsProtected = types.BoolPointerValue(from.DdnsProtected)
	m.Disable = types.BoolPointerValue(from.Disable)
	m.DiscoveredData = flex.FlattenStringPointer(from.DiscoveredData)
	m.DnsName = flex.FlattenStringPointer(from.DnsName)
	m.Extattrs = flex.FlattenFrameworkMapOfMapString(ctx, from.Extattrs, diags)
	m.ForbidReclamation = types.BoolPointerValue(from.ForbidReclamation)
	m.Ipv4addr = flex.FlattenString(from.Ipv4addr)
	m.LastQueried = flex.FlattenStringPointer(from.LastQueried)
	m.MsAdUserData = flex.FlattenStringPointer(from.MsAdUserData)
	m.Name = flex.FlattenString(from.Name)
	m.Reclaimable = types.BoolPointerValue(from.Reclaimable)
	m.RemoveAssociatedPtr = types.BoolPointerValue(from.RemoveAssociatedPtr)
	m.SharedRecordGroup = flex.FlattenStringPointer(from.SharedRecordGroup)
	m.Ttl = flex.FlattenInt32Pointer(from.Ttl)
	m.UseTtl = types.BoolPointerValue(from.UseTtl)
	m.View = flex.FlattenStringPointer(from.View)
	m.Zone = flex.FlattenStringPointer(from.Zone)
}
