package provider

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/apricote/hcloud-upload-image/hcloudimages"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

type imageResource struct {
	hcloudClient *hcloud.Client
}

var (
	_ resource.Resource                = (*imageResource)(nil)
	_ resource.ResourceWithConfigure   = (*imageResource)(nil)
	_ resource.ResourceWithImportState = (*imageResource)(nil)
)

func NewImageResource() resource.Resource {
	return &imageResource{}
}

func (r *imageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_image"
}

func (r *imageResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Uploads a Talos disk image (raw.xz) into Hetzner Cloud and creates a snapshot image.",
		Attributes: map[string]schema.Attribute{
			"image_url": schema.StringAttribute{
				Required:    true,
				Description: "Public HTTPS URL ending in .raw.xz.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^https://`), "must start with https://"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"architecture": schema.StringAttribute{
				Required:    true,
				Description: "Hetzner snapshot architecture. Valid values: x86 or arm.",
				Validators: []validator.String{
					stringvalidator.OneOf("x86", "arm"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"location": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("fsn1"),
				Description: "Hetzner location name for the temporary upload server. Defaults to fsn1.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"server_type": schema.StringAttribute{
				Optional:    true,
				Description: "Hetzner server type name for the temporary upload server.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"labels": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Additional labels applied to the resulting snapshot image.",
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description applied to the resulting snapshot image.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"debug_skip_cleanup": schema.BoolAttribute{
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Skip cleanup of temporary server and ssh key. Only use for debugging.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"effective_labels": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Labels currently set on the snapshot image.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Hetzner snapshot image ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"image_id": schema.StringAttribute{
				Computed:    true,
				Description: "Hetzner snapshot image ID (same as id).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Delete: true,
			}),
		},
	}
}

func (r *imageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected provider data type",
			fmt.Sprintf("Expected *providerData, got: %T.", req.ProviderData),
		)
		return
	}

	r.hcloudClient = data.hcloudClient
}

func (r *imageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan imageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, 10*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	parsedURL, err := validateAndParseImageURL(plan.ImageURL.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid image_url", err.Error())
		return
	}

	labels := map[string]string{}
	if !plan.Labels.IsNull() {
		resp.Diagnostics.Append(plan.Labels.ElementsAs(ctx, &labels, true)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	arch, archLabel, err := mapArchitecture(plan.Architecture.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid architecture", err.Error())
		return
	}

	effective := mergeLabels(defaultLabels(archLabel), labels)

	opts := hcloudimages.UploadOptions{
		ImageURL:                 parsedURL,
		ImageCompression:         hcloudimages.CompressionXZ,
		ImageFormat:              hcloudimages.FormatRaw,
		Architecture:             arch,
		Labels:                   effective,
		Description:              stringPointerFromTypes(plan.Description),
		DebugSkipResourceCleanup: boolValue(plan.DebugSkipCleanup),
	}

	if !plan.Location.IsNull() && plan.Location.ValueString() != "" {
		opts.Location = &hcloud.Location{Name: plan.Location.ValueString()}
	}
	if !plan.ServerType.IsNull() && plan.ServerType.ValueString() != "" {
		opts.ServerType = &hcloud.ServerType{Name: plan.ServerType.ValueString()}
	}

	tflog.Info(ctx, "uploading image to Hetzner Cloud")
	if !strings.Contains(parsedURL.Path, "/hcloud-") {
		tflog.Warn(ctx, "image_url does not look like an hcloud disk image (missing '/hcloud-'); continuing")
	}

	client := hcloudimages.NewClient(r.hcloudClient)
	image, uploadErr := client.Upload(ctx, opts)
	if uploadErr != nil {
		resp.Diagnostics.AddError("Upload failed", uploadErr.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%d", image.ID))
	plan.ImageID = plan.ID
	if image.Labels != nil {
		plan.EffectiveLabels = types.MapValueMust(types.StringType, stringMapToAttrValues(image.Labels))
	} else {
		plan.EffectiveLabels = types.MapValueMust(types.StringType, stringMapToAttrValues(effective))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *imageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state imageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readTimeout, diags := state.Timeouts.Read(ctx, 2*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	if state.ID.IsNull() || state.ID.ValueString() == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	imageID, err := parseImageID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid id", err.Error())
		return
	}

	image, _, err := r.hcloudClient.Image.GetByID(ctx, imageID)
	if err != nil {
		resp.Diagnostics.AddError("Read failed", err.Error())
		return
	}
	if image == nil {
		resp.Diagnostics.AddWarning("Snapshot image missing", "The snapshot image was not found in Hetzner Cloud; removing from state.")
		resp.State.RemoveResource(ctx)
		return
	}

	state.ID = types.StringValue(fmt.Sprintf("%d", image.ID))
	state.ImageID = state.ID
	if image.Labels != nil {
		state.EffectiveLabels = types.MapValueMust(types.StringType, stringMapToAttrValues(image.Labels))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *imageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"All configurable attributes require replacement. This operation should not be called; please report this as a provider bug.",
	)
}

func (r *imageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state imageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteTimeout, diags := state.Timeouts.Delete(ctx, 10*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, deleteTimeout)
	defer cancel()

	if state.ID.IsNull() || state.ID.ValueString() == "" {
		return
	}

	imageID, err := parseImageID(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid id", err.Error())
		return
	}

	image := &hcloud.Image{ID: imageID}
	_, err = r.hcloudClient.Image.Delete(ctx, image)
	if err != nil {
		if hcloud.IsError(err, hcloud.ErrorCodeNotFound) {
			return
		}
		resp.Diagnostics.AddError("Delete failed", err.Error())
	}
}

func (r *imageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func validateAndParseImageURL(raw string) (*url.URL, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	if parsed.Scheme != "https" {
		return nil, fmt.Errorf("scheme must be https")
	}
	if parsed.Hostname() == "" {
		return nil, fmt.Errorf("host must be set")
	}
	if !strings.HasSuffix(parsed.Path, ".raw.xz") {
		return nil, fmt.Errorf("path must end with .raw.xz")
	}
	return parsed, nil
}

func mapArchitecture(raw string) (hcloud.Architecture, string, error) {
	switch raw {
	case "x86":
		return hcloud.ArchitectureX86, "x86", nil
	case "arm":
		return hcloud.ArchitectureARM, "arm", nil
	default:
		return "", "", fmt.Errorf("expected x86 or arm")
	}
}
