package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
	"net/http"
)

var _ resource.Resource = &LimitResource{}
var _ resource.ResourceWithImportState = &LimitResource{}

func NewLimitResource() resource.Resource {

	return &LimitResource{}
}

type LimitResource struct {

	// Client for interaction
	// with RabbitMQ HTTP API
	client *rabbithole.Client
}

type LimitResourceModel struct {
	Id     types.String  `tfsdk:"id"`
	Scope  types.String  `tfsdk:"scope"`
	Name   types.String  `tfsdk:"name"`
	Limits types.MapType `tfsdk:"limits"`
}

const limitVhost = "vhost"
const limitUser = "user"

func (r *LimitResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_limit"
}

func (r *LimitResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				Computed: true,

				PlanModifiers: []planmodifier.String{

					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"scope": schema.StringAttribute{

				Required: true,

				Validators: []validator.String{

					stringvalidator.OneOf(limitVhost, limitUser),
				},
			},

			"name": schema.StringAttribute{

				Required: true,
			},

			"limits": schema.MapAttribute{

				Required: true,

				ElementType: types.NumberType,
			},
		},
	}
}

func (r *LimitResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	client, _ := req.ProviderData.(*rabbithole.Client)

	r.client = client
}

func (r *LimitResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data LimitResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return
	}

	/////////////
	/// LOGIC ///
	/////////////

	// TODO: ADD LOGS
	//tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Attempting to create limits: %v for %v in %v", data.Limit, data.Name, data.Scope))

	var res *http.Response
	var err error

	switch data.Scope.ValueString() {

	case "vhost":

		res, err = r.client.PutVhostLimits(data.Name.ValueString(), getLimitSettings(&data))
		break

	case "user":

		res, err = r.client.PutUserLimits(data.Name.ValueString(), getLimitSettings(&data))
		break
	}

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Limits creation response: %#v", res))

	/////////////
	/// LOGIC ///
	/////////////

	if err != nil {

		resp.Diagnostics.AddError("RabbitMQ: Limits creation error", err.Error())

		return
	}

	data.Id = types.StringValue(data.Scope.ValueString() + "@" + data.Name.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LimitResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data LimitResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return
	}

	/////////////
	/// LOGIC ///
	/////////////

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Attempting to extract limit: %v", data.Id))

	limits := make(map[string]int)
	var err error

	switch data.Scope.ValueString() {

	case limitVhost:

		var vhostLimits []rabbithole.VhostLimitsInfo
		vhostLimits, err = r.client.GetVhostLimits(data.Name.ValueString())

		if vhostLimits != nil {

			limits = vhostLimits[0].Value
		}

		break

	case limitUser:

		var userLimits []rabbithole.UserLimitsInfo
		userLimits, err = r.client.GetUserLimits(data.Name.ValueString())

		if userLimits != nil {

			limits = userLimits[0].Value
		}
		break
	}

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Limit extraction response: %#v", limits))

	/////////////
	/// LOGIC ///
	/////////////

	if err != nil {

		if checkDeleted(ctx, &resp.State, err) != nil {

			resp.Diagnostics.AddError("RabbitMQ: Limit extraction error", err.Error())
		}
	} else {

		_, diags := types.MapValueFrom(ctx, data.Limits, limits)

		if diags.HasError() {

			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *LimitResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// NO-OP
}

func (r *LimitResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data LimitResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return
	}

	/////////////
	/// LOGIC ///
	/////////////

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Attempting to delete limits: %#v", data.Id))

	var res *http.Response
	var err error

	switch data.Scope.ValueString() {

	case "vhost":

		res, err = r.client.DeleteVhostLimits(data.Name.ValueString(), nil)
		break

	case "user":

		res, err = r.client.DeleteUserLimits(data.Name.ValueString(), nil)
		break
	}

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Limits deletion response: %#v", res))

	/////////////
	/// LOGIC ///
	/////////////

	if checkDeleted(ctx, &req.State, err) != nil {

		resp.Diagnostics.AddError("RabbitMQ: Limits deletion error", err.Error())
	}
}

func (r *LimitResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resource.ImportStatePassthroughID(ctx, path.Root("scope"), req, resp)
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
	resource.ImportStatePassthroughID(ctx, path.Root("limit"), req, resp)
	resource.ImportStatePassthroughID(ctx, path.Root("value"), req, resp)
}

func getLimitSettings(settings *LimitResourceModel) map[string]int {

	arguments := make(map[string]int)

	_ = json.Unmarshal([]byte(settings.Limits.String()), &arguments)

	return arguments
}
