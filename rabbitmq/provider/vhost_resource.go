package provider

import (
	"context"
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

var _ resource.Resource = &VhostResource{}
var _ resource.ResourceWithImportState = &VhostResource{}

func NewVhostResource() resource.Resource {

	return &VhostResource{}
}

type VhostResource struct {

	// Client for interaction
	// with RabbitMQ HTTP API
	client *rabbithole.Client
}

type VhostResourceModel struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	QueueType types.String `tfsdk:"queue_type"`
}

func (r *VhostResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_vhost"
}

func (r *VhostResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{

				Computed: true,

				PlanModifiers: []planmodifier.String{

					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"name": schema.StringAttribute{

				Required: true,
			},

			"queue_type": schema.StringAttribute{

				Optional: true,
			},
		},
	}
}

func (r *VhostResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	client, _ := req.ProviderData.(*rabbithole.Client)

	r.client = client
}

func (r *VhostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data VhostResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return
	}

	/////////////
	/// LOGIC ///
	/////////////

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Attempting to create vhost: %v", data.Name))

	res, err := r.client.PutVhost(data.Name.ValueString(), getVhostSettings(&data))

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Vhost creation response: %#v", res))

	/////////////
	/// LOGIC ///
	/////////////

	if err != nil {

		resp.Diagnostics.AddError("RabbitMQ: Vhost creation error", err.Error())

		return
	}

	data.Id = data.Name

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VhostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var data VhostResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return
	}

	/////////////
	/// LOGIC ///
	/////////////

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Attempting to extract vhost: %v", data.Id))

	vhost, err := r.client.GetVhost(data.Id.ValueString())

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Vhost extraction response: %#v", vhost))

	/////////////
	/// LOGIC ///
	/////////////

	if err != nil {

		if checkDeleted(ctx, &resp.State, err) != nil {

			resp.Diagnostics.AddError("RabbitMQ: Vhost extraction error", err.Error())
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VhostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data VhostResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return
	}

	/////////////
	/// LOGIC ///
	/////////////

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Attempting to update vhost %v", data.Id))

	res, err := r.client.PutVhost(data.Name.ValueString(), getVhostSettings(&data))

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Vhost update response: %#v", res))

	/////////////
	/// LOGIC ///
	/////////////

	if err != nil {

		resp.Diagnostics.AddError("RabbitMQ: Vhost update error", err.Error())
	}
}

func (r *VhostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data VhostResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {

		return
	}

	/////////////
	/// LOGIC ///
	/////////////

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Attempting to delete vhost: %v", data.Id))

	res, err := r.client.DeleteVhost(data.Id.ValueString())

	tflog.Debug(ctx, fmt.Sprintf("RabbitMQ: Vhost deletion response: %#v", res))

	/////////////
	/// LOGIC ///
	/////////////

	if checkDeleted(ctx, &req.State, err) != nil {

		resp.Diagnostics.AddError("RabbitMQ: Vhost deletion error", err.Error())
	}
}

func (r *VhostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
	resource.ImportStatePassthroughID(ctx, path.Root("queue_type"), req, resp)
}

func getVhostSettings(data *VhostResourceModel) rabbithole.VhostSettings {

	return rabbithole.VhostSettings{DefaultQueueType: data.QueueType.ValueString()}
}
