// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfsdk

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschemadata"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// EphemeralResultData represents the data returned after opening a Terraform ephemeral resource.
type EphemeralResultData struct {
	Raw    tftypes.Value
	Schema fwschema.Schema
}

// Get populates the struct passed as `target` with the entire ephemeral result data object.
func (s EphemeralResultData) Get(ctx context.Context, target interface{}) diag.Diagnostics {
	return s.data().Get(ctx, target)
}

// GetAttribute retrieves the attribute or block found at `path` and populates
// the `target` with the value. This method is intended for top level schema
// attributes or blocks. Use `types` package methods or custom types to step
// into collections.
//
// Attributes or elements under null or unknown collections return null
// values, however this behavior is not protected by compatibility promises.
func (s EphemeralResultData) GetAttribute(ctx context.Context, path path.Path, target interface{}) diag.Diagnostics {
	return s.data().GetAtPath(ctx, path, target)
}

// PathMatches returns all matching path.Paths from the given path.Expression.
//
// If a parent path is null or unknown, which would prevent a full expression
// from matching, the parent path is returned rather than no match to prevent
// false positives.
func (s EphemeralResultData) PathMatches(ctx context.Context, pathExpr path.Expression) (path.Paths, diag.Diagnostics) {
	return s.data().PathMatches(ctx, pathExpr)
}

// Set populates the entire ephemeral result data object using the supplied Go value. The value `val`
// should be a struct whose values have one of the attr.Value types. Each field
// must be tagged with the corresponding schema field.
func (s *EphemeralResultData) Set(ctx context.Context, val interface{}) diag.Diagnostics {
	data := s.data()
	diags := data.Set(ctx, val)

	if diags.HasError() {
		return diags
	}

	s.Raw = data.TerraformValue

	return diags
}

// SetAttribute sets the attribute at `path` using the supplied Go value.
//
// The attribute path and value must be valid with the current schema. If the
// attribute path already has a value, it will be overwritten. If the attribute
// path does not have a value, it will be added, including any parent attribute
// paths as necessary.
//
// The value must not be an untyped nil. Use a typed nil or types package null
// value function instead. For example with a types.StringType attribute,
// use (*string)(nil) or types.StringNull().
//
// Lists can only have the next element added according to the current length.
func (s *EphemeralResultData) SetAttribute(ctx context.Context, path path.Path, val interface{}) diag.Diagnostics {
	data := s.data()
	diags := data.SetAtPath(ctx, path, val)

	if diags.HasError() {
		return diags
	}

	s.Raw = data.TerraformValue

	return diags
}

func (s EphemeralResultData) data() *fwschemadata.Data {
	return &fwschemadata.Data{
		Description:    fwschemadata.DataDescriptionEphemeralResultData,
		Schema:         s.Schema,
		TerraformValue: s.Raw,
	}
}
