package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func validateHcloudResourceLabels(labels map[string]string) error {
	raw := make(map[string]interface{}, len(labels))
	for key, value := range labels {
		raw[key] = value
	}

	ok, err := hcloud.ValidateResourceLabels(raw)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("labels are not correctly formatted")
	}
	return nil
}

type hcloudResourceLabelsValidator struct{}

var _ validator.Map = hcloudResourceLabelsValidator{}

func (hcloudResourceLabelsValidator) Description(_ context.Context) string {
	return "Ensures labels follow Hetzner Cloud label rules."
}

func (hcloudResourceLabelsValidator) MarkdownDescription(_ context.Context) string {
	return "Ensures labels follow Hetzner Cloud label rules."
}

func (hcloudResourceLabelsValidator) ValidateMap(ctx context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	labels := map[string]string{}
	resp.Diagnostics.Append(req.ConfigValue.ElementsAs(ctx, &labels, true)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := validateHcloudResourceLabels(labels); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid labels",
			fmt.Sprintf(
				"%s. Note: Hetzner Cloud label values may contain only letters/numbers plus '-', '_' and '.', and must start/end with a letter or number (no '/').",
				err.Error(),
			),
		)
	}
}
