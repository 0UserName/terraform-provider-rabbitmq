package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	rabbithole "github.com/michaelklishin/rabbit-hole/v2"
)

var _ resource.Resource = &BindingResource{}
var _ resource.ResourceWithImportState = &BindingResource{}

func NewBindingResource() resource.Resource {

	return &BindingResource{}
}

type BindingResource struct {

	// Client for interaction
	// with RabbitMQ HTTP API
	client *rabbithole.Client
}

type BindingSettingsResourceModel struct {
	Type      types.String  `tfsdk:"type"`
	Arguments types.Dynamic `tfsdk:"arguments"`
}

type BindingResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Vhost       types.String `tfsdk:"vhost"`
	Source      types.String `tfsdk:"source"`
	Destination types.String `tfsdk:"destination"`

	// Consists of a routing
	// key and a hash of its
	// arguments.
	PropertiesKey types.String                 `tfsdk:"properties_key"`
	RoutingKey    types.String                 `tfsdk:"routing_key"`
	Settings      BindingSettingsResourceModel `tfsdk:"arguments"`
}

func (r *BindingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_binding"
}

func (r *BindingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				Computed: true,

				PlanModifiers: []planmodifier.String{

					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"vhost": schema.StringAttribute{

				Required: true,
			},

			"source": schema.StringAttribute{

				Required: true,
			},

			"destination": schema.StringAttribute{

				Required: true,
			},

			"properties_key": schema.StringAttribute{

				Required: true,
			},

			"routing_key": schema.StringAttribute{

				Computed: true,
			},

			"arguments": schema.DynamicAttribute{

				Optional: true,
			},
		},
	}
}

func (r *BindingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	client, _ := req.ProviderData.(*rabbithole.Client)

	r.client = client
}

func (r *BindingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data BindingResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return
	}

	/////////////
	/// LOGIC ///
	/////////////

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Attempting to create binding: %v", data.Name))

	res, err := r.client.DeclareBinding(data.Vhost.ValueString(), getBindingSettings(ctx, &data.Settings))

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Binding creation response: %#v", res))

	/////////////
	/// LOGIC ///
	/////////////

	if err != nil {

		resp.Diagnostics.AddError("RabbitMQ: Binding creation error", err.Error())

		return
	}

	data.Id = types.StringValue(data.Vhost.ValueString() + data.Name.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BindingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data BindingResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return
	}

	/////////////
	/// LOGIC ///
	/////////////

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Attempting to extract binding: %v", data.Id))

	bindings, err := r.client.ListBindingsIn(data.Vhost.ValueString())

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Binding extraction response: %#v", binding))

	/////////////
	/// LOGIC ///
	/////////////

	if err != nil {

		if checkDeleted(ctx, &resp.State, err) != nil {

			resp.Diagnostics.AddError("RabbitMQ: Binding extraction error", err.Error())
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BindingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// NO-OP
}

func (r *BindingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data BindingResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return
	}

	/////////////
	/// LOGIC ///
	/////////////

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Attempting to delete binding: %#v", data.Id))

	res, err := r.client.DeleteBinding(data.Vhost.ValueString(), getBindingSettings(&data))

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Binding deletion response: %#v", res))

	/////////////
	/// LOGIC ///
	/////////////

	if checkDeleted(ctx, &req.State, err) != nil {

		resp.Diagnostics.AddError("RabbitMQ: Binding deletion error", err.Error())
	}
}

func (r *BindingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func getBindingSettings(ctx context.Context, data *BindingResourceModel) rabbithole.BindingInfo {

	arguments := make(map[string]any)

	_ = json.Unmarshal([]byte(data.Settings.Arguments.String()), &arguments)

	return rabbithole.BindingInfo{

		Vhost: data.Vhost.ValueString(),

		Source:      data.Source.ValueString(),
		Destination: data.Destination.ValueString(),

		DestinationType: data.Settings.Type.ValueString(),

		RoutingKey: data.RoutingKey.ValueString(),

		PropertiesKey: data.PropertiesKey.ValueString(),

		Arguments: arguments,
	}
}
