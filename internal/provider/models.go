package provider

import (
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type imageResourceModel struct {
	ImageURL         types.String   `tfsdk:"image_url"`
	Architecture     types.String   `tfsdk:"architecture"`
	Location         types.String   `tfsdk:"location"`
	ServerType       types.String   `tfsdk:"server_type"`
	Labels           types.Map      `tfsdk:"labels"`
	Description      types.String   `tfsdk:"description"`
	DebugSkipCleanup types.Bool     `tfsdk:"debug_skip_cleanup"`
	EffectiveLabels  types.Map      `tfsdk:"effective_labels"`
	ID               types.String   `tfsdk:"id"`
	ImageID          types.String   `tfsdk:"image_id"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
}

func parseImageID(raw string) (int64, error) {
	return strconv.ParseInt(raw, 10, 64)
}

func stringPointerFromTypes(v types.String) *string {
	if v.IsNull() || v.ValueString() == "" {
		return nil
	}
	s := v.ValueString()
	return &s
}

func boolValue(v types.Bool) bool {
	if v.IsNull() {
		return false
	}
	return v.ValueBool()
}

func stringMapToAttrValues(in map[string]string) map[string]attr.Value {
	out := make(map[string]attr.Value, len(in))
	for k, v := range in {
		out[k] = types.StringValue(v)
	}
	return out
}
